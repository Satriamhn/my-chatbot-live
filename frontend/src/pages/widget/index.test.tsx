import { describe, it, expect, vi, beforeEach } from "vitest";
import { act, render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router-dom";

const getSettingsMock = vi.fn();
const postMock = vi.fn();

type WidgetWebSocketOptions = {
  onToken: (chunk: { content: string; is_last: boolean }) => void;
  onError: (msg: string) => void;
};

const widgetWebSocketState = vi.hoisted(() => ({
  options: undefined as WidgetWebSocketOptions | undefined,
}));

vi.mock("../../hooks/useWidgetWebSocket", () => ({
  useWidgetWebSocket: (options: WidgetWebSocketOptions) => {
    widgetWebSocketState.options = options;

    return {
      isConnected: false,
      sendMessage: vi.fn(),
    };
  },
}));

vi.mock("../../services/widget", () => ({
  widgetService: {
    getSettings: (...args: unknown[]) => getSettingsMock(...args),
  },
}));

vi.mock("axios", () => ({
  default: {
    post: (...args: unknown[]) => postMock(...args),
    isAxiosError: (error: unknown) => Boolean((error as { isAxiosError?: boolean } | undefined)?.isAxiosError),
  },
}));

import Widget, { MessageBubble } from "./index";

function renderWidget(entry: string) {
  return render(
    <MemoryRouter initialEntries={[entry]}>
      <Routes>
        <Route path="/widget" element={<Widget />} />
      </Routes>
    </MemoryRouter>,
  );
}

describe("Widget page", () => {
  beforeEach(() => {
    getSettingsMock.mockReset();
    postMock.mockReset();
    widgetWebSocketState.options = undefined;
    document.title = "";
  });

  it("shows a deterministic error when org_id is missing", () => {
    renderWidget("/widget");

    expect(screen.getByRole("alert")).toHaveTextContent("Missing org_id");
    expect(screen.getByRole("alert")).toHaveTextContent("tenant id in the URL");
  });

  it("keeps safe fallback widget content when settings fetch fails", async () => {
    getSettingsMock.mockRejectedValueOnce(new Error("settings unavailable"));
    postMock.mockResolvedValueOnce({ data: { id: "session-1" } });

    renderWidget("/widget?org_id=tenant-123");

    await waitFor(() => expect(getSettingsMock).toHaveBeenCalledWith("tenant-123"));
    await waitFor(() => expect(postMock).toHaveBeenCalled());

    expect(screen.getByRole("heading", { level: 3, name: "Customer Support" })).toBeVisible();
    expect(screen.getByText("Hello! How can we help you today?")).toBeVisible();
    expect(document.title).toBe("Customer Support");
  });

  it("surfaces anonymous session init failures in the UI", async () => {
    getSettingsMock.mockResolvedValueOnce({
      bot_name: "Acme Bot",
      welcome_message: "Welcome to Acme",
    });
    postMock.mockRejectedValueOnce({
      isAxiosError: true,
      response: { data: { error: "Unable to initialize chat session." } },
    });

    renderWidget("/widget?org_id=tenant-123");

    await waitFor(() => expect(getSettingsMock).toHaveBeenCalledWith("tenant-123"));
    await waitFor(() => expect(screen.getByRole("alert")).toHaveTextContent("Unable to initialize chat session."));

    expect(screen.getByRole("heading", { level: 3, name: "Acme Bot" })).toBeVisible();
    expect(screen.getByText("Welcome to Acme")).toBeVisible();
  });

  it("does not show any bot bubble for markdown that renders nothing", async () => {
    getSettingsMock.mockResolvedValueOnce({
      bot_name: "Acme Bot",
      welcome_message: "Welcome to Acme",
    });
    postMock.mockResolvedValueOnce({ data: { id: "session-1" } });

    const { container } = renderWidget("/widget?org_id=tenant-123");

    await waitFor(() => expect(getSettingsMock).toHaveBeenCalledWith("tenant-123"));
    await waitFor(() => expect(postMock).toHaveBeenCalled());
    await waitFor(() => expect(widgetWebSocketState.options).toBeDefined());

    const streamingBubbleSelector = ".mr-auto .bg-white.border.border-gray-100.text-gray-700.rounded-tl-sm.shadow-sm.min-w-0";

    act(() => {
      widgetWebSocketState.options?.onToken({ content: "[foo]: /bar", is_last: true });
    });

    expect(container.querySelector(streamingBubbleSelector)).toBeNull();
    expect(container.querySelectorAll(".mr-auto")).toHaveLength(0);
  });

  it("suppresses an already-queued bot row when markdown renders nothing", () => {
    const { container } = render(
      <MessageBubble
        message={{
          id: "bot-hidden",
          sender: "bot",
          content: "[foo]: /bar",
        }}
      />,
    );

    expect(container).toBeEmptyDOMElement();
  });

  it("still renders visible bot markdown normally", () => {
    const { container } = render(
      <MessageBubble
        message={{
          id: "bot-visible",
          sender: "bot",
          content: "Hello **world**",
        }}
      />,
    );

    expect(container.querySelectorAll(".mr-auto")).toHaveLength(1);
    expect(screen.getByText("Hello")).toBeVisible();
    expect(screen.getByText("world")).toBeVisible();
  });
});
