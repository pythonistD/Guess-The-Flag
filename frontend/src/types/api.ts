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
  variant: GameVariant;
  max_flags: number;
}

export type GameVariant = 'text_input' | 'multiple_choice';

export interface QuestionRequest {
  gameId: string;
}

export interface AnswerOption {
  country_id: number;
  name: string;
}

export interface QuestionResponse {
  question_text: string;
  flag_svg: string;
  question_id: string;
  variant: GameVariant;
  options?: AnswerOption[];
  step: number;
  max_flags: number;
  steps_completed: number;
}

export interface AnswerRequest {
  gameId: string;
  questionId: string;
  answer?: string;
  selected_country?: number;
  skipped?: boolean;
}

export interface AnswerResponse {
  is_correct: boolean;
  step: number;
  max_flags: number;
  steps_completed: number;
  game_completed: boolean;
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
