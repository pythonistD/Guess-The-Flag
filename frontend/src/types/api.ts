export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface TokenResponse {
  token: string;
}

export interface StartGameResponse {
  game_id: string;
}

export interface QuestionRequest {
  gameId: string;
  questionNum: number;
}

export interface QuestionResponse {
  question_text: string;
  flag_url: string;
}

export interface AnswerRequest {
  gameId: string;
  questionNum: number;
  answer: string;
}

export interface AnswerResponse {
  is_correct: boolean;
}

export interface EndGameRequest {
  gameId: string;
}

export interface QuestionWithAnswer {
  answer: string;
  is_correct: boolean;
  name: string;
  code: string;
  flag_url: string;
}

export type EndGameResponse = QuestionWithAnswer[];

export interface ErrorResponse {
  error: string;
} 