import { describe, it, expect, vi, beforeEach } from "vitest";
import { AxiosError } from "axios";
import api from "../api";

// Mock localStorage
const localStorageMock = (() => {
  let store: Record<string, string> = {};
  return {
    getItem: (key: string) => store[key] || null,
    setItem: (key: string, value: string) => { store[key] = value.toString(); },
    removeItem: (key: string) => { delete store[key]; },
    clear: () => { store = {}; }
  };
})();
Object.defineProperty(window, "localStorage", { value: localStorageMock });

// Mock window.location
const locationMock = { href: "" };
Object.defineProperty(window, "location", { value: locationMock, writable: true });

describe("API Interceptors", () => {
  beforeEach(() => {
    localStorage.clear();
    locationMock.href = "";
    vi.clearAllMocks();
  });

  it("should add Authorization header if token exists", async () => {
    localStorage.setItem("token", "test-token");
    
    // We can't easily test the internal axios instance state without making a request
    // and intercepting it, but we can check the interceptor logic.
    // Here we'll simulate a request config
    const config = { headers: {} as Record<string, string> };
    const requestInterceptor = (api.interceptors.request as { handlers: { fulfilled: (config: unknown) => unknown }[] }).handlers[0].fulfilled;
    const result = requestInterceptor(config) as { headers: Record<string, string> };

    expect(result.headers.Authorization).toBe("Bearer test-token");
  });

  it("should add X-Org-Id header if orgId exists", async () => {
    localStorage.setItem("orgId", "test-org");
    
    const config = { headers: {} as Record<string, string> };
    const requestInterceptor = (api.interceptors.request as { handlers: { fulfilled: (config: unknown) => unknown }[] }).handlers[0].fulfilled;
    const result = requestInterceptor(config) as { headers: Record<string, string> };

    expect(result.headers["X-Org-ID"]).toBe("test-org");
  });

  it("should handle 401 response and redirect to signin", async () => {
    localStorage.setItem("token", "old-token");
    
    const error = {
      response: {
        status: 401,
        data: { message: "Unauthorized" }
      },
      isAxiosError: true,
      request: {},
      config: {
        url: "/test"
      }
    } as unknown as AxiosError;

    const responseInterceptor = (api.interceptors.response as { handlers: { rejected: (error: unknown) => unknown }[] }).handlers[0].rejected;
    
    vi.stubGlobal('location', { ...locationMock, pathname: '/dashboard' });
    
    try {
      await responseInterceptor(error);
    } catch {
      expect(localStorage.getItem("token")).toBeNull();
      expect(window.location.href).toBe("/signin");
    }
  });

  it("should handle ORG_MISMATCH 403 error", async () => {
    localStorage.setItem("orgId", "wrong-org");
    
    const error = {
      response: {
        status: 403,
        data: { message: "Org Mismatch", code: "ORG_MISMATCH" }
      },
      isAxiosError: true
    } as unknown as AxiosError;

    const responseInterceptor = (api.interceptors.response as { handlers: { rejected: (error: unknown) => unknown }[] }).handlers[0].rejected;
    
    try {
      await responseInterceptor(error);
    } catch (err) {
      const errorObj = err as { message: string };
      expect(localStorage.getItem("orgId")).toBeNull();
      expect(errorObj.message).toBe("Organization access mismatch. Please switch organization.");
    }
  });

  it("should handle request with missing token", async () => {
    // Ensure token is empty
    localStorage.removeItem("token");
    
    const config = { headers: {} as Record<string, string> };
    const requestInterceptor = (api.interceptors.request as { handlers: { fulfilled: (config: unknown) => unknown }[] }).handlers[0].fulfilled;
    const result = requestInterceptor(config) as { headers: Record<string, string> };

    expect(result.headers.Authorization).toBeUndefined();
  });


  it("should normalize error responses", async () => {
    const error = {
      response: {
        status: 400,
        data: { message: "Bad Request", code: "ERR_BAD_REQUEST" }
      },
      isAxiosError: true
    } as unknown as AxiosError;

    const responseInterceptor = (api.interceptors.response as { handlers: { rejected: (error: unknown) => unknown }[] }).handlers[0].rejected;
    
    try {
      await responseInterceptor(error);
    } catch (err) {
      const errorObj = err as { message: string; status: number; code: string };
      expect(errorObj.message).toBe("Bad Request");
      expect(errorObj.status).toBe(400);
      expect(errorObj.code).toBe("ERR_BAD_REQUEST");
    }
  });
});
