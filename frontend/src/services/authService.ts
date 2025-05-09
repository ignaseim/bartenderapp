import { LoginResponse, User } from '../types/models';
import { authApi, apiRequest } from './apiClient';

// Auth service methods
export const authService = {
  // Login with username and password
  async login(username: string, password: string): Promise<LoginResponse> {
    return apiRequest<LoginResponse>(authApi, {
      method: 'POST',
      url: '/login',
      data: { username, password },
    });
  },

  // Get current user information
  async getCurrentUser(token?: string): Promise<User> {
    const config: any = {
      method: 'GET',
      url: '/users/me',
    };

    // If token provided directly, use it instead of from localStorage
    if (token) {
      config.headers = { Authorization: `Bearer ${token}` };
    }

    return apiRequest<User>(authApi, config);
  },

  // Refresh the access token
  async refreshToken(refreshToken: string): Promise<LoginResponse> {
    return apiRequest<LoginResponse>(authApi, {
      method: 'POST',
      url: '/refresh',
      data: { refresh_token: refreshToken },
    });
  },

  // Verify if a token is valid
  async verifyToken(token: string): Promise<boolean> {
    try {
      await apiRequest<{ valid: boolean }>(authApi, {
        method: 'POST',
        url: '/verify',
        data: { token },
      });
      return true;
    } catch (error) {
      return false;
    }
  },

  // Get a list of all users (admin only)
  async getUsers(): Promise<User[]> {
    return apiRequest<User[]>(authApi, {
      method: 'GET',
      url: '/users',
    });
  },

  // Get a specific user by ID
  async getUserById(id: number): Promise<User> {
    return apiRequest<User>(authApi, {
      method: 'GET',
      url: `/users/${id}`,
    });
  },

  // Create a new user (admin only)
  async createUser(userData: Partial<User> & { password: string }): Promise<User> {
    return apiRequest<User>(authApi, {
      method: 'POST',
      url: '/users',
      data: userData,
    });
  },

  // Update a user
  async updateUser(id: number, userData: Partial<User> & { password?: string }): Promise<User> {
    return apiRequest<User>(authApi, {
      method: 'PUT',
      url: `/users/${id}`,
      data: userData,
    });
  },

  // Delete a user (admin only)
  async deleteUser(id: number): Promise<void> {
    return apiRequest<void>(authApi, {
      method: 'DELETE',
      url: `/users/${id}`,
    });
  },
}; 