const API_BASE = '/api';

class ApiClient {
  private getToken(): string | null {
    return localStorage.getItem('access_token');
  }

  private getHeaders(authenticated = false): Record<string, string> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };
    if (authenticated) {
      const token = this.getToken();
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }
    }
    return headers;
  }

  async get<T>(path: string, auth = false): Promise<T> {
    const res = await fetch(`${API_BASE}${path}`, {
      headers: this.getHeaders(auth),
    });
    if (!res.ok) {
      const err = await res.json().catch(() => ({ error: 'Request failed' }));
      throw new Error(err.error || 'Request failed');
    }
    return res.json();
  }

  async post<T>(path: string, body?: unknown, auth = false): Promise<T> {
    const res = await fetch(`${API_BASE}${path}`, {
      method: 'POST',
      headers: this.getHeaders(auth),
      body: body ? JSON.stringify(body) : undefined,
    });
    if (!res.ok) {
      const err = await res.json().catch(() => ({ error: 'Request failed' }));
      throw new Error(err.error || 'Request failed');
    }
    if (res.status === 204) return undefined as T;
    return res.json();
  }

  async put<T>(path: string, body?: unknown, auth = true): Promise<T> {
    const res = await fetch(`${API_BASE}${path}`, {
      method: 'PUT',
      headers: this.getHeaders(auth),
      body: body ? JSON.stringify(body) : undefined,
    });
    if (!res.ok) {
      const err = await res.json().catch(() => ({ error: 'Request failed' }));
      throw new Error(err.error || 'Request failed');
    }
    return res.json();
  }

  async delete(path: string, auth = true): Promise<void> {
    const res = await fetch(`${API_BASE}${path}`, {
      method: 'DELETE',
      headers: this.getHeaders(auth),
    });
    if (!res.ok) {
      const err = await res.json().catch(() => ({ error: 'Request failed' }));
      throw new Error(err.error || 'Request failed');
    }
  }

  async upload(file: File): Promise<{ id: string; file_name: string; file_path: string; mime_type: string; size: string }> {
    const formData = new FormData();
    formData.append('file', file);
    const token = this.getToken();
    const headers: Record<string, string> = {};
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }
    const res = await fetch(`${API_BASE}/upload`, {
      method: 'POST',
      headers,
      body: formData,
    });
    if (!res.ok) {
      const err = await res.json().catch(() => ({ error: 'Upload failed' }));
      throw new Error(err.error || 'Upload failed');
    }
    return res.json();
  }

  async refreshToken(): Promise<boolean> {
    const refreshToken = localStorage.getItem('refresh_token');
    if (!refreshToken) return false;
    try {
      const data = await this.post<{ access_token: string; refresh_token: string; expires_in: number }>(
        '/auth/refresh',
        { refresh_token: refreshToken },
        false
      );
      localStorage.setItem('access_token', data.access_token);
      localStorage.setItem('refresh_token', data.refresh_token);
      return true;
    } catch {
      this.logout();
      return false;
    }
  }

  logout(): void {
    const refreshToken = localStorage.getItem('refresh_token');
    if (refreshToken) {
      this.post('/auth/logout', { refresh_token: refreshToken }, false).catch(() => {});
    }
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('user');
  }
}

export const api = new ApiClient();
