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
} from '../types/api';

const API_BASE_URL = 'http://localhost:8080';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Interceptor для добавления токена к запросам
api.interceptors.request.use((config) => {
  console.log(`API Request: ${config.method?.toUpperCase()} ${config.url}`, config.data);
  const token = localStorage.getItem('jwt_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
    console.log('Token added to request');
  }
  return config;
});

// Interceptor для обработки ошибок аутентификации
api.interceptors.response.use(
  (response) => {
    console.log(`API Response: ${response.status}`, response.data);
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
  // Аутентификация
  static async register(data: RegisterRequest): Promise<TokenResponse> {
    console.log('Registering user:', data.username);
    const response: AxiosResponse<TokenResponse> = await api.post('/auth/register', data);
    return response.data;
  }

  static async login(data: LoginRequest): Promise<TokenResponse> {
    console.log('Logging in user:', data.username);
    const response: AxiosResponse<TokenResponse> = await api.post('/auth/login', data);
    return response.data;
  }

  // Игровые операции
  static async startGame(): Promise<StartGameResponse> {
    console.log('Starting new game');
    const response: AxiosResponse<StartGameResponse> = await api.post('/game/start');
    return response.data;
  }

  static async getQuestion(data: QuestionRequest): Promise<QuestionResponse> {
    console.log('Getting question:', data.questionNum);
    const response: AxiosResponse<QuestionResponse> = await api.post('/game/question', data);
    return response.data;
  }

  static async answerQuestion(data: AnswerRequest): Promise<AnswerResponse> {
    console.log('Submitting answer:', data.answer);
    const response: AxiosResponse<AnswerResponse> = await api.post('/game/answer', data);
    return response.data;
  }

  static async endGame(data: EndGameRequest): Promise<EndGameResponse> {
    console.log('Ending game:', data.gameId);
    const response: AxiosResponse<EndGameResponse> = await api.post('/game/end', data);
    return response.data;
  }
}

// Утилиты для работы с токеном
export const TokenUtils = {
  setToken: (token: string) => {
    localStorage.setItem('jwt_token', token);
    console.log('Token saved to localStorage');
  },

  getToken: (): string | null => {
    return localStorage.getItem('jwt_token');
  },

  removeToken: () => {
    localStorage.removeItem('jwt_token');
    console.log('Token removed from localStorage');
  },

  isAuthenticated: (): boolean => {
    const hasToken = !!localStorage.getItem('jwt_token');
    console.log('Authentication check:', hasToken);
    return hasToken;
  },
}; 