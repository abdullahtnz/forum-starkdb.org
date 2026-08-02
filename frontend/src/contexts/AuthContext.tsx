import React, { createContext, useContext, useState, useEffect, useCallback } from 'react';
import { api } from '../api';
import type { User, AuthTokens } from '../types';

interface AuthContextType {
  user: User | null;
  token: string | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<void>;
  signup: (username: string, email: string, password: string, acceptedTerms: boolean) => Promise<void>;
  logout: () => void;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(localStorage.getItem('access_token'));
  const [loading, setLoading] = useState(true);

  const refreshUser = useCallback(async () => {
    if (!token) {
      setLoading(false);
      return;
    }
    try {
      const profile = await api.get<User>('/users/me', true);
      setUser(profile);
    } catch {
      const refreshed = await api.refreshToken();
      if (refreshed) {
        setToken(localStorage.getItem('access_token'));
        try {
          const profile = await api.get<User>('/users/me', true);
          setUser(profile);
        } catch {
          api.logout();
          setUser(null);
          setToken(null);
        }
      } else {
        setUser(null);
        setToken(null);
      }
    } finally {
      setLoading(false);
    }
  }, [token]);

  useEffect(() => {
    if (token) {
      refreshUser();
    } else {
      setLoading(false);
    }
  }, []);

  const login = async (email: string, password: string) => {
    const data: AuthTokens = await api.post('/auth/login', { email, password });
    localStorage.setItem('access_token', data.access_token);
    localStorage.setItem('refresh_token', data.refresh_token);
    setToken(data.access_token);
    const profile = await api.get<User>('/users/me', true);
    setUser(profile);
  };

  const signup = async (username: string, email: string, password: string, acceptedTerms: boolean) => {
    const data: AuthTokens = await api.post('/auth/signup', { username, email, password, accepted_terms: acceptedTerms });
    localStorage.setItem('access_token', data.access_token);
    localStorage.setItem('refresh_token', data.refresh_token);
    setToken(data.access_token);
    const profile = await api.get<User>('/users/me', true);
    setUser(profile);
  };

  const logout = () => {
    api.logout();
    setUser(null);
    setToken(null);
  };

  return (
    <AuthContext.Provider value={{ user, token, loading, login, signup, logout, refreshUser }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
}
