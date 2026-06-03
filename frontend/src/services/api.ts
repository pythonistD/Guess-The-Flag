import axios, { AxiosResponse } from 'axios';
import {
  RegisterRequest,
  LoginRequest,
  TokenResponse,
  StartGameResponse,
  QuestionRequest,
  QuestionResponse,
  AnswerRequest,
  AnswerResponse,
  EndGameRequest,
  EndGameResponse,
  FlagDebugItem,
} from '../types/api';

// Empty string is valid (same-origin via nginx proxy in Docker/Traefik).
const API_BASE_URL =
  process.env.REACT_APP_API_BASE_URL !== undefined
    ? process.env.REACT_APP_API_BASE_URL
    : 'http://localhost:8080';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use((config) => {
  console.log(`API Request: ${config.method?.toUpperCase()} ${config.url}`, config.data);
  const token = localStorage.getItem('jwt_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => {
    console.log(`API Response: ${response.status}`, response.config.url);
    return response;
  },
  (error) => {
    console.error(`API Error: ${error.response?.status}`, error.response?.data || error.message);
    if (error.response?.status === 401) {
      localStorage.removeItem('jwt_token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export class ApiService {
  static async register(data: RegisterRequest): Promise<TokenResponse> {
    const response: AxiosResponse<TokenResponse> = await api.post('/auth/register', data);
    return response.data;
  }

  static async login(data: LoginRequest): Promise<TokenResponse> {
    const response: AxiosResponse<TokenResponse> = await api.post('/auth/login', data);
    return response.data;
  }

  static async startGame(langCode: string = 'rus'): Promise<StartGameResponse> {
    const response: AxiosResponse<StartGameResponse> = await api.post(
      `/game/start?lang_code=${encodeURIComponent(langCode)}`
    );
    return response.data;
  }

  static async getQuestion(data: QuestionRequest): Promise<QuestionResponse> {
    const response: AxiosResponse<QuestionResponse> = await api.post(
      `/game/${data.gameId}/questions/next`
    );
    return response.data;
  }

  static async answerQuestion(data: AnswerRequest): Promise<AnswerResponse> {
    const body: Record<string, unknown> = {};
    if (data.skipped) {
      body.skipped = true;
    } else if (data.selected_country !== undefined) {
      body.selected_country = data.selected_country;
      body.skipped = false;
    } else if (data.answer !== undefined) {
      body.answer = data.answer;
      body.skipped = false;
    }

    const response: AxiosResponse<AnswerResponse> = await api.post(
      `/game/${data.gameId}/questions/${data.questionId}/answer`,
      body
    );
    return response.data;
  }

  static async endGame(data: EndGameRequest): Promise<EndGameResponse> {
    const response: AxiosResponse<EndGameResponse> = await api.post(`/game/${data.gameId}/end`);
    return response.data;
  }

  static async getAllFlags(): Promise<FlagDebugItem[]> {
    const response: AxiosResponse<FlagDebugItem[]> = await api.get('/debug/flags');
    return response.data;
  }
}

export const TokenUtils = {
  setToken: (token: string) => {
    localStorage.setItem('jwt_token', token);
  },

  getToken: (): string | null => {
    return localStorage.getItem('jwt_token');
  },

  removeToken: () => {
    localStorage.removeItem('jwt_token');
  },

  isAuthenticated: (): boolean => {
    return !!localStorage.getItem('jwt_token');
  },
};
