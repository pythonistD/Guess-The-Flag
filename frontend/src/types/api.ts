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
}

export interface QuestionResponse {
  question_text: string;
  flag_svg: string;
  question_id: string;
}

export interface AnswerRequest {
  gameId: string;
  questionId: string;
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
  flag_svg: string;
}

export type EndGameResponse = QuestionWithAnswer[];

export interface ErrorResponse {
  error: string;
}

export interface FlagDebugItem {
  country_id: number;
  flag_svg: string;
}
