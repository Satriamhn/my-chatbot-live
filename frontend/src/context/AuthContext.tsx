import React, { createContext, useContext, useState, useEffect, useCallback } from "react";
import api from "../services/api";

// ─── Types ─────────────────────────────────────────────────────────────────

export interface AuthUser {
  user_id: string;
  email: string;
  role: string;
  org_id: string;
}

interface AuthContextValue {
  user: AuthUser | null;
  token: string | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string, orgName: string) => Promise<void>;
  logout: () => void;
}

// ─── Context ────────────────────────────────────────────────────────────────

const AuthContext = createContext<AuthContextValue | null>(null);

// ─── Provider ───────────────────────────────────────────────────────────────

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [token, setToken] = useState<string | null>(() => localStorage.getItem("token"));
  const [isLoading, setIsLoading] = useState(true);

  // Hydrate user from stored token on mount
  useEffect(() => {
    const storedToken = localStorage.getItem("token");
    const storedOrgId = localStorage.getItem("orgId");

    if (storedToken && storedOrgId) {
      // Fetch fresh user info - headers are handled by interceptors in api.ts
      api
        .get("/api/v1/me")
        .then((res) => {
          setUser(res.data);
          setToken(storedToken);
        })
        .catch(() => {
          // Token expired or invalid — clear storage
          clearAuth();
        })
        .finally(() => setIsLoading(false));
    } else {
      setIsLoading(false);
    }
  }, []);

  const persistAuth = useCallback((tokenStr: string, orgId: string, userData: AuthUser) => {
    localStorage.setItem("token", tokenStr);
    localStorage.setItem("orgId", orgId);
    setToken(tokenStr);
    setUser(userData);
  }, []);

  const clearAuth = useCallback(() => {
    localStorage.removeItem("token");
    localStorage.removeItem("orgId");
    setToken(null);
    setUser(null);
  }, []);

  const login = useCallback(async (email: string, password: string) => {
    const res = await api.post("/api/v1/auth/login", { email, password });
    const { token: t, org_id, user_id, role } = res.data;
    const userData: AuthUser = { user_id, email, role, org_id };
    persistAuth(t, org_id, userData);
  }, [persistAuth]);

  const register = useCallback(async (email: string, password: string, orgName: string) => {
    const res = await api.post("/api/v1/auth/register", {
      email,
      password,
      org_name: orgName,
    });
    const { token: t, org_id, user_id, role } = res.data;
    const userData: AuthUser = { user_id, email, role, org_id };
    persistAuth(t, org_id, userData);
  }, [persistAuth]);

  const logout = useCallback(() => {
    clearAuth();
    window.location.href = "/signin";
  }, [clearAuth]);

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        isLoading,
        isAuthenticated: !!token && !!user,
        login,
        register,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

// ─── Hook ───────────────────────────────────────────────────────────────────

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used within <AuthProvider>");
  }
  return ctx;
}
