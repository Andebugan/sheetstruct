import axios, { AxiosInstance, AxiosError } from 'axios';
import { User, Sheet, Component, TokenResponse, SheetFilter, ApiError, SheetID, ComponentID } from '../types';

const getApiBaseUrl = (): string => {
  // @ts-ignore - process.env is available in CRA
  return process.env.REACT_APP_API_URL || 'http://localhost:5000/api/v1';
};

const API_BASE_URL = getApiBaseUrl();

class ApiService {
  private api: AxiosInstance;

  constructor() {
    this.api = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    this.api.interceptors.request.use((config: any) => {
      const token = localStorage.getItem('token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });

    this.api.interceptors.response.use(
      (response) => response,
      async (error: AxiosError) => {
        const originalRequest = error.config as any;
        if (error.response?.status === 401 && !originalRequest._retry) {
          originalRequest._retry = true;
          const refreshToken = localStorage.getItem('refreshToken');
          if (refreshToken) {
            try {
              const newToken = await this.refreshToken(refreshToken);
              originalRequest.headers.Authorization = `Bearer ${newToken}`;
              return this.api(originalRequest);
            } catch {
              localStorage.removeItem('token');
              localStorage.removeItem('refreshToken');
              window.location.href = '/login';
            }
          }
        }
        return Promise.reject(error);
      }
    );
  }

  async healthCheck(): Promise<boolean> {
    try {
      await this.api.get('/health');
      return true;
    } catch {
      return false;
    }
  }

  async getCurrentUser(): Promise<User> {
    const response = await this.api.get<User>('/user');
    return response.data;
  }

  async createUser(userData: { Name: string; Email: string; Password: string }): Promise<void> {
    await this.api.post('/user', userData);
  }

  async updateUser(userData: User): Promise<User> {
    const response = await this.api.patch<User>('/user', userData);
    return response.data;
  }

  async deleteUser(): Promise<void> {
    await this.api.delete('/user');
  }

  async authenticate(credentials: { Login: string; Password: string }): Promise<TokenResponse> {
    const response = await this.api.post<TokenResponse>('/user/login', credentials);
    if (response.data.Token) {
      localStorage.setItem('token', response.data.Token);
      localStorage.setItem('refreshToken', response.data.RefreshToken);
    }
    return response.data;
  }

  async refreshToken(refreshToken: string): Promise<string> {
    const response = await this.api.post<string>('/user/refresh_token', { refresh_token: refreshToken });
    const newToken = typeof response.data === 'string' ? response.data : (response.data as any).token || (response.data as any).Token;
    if (newToken) {
      localStorage.setItem('token', newToken);
    }
    return newToken;
  }

  async logout(): Promise<void> {
    await this.api.post('/user/logout');
    localStorage.removeItem('token');
    localStorage.removeItem('refreshToken');
  }

  async resetPassword(data: { Name?: string; Email?: string }): Promise<void> {
    await this.api.post('/user/reset_password', data);
  }

  async getSheets(filter?: SheetFilter): Promise<Sheet[]> {
    const response = await this.api.get<Sheet[]>('/sheets');
    return Array.isArray(response.data) ? response.data : [];
  }

  async getSheet(sid: number): Promise<Sheet> {
    const response = await this.api.get<Sheet>(`/sheet/${sid}?sid=${sid}`);
    return response.data;
  }

  async createSheet(): Promise<Sheet> {
    const response = await this.api.post<Sheet>('/sheet');
    return response.data;
  }

  async updateSheet(sheet: Sheet): Promise<Sheet> {
    const response = await this.api.patch<Sheet>('/sheet', sheet);
    return response.data;
  }

  async deleteSheet(sid: number): Promise<void> {
    await this.api.delete(`/sheet/${sid}?sid=${sid}`);
  }

  async cloneSheet(sid: number): Promise<Sheet> {
    const response = await this.api.put<Sheet>(`/sheet/${sid}?sid=${sid}`);
    return response.data;
  }

  async getComponents(sid: number): Promise<Component[]> {
    const response = await this.api.get<Component[]>(`/sheet/${sid}/components?sid=${sid}`);
    const components = Array.isArray(response.data) ? response.data : [];
    return components
  }

  async getComponent(sid: number, cid: number): Promise<Component> {
    const response = await this.api.get<Component>(`/sheet/${sid}/component/${cid}?sid=${sid}&cid=${cid}`);
    return response.data;
  }

  async createComponent(sid: number): Promise<Component> {
    const response = await this.api.post<Component>(`/sheet/${sid}/component?sid=${sid}`);
    return response.data;
  }

  async updateComponent(component: Component): Promise<Component> {
    const sid = typeof component.SId === 'object' ? component.SId.Value : component.SId;
    let varValueBytes: number[] = [];
    if (typeof component.VarValue === 'string') {
      varValueBytes = Array.from(new TextEncoder().encode(component.VarValue));
    } else if (component.VarValue instanceof Uint8Array) {
      varValueBytes = Array.from(component.VarValue);
    } else if (Array.isArray(component.VarValue)) {
      varValueBytes = component.VarValue;
    }
    
    const componentData = {
      ...component,
      VarValue: varValueBytes,
    };
    const response = await this.api.patch<Component>(`/sheet/${sid}/component?sid=${sid}`, componentData);
    return response.data;
  }

  async deleteComponent(sid: number, cid: number): Promise<void> {
    await this.api.delete(`/sheet/${sid}/component/${cid}?sid=${sid}&cid=${cid}`);
  }

  async cloneComponent(sid: number, cid: number): Promise<Component> {
    const response = await this.api.put<Component>(`/sheet/${sid}/component/${cid}?sid=${sid}&cid=${cid}`);
    return response.data;
  }

  handleError(error: unknown): ApiError {
    if (axios.isAxiosError(error)) {
      const axiosError = error as AxiosError<{ message?: string }>;
      return {
        message: axiosError.response?.data?.message || axiosError.message || 'An error occurred',
        status: axiosError.response?.status || 500,
      };
    }
    return {
      message: 'An unexpected error occurred',
      status: 500,
    };
  }
}

export const apiService = new ApiService();

