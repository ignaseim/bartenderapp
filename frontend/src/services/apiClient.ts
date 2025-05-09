import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';

// Define the base URLs from environment variables
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8081';
const AUTH_API_URL = import.meta.env.VITE_AUTH_API_URL || 'http://localhost:8081';
const INVENTORY_API_URL = import.meta.env.VITE_INVENTORY_API_URL || 'http://localhost:8082';
const ORDER_API_URL = import.meta.env.VITE_ORDER_API_URL || 'http://localhost:8083';
const PRICING_API_URL = import.meta.env.VITE_PRICING_API_URL || 'http://localhost:8084';

// Create axios instances for each service
const createAxiosInstance = (baseURL: string): AxiosInstance => {
  const instance = axios.create({
    baseURL,
    headers: {
      'Content-Type': 'application/json',
    },
  });

  // Add request interceptor to add auth token
  instance.interceptors.request.use(
    (config) => {
      const token = localStorage.getItem('token');
      if (token && config.headers) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    },
    (error) => Promise.reject(error)
  );

  // Add response interceptor for error handling
  instance.interceptors.response.use(
    (response) => response,
    async (error) => {
      const originalRequest = error.config;
      
      // If 401 Unauthorized and not a retry
      if (error.response?.status === 401 && !originalRequest._retry) {
        originalRequest._retry = true;
        
        try {
          const refreshToken = localStorage.getItem('refresh_token');
          if (!refreshToken) {
            // No refresh token, user needs to log in again
            localStorage.removeItem('token');
            window.location.href = '/login';
            return Promise.reject(error);
          }
          
          // Try to get a new token
          const response = await axios.post(`${AUTH_API_URL}/refresh`, {
            refresh_token: refreshToken,
          });
          
          // Save the new tokens
          const { token, refresh_token } = response.data;
          localStorage.setItem('token', token);
          localStorage.setItem('refresh_token', refresh_token);
          
          // Update the failed request with new token and retry
          originalRequest.headers.Authorization = `Bearer ${token}`;
          return axios(originalRequest);
        } catch (refreshError) {
          // Refresh failed, user needs to log in again
          localStorage.removeItem('token');
          localStorage.removeItem('refresh_token');
          window.location.href = '/login';
          return Promise.reject(refreshError);
        }
      }
      
      return Promise.reject(error);
    }
  );

  return instance;
};

// Create instances for each service
export const authApi = createAxiosInstance(AUTH_API_URL);
export const inventoryApi = createAxiosInstance(INVENTORY_API_URL);
export const orderApi = createAxiosInstance(ORDER_API_URL);
export const pricingApi = createAxiosInstance(PRICING_API_URL);

// Generic API request function
export const apiRequest = async <T>(
  instance: AxiosInstance,
  config: AxiosRequestConfig
): Promise<T> => {
  try {
    const response: AxiosResponse<T> = await instance(config);
    return response.data;
  } catch (error: any) {
    if (error.response) {
      // Server responded with an error status
      const message = error.response.data?.error || error.response.statusText;
      throw new Error(message);
    } else if (error.request) {
      // Request was made but no response
      throw new Error('No response from server. Please try again later.');
    } else {
      // Something else happened
      throw new Error(error.message || 'An unexpected error occurred');
    }
  }
}; 