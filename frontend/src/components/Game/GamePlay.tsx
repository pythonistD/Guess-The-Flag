import React, { useState, useEffect, useRef, useMemo } from 'react';
import styled from 'styled-components';
import { ApiService } from '../../services/api';
import { useNavigate } from 'react-router-dom';
import { QuestionResponse, AnswerResponse, GameVariant } from '../../types/api';

function svgToDataUrl(svg: string): string {
  const encoded = encodeURIComponent(svg)
    .replace(/'/g, '%27')
    .replace(/"/g, '%22');
  return `data:image/svg+xml;charset=utf-8,${encoded}`;
}

function variantLabel(variant: GameVariant): string {
  return variant === 'multiple_choice' ? 'Выбор из вариантов' : 'Ввод с клавиатуры';
}

const Container = styled.div`
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: #e9edf2;
  padding: 20px;
`;

const Header = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  gap: 1rem;
  flex-wrap: wrap;
`;

const HeaderLeft = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
`;

const QuestionCounter = styled.div`
  color: #1f2937;
  font-size: 1.2rem;
  font-weight: bold;
`;

const VariantBadge = styled.div`
  display: inline-block;
  width: fit-content;
  padding: 0.25rem 0.65rem;
  border-radius: 999px;
  background: #eef2ff;
  color: #1f4b99;
  font-size: 0.8rem;
  font-weight: 600;
`;

const EndGameButton = styled.button`
  padding: 0.5rem 1rem;
  background: #ffffff;
  color: #1f2937;
  border: 1px solid #c8d0db;
  border-radius: 5px;
  cursor: pointer;
  transition: all 0.3s;

  &:hover {
    background: #f4f6f9;
  }
`;

const GameCard = styled.div`
  background: white;
  border-radius: 15px;
  border: 1px solid #d5dce6;
  box-shadow: 0 8px 24px rgba(31, 41, 55, 0.12);
  overflow: hidden;
  max-width: 600px;
  width: 100%;
  margin: 0 auto;
`;

const FlagContainer = styled.div`
  padding: 2rem;
  text-align: center;
  background: #d9dee5;
  border-bottom: 1px solid #c5ccd6;
`;

const FLAG_MAX_WIDTH = 420;
const FLAG_MAX_HEIGHT = 260;
const TEXT_ANSWER_MIN_HEIGHT = 96;
const MC_ANSWER_MIN_HEIGHT = 260;

const FlagFrame = styled.div`
  width: 100%;
  max-width: ${FLAG_MAX_WIDTH}px;
  height: ${FLAG_MAX_HEIGHT}px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: center;
`;

const FlagImage = styled.img`
  display: block;
  max-width: 100%;
  max-height: 100%;
  width: auto;
  height: auto;
  object-fit: contain;
  box-shadow: 0 0 0 1px #c5ccd6;
`;

const QuestionText = styled.h2`
  color: #333;
  margin: 0 0 1.25rem 0;
  font-size: 1.5rem;
  line-height: 1.3;
  min-height: 2.6rem;
  text-align: center;
`;

const AnswerSection = styled.div`
  padding: 2rem;
`;

const AnswerContent = styled.div<{ $minHeight: number }>`
  min-height: ${({ $minHeight }) => $minHeight}px;
  margin-bottom: 1rem;
  display: flex;
  flex-direction: column;
  justify-content: center;
`;

const AnswerInput = styled.input`
  width: 100%;
  padding: 1rem;
  border: 2px solid #e1e5e9;
  border-radius: 8px;
  font-size: 1.1rem;
  transition: border-color 0.3s;

  &:focus {
    outline: none;
    border-color: #1f4b99;
  }
`;

const OptionsGrid = styled.div`
  display: grid;
  grid-template-columns: 1fr;
  gap: 0.75rem;
`;

const OptionButton = styled.button<{ selected?: boolean }>`
  width: 100%;
  padding: 1rem;
  border-radius: 8px;
  border: 2px solid ${({ selected }) => (selected ? '#1f4b99' : '#e1e5e9')};
  background: ${({ selected }) => (selected ? '#eef2ff' : '#ffffff')};
  color: #1f2937;
  font-size: 1rem;
  text-align: left;
  cursor: pointer;
  transition: all 0.2s;

  &:hover:not(:disabled) {
    border-color: #1f4b99;
    background: #f8faff;
  }

  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
`;

const ButtonGroup = styled.div`
  display: flex;
  gap: 1rem;
`;

const Button = styled.button<{ variant?: 'primary' | 'secondary' }>`
  flex: 1;
  padding: 1rem;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.3s;

  ${({ variant = 'primary' }) =>
    variant === 'primary'
      ? `
        background: #1f4b99;
        color: white;
        &:hover {
          background: #183d80;
          transform: translateY(-2px);
        }
      `
      : `
        background: #f8f9fa;
        color: #333;
        border: 2px solid #e9ecef;
        &:hover {
          background: #e9ecef;
        }
      `}

  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
    transform: none;
  }
`;

const ResultBanner = styled.div<{ $isCorrect: boolean }>`
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  padding: 1.1rem 1.25rem;
  border-radius: 12px;
  font-size: 1.05rem;
  font-weight: 600;
  letter-spacing: 0.01em;
  width: 100%;

  ${({ $isCorrect }) =>
    $isCorrect
      ? `
        background: #f0fdf4;
        color: #166534;
        border: 1px solid #bbf7d0;
        box-shadow: inset 4px 0 0 #22c55e;
      `
      : `
        background: #fff1f2;
        color: #9f1239;
        border: 1px solid #fecdd3;
        box-shadow: inset 4px 0 0 #f43f5e;
      `}
`;

const ResultIcon = styled.span<{ $isCorrect: boolean }>`
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  flex-shrink: 0;
  font-size: 1.1rem;
  font-weight: 700;

  ${({ $isCorrect }) =>
    $isCorrect
      ? `
        background: #22c55e;
        color: #ffffff;
      `
      : `
        background: #f43f5e;
        color: #ffffff;
      `}
`;

const LoadingSpinner = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 2rem;
  font-size: 1.1rem;
  color: #666;
`;

const ErrorMessage = styled.div`
  background: #fee;
  color: #c33;
  padding: 1rem;
  border-radius: 5px;
  border: 1px solid #fcc;
  margin-bottom: 1rem;
`;

const GamePlay: React.FC = () => {
  const [currentQuestion, setCurrentQuestion] = useState<QuestionResponse | null>(null);
  const [questionNumber, setQuestionNumber] = useState(1);
  const [answer, setAnswer] = useState('');
  const [selectedCountryId, setSelectedCountryId] = useState<number | null>(null);
  const [lastResult, setLastResult] = useState<AnswerResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [gameId, setGameId] = useState<string | null>(null);
  const [gameVariant, setGameVariant] = useState<GameVariant>('text_input');
  const navigate = useNavigate();
  const didInitRef = useRef(false);

  const isMultipleChoice = gameVariant === 'multiple_choice';

  const flagSvgUrl = useMemo(() => {
    if (!currentQuestion?.flag_svg) return '';
    return svgToDataUrl(currentQuestion.flag_svg);
  }, [currentQuestion?.flag_svg]);

  useEffect(() => {
    if (didInitRef.current) return;
    didInitRef.current = true;
    const storedGameId = localStorage.getItem('current_game_id');
    if (!storedGameId) {
      navigate('/');
      return;
    }
    const storedVariant = localStorage.getItem('current_game_variant') as GameVariant | null;
    if (storedVariant === 'text_input' || storedVariant === 'multiple_choice') {
      setGameVariant(storedVariant);
    }
    setGameId(storedGameId);
    loadQuestion(storedGameId);
  }, [navigate]);

  const loadQuestion = async (id: string) => {
    setLoading(true);
    setError(null);
    setLastResult(null);
    setAnswer('');
    setSelectedCountryId(null);

    try {
      const response = await ApiService.getQuestion({ gameId: id });
      setCurrentQuestion(response);
      setGameVariant(response.variant);
      localStorage.setItem('current_game_variant', response.variant);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Не удалось загрузить вопрос');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmitAnswer = async () => {
    if (!gameId || !currentQuestion) return;
    if (isMultipleChoice && selectedCountryId === null) return;
    if (!isMultipleChoice && !answer.trim()) return;

    setLoading(true);
    setError(null);

    try {
      const response = await ApiService.answerQuestion(
        isMultipleChoice
          ? {
              gameId,
              questionId: currentQuestion.question_id,
              selected_country: selectedCountryId!,
            }
          : {
              gameId,
              questionId: currentQuestion.question_id,
              answer: answer.trim(),
            }
      );
      setLastResult(response);
      setAnswer('');
      setSelectedCountryId(null);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Не удалось отправить ответ');
    } finally {
      setLoading(false);
    }
  };

  const handleNextQuestion = () => {
    if (!gameId) return;
    setQuestionNumber((n) => n + 1);
    loadQuestion(gameId);
  };

  const handleSkipQuestion = async () => {
    if (!gameId || !currentQuestion || loading) return;

    setLoading(true);
    setError(null);

    try {
      await ApiService.answerQuestion({
        gameId,
        questionId: currentQuestion.question_id,
        skipped: true,
      });
      setQuestionNumber((n) => n + 1);
      const response = await ApiService.getQuestion({ gameId });
      setCurrentQuestion(response);
      setGameVariant(response.variant);
      localStorage.setItem('current_game_variant', response.variant);
      setAnswer('');
      setSelectedCountryId(null);
      setLastResult(null);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Не удалось пропустить вопрос');
    } finally {
      setLoading(false);
    }
  };

  const handleEndGame = async () => {
    if (!gameId) return;

    try {
      const results = await ApiService.endGame({ gameId });
      localStorage.setItem('game_results', JSON.stringify(Array.isArray(results) ? results : []));
      localStorage.removeItem('current_game_id');
      localStorage.removeItem('current_game_variant');
      navigate('/results');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Не удалось завершить игру');
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (isMultipleChoice) return;
    if (e.key === 'Enter' && !loading && answer.trim()) {
      if (lastResult) {
        handleNextQuestion();
      } else {
        handleSubmitAnswer();
      }
    }
  };

  const canSubmit = isMultipleChoice
    ? selectedCountryId !== null
    : answer.trim().length > 0;

  const answerContentMinHeight = isMultipleChoice ? MC_ANSWER_MIN_HEIGHT : TEXT_ANSWER_MIN_HEIGHT;

  if (!gameId) {
    return <LoadingSpinner>Загрузка игры...</LoadingSpinner>;
  }

  return (
    <Container>
      <Header>
        <HeaderLeft>
          <QuestionCounter>Вопрос {questionNumber}</QuestionCounter>
          <VariantBadge>{variantLabel(gameVariant)}</VariantBadge>
        </HeaderLeft>
        <EndGameButton onClick={handleEndGame}>
          Завершить игру
        </EndGameButton>
      </Header>

      <GameCard>
        {loading && !currentQuestion ? (
          <LoadingSpinner>Загрузка вопроса...</LoadingSpinner>
        ) : currentQuestion ? (
          <>
            <FlagContainer>
              <QuestionText>{currentQuestion.question_text}</QuestionText>
              <FlagFrame>
                {flagSvgUrl && <FlagImage src={flagSvgUrl} alt="Флаг страны" />}
              </FlagFrame>
            </FlagContainer>

            <AnswerSection>
              {error && <ErrorMessage>{error}</ErrorMessage>}

              <AnswerContent $minHeight={answerContentMinHeight}>
                {lastResult ? (
                  <ResultBanner $isCorrect={lastResult.is_correct}>
                    <ResultIcon $isCorrect={lastResult.is_correct}>
                      {lastResult.is_correct ? '✓' : '✗'}
                    </ResultIcon>
                    {lastResult.is_correct ? 'Правильно!' : 'Неправильно'}
                  </ResultBanner>
                ) : isMultipleChoice ? (
                  <OptionsGrid>
                    {(currentQuestion.options ?? []).map((option) => (
                      <OptionButton
                        key={option.country_id}
                        type="button"
                        selected={selectedCountryId === option.country_id}
                        onClick={() => setSelectedCountryId(option.country_id)}
                        disabled={loading}
                      >
                        {option.name}
                      </OptionButton>
                    ))}
                  </OptionsGrid>
                ) : (
                  <AnswerInput
                    type="text"
                    placeholder="Введите название страны..."
                    value={answer}
                    onChange={(e) => setAnswer(e.target.value)}
                    onKeyPress={handleKeyPress}
                    disabled={loading}
                  />
                )}
              </AnswerContent>

              <ButtonGroup>
                {lastResult ? (
                  <>
                    <Button onClick={handleNextQuestion}>
                      Следующий вопрос
                    </Button>
                    <Button variant="secondary" onClick={handleEndGame}>
                      Завершить игру
                    </Button>
                  </>
                ) : (
                  <>
                    <Button
                      onClick={handleSubmitAnswer}
                      disabled={loading || !canSubmit}
                    >
                      {loading ? 'Отправка...' : 'Ответить'}
                    </Button>
                    <Button
                      variant="secondary"
                      onClick={handleSkipQuestion}
                      disabled={loading}
                    >
                      Пропустить
                    </Button>
                  </>
                )}
              </ButtonGroup>
            </AnswerSection>
          </>
        ) : (
          <LoadingSpinner>Загрузка...</LoadingSpinner>
        )}
      </GameCard>
    </Container>
  );
};

export default GamePlay;
