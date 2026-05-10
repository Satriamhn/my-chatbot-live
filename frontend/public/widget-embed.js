(function () {
  const ROOT_ID = "mychatbot-widget-embed-root";
  const IFRAME_ID = "mychatbot-widget-embed-frame";
  const DEFAULT_WIDTH = "380px";
  const DEFAULT_HEIGHT = "600px";
  const DEFAULT_TITLE = "Customer support widget";

  const state = {
    baseUrl: null,
    pendingOptions: null,
  };

  function normalizeSize(value, fallback) {
    if (value === undefined || value === null || value === "") {
      return fallback;
    }

    if (typeof value === "number" && Number.isFinite(value)) {
      return `${value}px`;
    }

    const raw = String(value).trim();
    if (raw === "") {
      return fallback;
    }

    return /^\d+(\.\d+)?$/.test(raw) ? `${raw}px` : raw;
  }

  function resolveBaseUrl(candidate) {
    if (candidate) {
      return candidate;
    }

    if (state.baseUrl) {
      return state.baseUrl;
    }

    const currentScript = document.currentScript;
    if (currentScript && currentScript.src) {
      return new URL(".", currentScript.src).href;
    }

    const scripts = document.getElementsByTagName("script");
    const lastScript = scripts[scripts.length - 1];
    if (lastScript && lastScript.src) {
      return new URL(".", lastScript.src).href;
    }

    return window.location.origin + "/";
  }

  function buildWidgetUrl(orgId, baseUrl) {
    const url = new URL("/widget", baseUrl);
    url.searchParams.set("org_id", orgId);
    return url.toString();
  }

  function getOrCreateRoot(position) {
    let root = document.getElementById(ROOT_ID);

    if (!root) {
      root = document.createElement("div");
      root.id = ROOT_ID;
      root.style.position = "fixed";
      root.style.bottom = "24px";
      root.style.zIndex = "2147483647";
      root.style.pointerEvents = "none";
      document.body.appendChild(root);
    }

    const isLeft = position === "bottom-left";
    root.style.left = isLeft ? "24px" : "auto";
    root.style.right = isLeft ? "auto" : "24px";

    return root;
  }

  function getOrCreateFrame(root) {
    let frame = root.querySelector(`#${IFRAME_ID}`);

    if (!frame) {
      frame = document.createElement("iframe");
      frame.id = IFRAME_ID;
      frame.setAttribute("allow", "clipboard-write");
      frame.setAttribute("title", DEFAULT_TITLE);
      frame.style.display = "block";
      frame.style.border = "0";
      frame.style.background = "#fff";
      frame.style.borderRadius = "16px";
      frame.style.boxShadow = "0 24px 64px rgba(15, 23, 42, 0.18)";
      frame.style.overflow = "hidden";
      frame.style.pointerEvents = "auto";
      frame.style.maxWidth = "calc(100vw - 24px)";
      frame.style.maxHeight = "calc(100vh - 48px)";
      root.appendChild(frame);
    }

    return frame;
  }

  function mount(options) {
    if (!document.body) {
      state.pendingOptions = options;
      document.addEventListener(
        "DOMContentLoaded",
        function onReady() {
          document.removeEventListener("DOMContentLoaded", onReady);
          mount(state.pendingOptions || options);
        },
        { once: true },
      );
      return;
    }

    const orgId = options.orgId || options.org_id;
    if (!orgId) {
      console.warn("[widget-embed] Missing org_id.");
      return;
    }

    state.baseUrl = resolveBaseUrl(options.baseUrl);

    const root = getOrCreateRoot(options.position || "bottom-right");
    const frame = getOrCreateFrame(root);

    frame.src = buildWidgetUrl(orgId, state.baseUrl);
    frame.title = options.title || DEFAULT_TITLE;
    frame.style.width = `min(${normalizeSize(options.width, DEFAULT_WIDTH)}, calc(100vw - 24px))`;
    frame.style.height = `min(${normalizeSize(options.height, DEFAULT_HEIGHT)}, calc(100vh - 48px))`;
  }

  function readCurrentScriptOptions() {
    const script = document.currentScript;
    if (!script) {
      return null;
    }

    return {
      org_id: script.getAttribute("data-org-id") || "",
      position: script.getAttribute("data-position") || "bottom-right",
      width: script.getAttribute("data-width") || "",
      height: script.getAttribute("data-height") || "",
      baseUrl: script.getAttribute("data-base-url") || "",
      title: script.getAttribute("data-title") || "",
    };
  }

  window.ChatbotWidgetEmbed = {
    init: mount,
  };

  const autoOptions = readCurrentScriptOptions();
  if (autoOptions) {
    mount(autoOptions);
  }
})();
