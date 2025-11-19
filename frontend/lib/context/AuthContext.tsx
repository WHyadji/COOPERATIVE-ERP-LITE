// ============================================================================
// Authentication Context - User State Management
// Provides authentication state and methods throughout the application
// ============================================================================

"use client";

import React, {
  createContext,
  useContext,
  useState,
  useEffect,
  ReactNode,
} from "react";
import { useRouter } from "next/navigation";
import apiClient, { tokenManager } from "@/lib/api/client";
import type { User, LoginRequest, LoginResponse, APIResponse } from "@/types";

// ============================================================================
// Context Types
// ============================================================================

interface AuthContextType {
  user: User | null;
  loading: boolean;
  isAuthenticated: boolean;
  login: (credentials: LoginRequest) => Promise<void>;
  logout: () => void;
  refreshUser: () => Promise<void>;
}

// ============================================================================
// Create Context
// ============================================================================

const AuthContext = createContext<AuthContextType | undefined>(undefined);

// ============================================================================
// Auth Provider Component
// ============================================================================

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const router = useRouter();

  const isAuthenticated = !!user;

  // ============================================================================
  // Initialize: Check if user is already logged in
  // ============================================================================

  useEffect(() => {
    const initAuth = async () => {
      const token = tokenManager.getToken();

      if (token) {
        try {
          // Fetch user profile with existing token
          await refreshUser();
        } catch (error) {
          console.error("Failed to initialize auth:", error);
          tokenManager.removeToken();
          setUser(null);
        }
      }

      setLoading(false);
    };

    initAuth();
  }, []);

  // ============================================================================
  // Login Function
  // ============================================================================

  const login = async (credentials: LoginRequest): Promise<void> => {
    try {
      setLoading(true);

      const response = await apiClient.post<APIResponse<LoginResponse>>(
        "/auth/login",
        credentials
      );

      if (response.data.success && response.data.data) {
        const { token, pengguna } = response.data.data;

        // Store token
        tokenManager.setToken(token);

        // Set user state
        setUser(pengguna);

        // Redirect to dashboard
        router.push("/dashboard");
      } else {
        throw new Error("Login failed: Invalid response");
      }
    } catch (error: unknown) {
      console.error("Login error:", error);
      tokenManager.removeToken();
      setUser(null);
      throw error;
    } finally {
      setLoading(false);
    }
  };

  // ============================================================================
  // Logout Function
  // ============================================================================

  const logout = (): void => {
    try {
      // Optional: Call logout endpoint
      apiClient.post("/auth/logout").catch((err) => {
        console.warn("Logout endpoint failed:", err);
      });
    } catch (error) {
      console.error("Logout error:", error);
    } finally {
      // Clear token and user state
      tokenManager.removeToken();
      setUser(null);

      // Redirect to login
      router.push("/login");
    }
  };

  // ============================================================================
  // Refresh User Profile
  // ============================================================================

  const refreshUser = async (): Promise<void> => {
    try {
      const response = await apiClient.get<APIResponse<User>>("/auth/profile");

      if (response.data.success && response.data.data) {
        setUser(response.data.data);
      } else {
        throw new Error("Failed to fetch user profile");
      }
    } catch (error) {
      console.error("Refresh user error:", error);
      tokenManager.removeToken();
      setUser(null);
      throw error;
    }
  };

  // ============================================================================
  // Context Value
  // ============================================================================

  const value: AuthContextType = {
    user,
    loading,
    isAuthenticated,
    login,
    logout,
    refreshUser,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

// ============================================================================
// useAuth Hook - Access Auth Context
// ============================================================================

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);

  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }

  return context;
};

// ============================================================================
// Protected Route HOC (Higher-Order Component)
// ============================================================================

interface ProtectedRouteProps {
  children: ReactNode;
  requiredRoles?: string[];
}

export const ProtectedRoute: React.FC<ProtectedRouteProps> = ({
  children,
  requiredRoles,
}) => {
  const { user, loading, isAuthenticated } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading) {
      if (!isAuthenticated) {
        // Not authenticated - redirect to login
        router.push("/login");
      } else if (requiredRoles && user) {
        // Check if user has required role
        const hasRequiredRole = requiredRoles.includes(user.peran);

        if (!hasRequiredRole) {
          // User doesn't have required role - redirect to dashboard or show error
          console.error("Insufficient permissions");
          router.push("/dashboard");
        }
      }
    }
  }, [loading, isAuthenticated, user, requiredRoles, router]);

  // Show loading state
  if (loading) {
    return (
      <div
        style={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          height: "100vh",
        }}
      >
        Loading...
      </div>
    );
  }

  // Show nothing if not authenticated (will redirect)
  if (!isAuthenticated) {
    return null;
  }

  // Check role authorization
  if (requiredRoles && user) {
    const hasRequiredRole = requiredRoles.includes(user.peran);
    if (!hasRequiredRole) {
      return null;
    }
  }

  return <>{children}</>;
};
