import axios, { AxiosError } from "axios";

export interface ApiError {
  message: string;
  code?: string;
  status?: number;
}

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || "http://localhost:8080",
  headers: {
    "Content-Type": "application/json",
  },
});

/**
 * Request interceptor to inject JWT and X-Org-ID from localStorage.
 */
api.interceptors.request.use(
  (config) => {
    // Inject JWT if available
    const token = localStorage.getItem("token");
    if (token && !config.headers.Authorization) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    // Inject X-Org-ID if available
    const orgId = localStorage.getItem("orgId");
    if (orgId && !config.headers["X-Org-ID"]) {
      config.headers["X-Org-ID"] = orgId;
    }

    return config;
  },
  (error) => Promise.reject(error)
);

/**
 * Response interceptor to handle global errors (401/403).
 */
api.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => {
    const normalizedError: ApiError = {
      message: "An unexpected error occurred",
      status: error.response?.status,
    };

    if (error.response) {
      const data = error.response.data as { message?: string; code?: string };
      normalizedError.message = data?.message || error.message;
      normalizedError.code = data?.code;

      if (error.response.status === 401) {
        // Handle unauthorized access: clear local session and redirect
        localStorage.removeItem("token");
        localStorage.removeItem("orgId");
        
        if (!window.location.pathname.includes("/signin") && !window.location.pathname.includes("/signup")) {
          window.location.href = "/signin";
        }
      } else if (error.response.status === 403) {
        if (data?.code === "ORG_MISMATCH") {
          localStorage.removeItem("orgId");
          normalizedError.message = "Organization access mismatch. Please switch organization.";
        } else {
          normalizedError.message = "You do not have permission to perform this action.";
        }
      }
    } else if (error.request) {
      normalizedError.message = "No response from server. Please check your connection.";
    } else {
      normalizedError.message = error.message;
    }

    return Promise.reject(normalizedError);
  }
);

export default api;
