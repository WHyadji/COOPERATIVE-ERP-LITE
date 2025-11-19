// ============================================================================
// API Client - Axios Instance with Interceptors
// Handles authentication, error handling, and API communication
// ============================================================================

import axios, {
  AxiosError,
  AxiosInstance,
  InternalAxiosRequestConfig,
} from "axios";
import type { APIErrorResponse } from "@/types";

// ============================================================================
// Constants
// ============================================================================

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080/api/v1";
const TOKEN_KEY = "auth_token";

// ============================================================================
// Token Management
// ============================================================================

/**
 * Store JWT token in memory (more secure than localStorage)
 * In production, consider using httpOnly cookies
 */
let authToken: string | null = null;

export const tokenManager = {
  getToken: (): string | null => {
    if (typeof window !== "undefined" && !authToken) {
      // Fallback to localStorage for page refreshes
      authToken = localStorage.getItem(TOKEN_KEY);
    }
    return authToken;
  },

  setToken: (token: string): void => {
    authToken = token;
    if (typeof window !== "undefined") {
      localStorage.setItem(TOKEN_KEY, token);
    }
  },

  removeToken: (): void => {
    authToken = null;
    if (typeof window !== "undefined") {
      localStorage.removeItem(TOKEN_KEY);
    }
  },
};

// ============================================================================
// Axios Instance Creation
// ============================================================================

const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
  timeout: 30000, // 30 seconds
  withCredentials: false, // Set to true if using httpOnly cookies
});

// ============================================================================
// Request Interceptor - Add Authorization Header
// ============================================================================

apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = tokenManager.getToken();

    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    // Log request in development
    if (process.env.NODE_ENV === "development") {
      console.log(
        `[API Request] ${config.method?.toUpperCase()} ${config.url}`,
        {
          data: config.data,
          params: config.params,
        }
      );
    }

    return config;
  },
  (error) => {
    console.error("[API Request Error]", error);
    return Promise.reject(error);
  }
);

// ============================================================================
// Response Interceptor - Handle Errors & Token Expiration
// ============================================================================

apiClient.interceptors.response.use(
  (response) => {
    // Log response in development
    if (process.env.NODE_ENV === "development") {
      console.log(`[API Response] ${response.config.url}`, response.data);
    }
    return response;
  },
  (error: AxiosError<APIErrorResponse>) => {
    // Handle different error scenarios
    if (error.response) {
      const { status, data } = error.response;

      // Handle specific status codes
      switch (status) {
        case 401:
          // Unauthorized - token expired or invalid
          console.error("[API Error] Unauthorized - clearing token");
          tokenManager.removeToken();

          // Redirect to login if not already there
          if (
            typeof window !== "undefined" &&
            !window.location.pathname.includes("/login")
          ) {
            window.location.href = "/login";
          }
          break;

        case 403:
          // Forbidden - insufficient permissions
          console.error("[API Error] Forbidden - insufficient permissions");
          break;

        case 404:
          // Not found
          console.error("[API Error] Resource not found");
          break;

        case 422:
          // Validation error
          console.error("[API Error] Validation error", data?.error);
          break;

        case 500:
        case 502:
        case 503:
          // Server errors
          console.error("[API Error] Server error", status);
          break;

        default:
          console.error("[API Error]", status, data);
      }

      // Return structured error
      return Promise.reject({
        status,
        message: data?.message || "An error occurred",
        error: data?.error || {
          code: `HTTP_${status}`,
          message: error.message,
        },
      });
    } else if (error.request) {
      // Request made but no response received
      console.error("[API Error] No response received", error.request);
      return Promise.reject({
        status: 0,
        message: "Network error - no response from server",
        error: { code: "NETWORK_ERROR", message: error.message },
      });
    } else {
      // Error in request configuration
      console.error("[API Error] Request configuration error", error.message);
      return Promise.reject({
        status: 0,
        message: "Request configuration error",
        error: { code: "CONFIG_ERROR", message: error.message },
      });
    }
  }
);

// ============================================================================
// Export API Client
// ============================================================================

export default apiClient;

// ============================================================================
// Utility Functions
// ============================================================================

/**
 * Check if user is authenticated
 */
export const isAuthenticated = (): boolean => {
  return !!tokenManager.getToken();
};

/**
 * Get API base URL (useful for constructing URLs)
 */
export const getAPIBaseURL = (): string => {
  return API_BASE_URL;
};
