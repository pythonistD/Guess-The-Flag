import React, { useState, useEffect } from 'react';
import styled from 'styled-components';
import { ApiService } from '../../services/api';
import { useNavigate } from 'react-router-dom';
import { QuestionResponse, AnswerResponse } from '../../types/api';

const Container = styled.div`
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
`;

const Header = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
`;

const QuestionCounter = styled.div`
  color: white;
  font-size: 1.2rem;
  font-weight: bold;
`;

const EndGameButton = styled.button`
  padding: 0.5rem 1rem;
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 5px;
  cursor: pointer;
  transition: all 0.3s;

  &:hover {
    background: rgba(255, 255, 255, 0.3);
  }
`;

const GameCard = styled.div`
  background: white;
  border-radius: 15px;
  box-shadow: 0 15px 35px rgba(0, 0, 0, 0.2);
  overflow: hidden;
  max-width: 600px;
  width: 100%;
  margin: 0 auto;
`;

const FlagContainer = styled.div`
  padding: 2rem;
  text-align: center;
  background: #f8f9fa;
`;

const FlagImage = styled.img`
  max-width: 100%;
  height: 200px;
  object-fit: contain;
  border: 3px solid #e9ecef;
  border-radius: 10px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
`;

const QuestionText = styled.h2`
  color: #333;
  margin-top: 1rem;
  font-size: 1.5rem;
`;

const AnswerSection = styled.div`
  padding: 2rem;
`;

const AnswerInput = styled.input`
  width: 100%;
  padding: 1rem;
  border: 2px solid #e1e5e9;
  border-radius: 8px;
  font-size: 1.1rem;
  margin-bottom: 1rem;
  transition: border-color 0.3s;

  &:focus {
    outline: none;
    border-color: #667eea;
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
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: white;
        &:hover {
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

const ResultMessage = styled.div<{ isCorrect: boolean }>`
  padding: 1rem;
  border-radius: 8px;
  margin-bottom: 1rem;
  font-weight: bold;
  text-align: center;

  ${({ isCorrect }) =>
    isCorrect
      ? `
        background: #d4edda;
        color: #155724;
        border: 1px solid #c3e6cb;
      `
      : `
        background: #f8d7da;
        color: #721c24;
        border: 1px solid #f5c6cb;
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
  const [lastResult, setLastResult] = useState<AnswerResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [gameId, setGameId] = useState<string | null>(null);
  const [wasSkipped, setWasSkipped] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    const storedGameId = localStorage.getItem('current_game_id');
    if (!storedGameId) {
      navigate('/game');
      return;
    }
    setGameId(storedGameId);
    loadQuestion(storedGameId, 1);
  }, [navigate]);

  const loadQuestion = async (gameId: string, questionNum: number) => {
    setLoading(true);
    setError(null);
    setLastResult(null);
    setWasSkipped(false);

    try {
      const response = await ApiService.getQuestion({
        gameId,
        questionNum,
      });
      setCurrentQuestion(response);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Не удалось загрузить вопрос');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmitAnswer = async () => {
    if (!gameId || !answer.trim()) return;

    setLoading(true);
    setError(null);

    try {
      const response = await ApiService.answerQuestion({
        gameId,
        questionNum: questionNumber,
        answer: answer.trim(),
      });
      setLastResult(response);
      setAnswer('');
      setWasSkipped(false);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Не удалось отправить ответ');
    } finally {
      setLoading(false);
    }
  };

  const handleNextQuestion = () => {
    if (!gameId) return;
    
    const nextQuestionNum = questionNumber + 1;
    setQuestionNumber(nextQuestionNum);
    loadQuestion(gameId, nextQuestionNum);
  };

  const handleSkipQuestion = async () => {
    if (!gameId) return;

    setLoading(true);
    setError(null);

    try {
      // Отправляем пустой ответ как "пропуск"
      const response = await ApiService.answerQuestion({
        gameId,
        questionNum: questionNumber,
        answer: '',
      });
      setLastResult(response);
      setAnswer('');
      setWasSkipped(true);
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
      localStorage.setItem('game_results', JSON.stringify(results));
      localStorage.removeItem('current_game_id');
      navigate('/game/results');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Не удалось завершить игру');
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !loading && answer.trim()) {
      if (lastResult) {
        handleNextQuestion();
      } else {
        handleSubmitAnswer();
      }
    }
  };

  if (!gameId) {
    return <LoadingSpinner>Загрузка игры...</LoadingSpinner>;
  }

  return (
    <Container>
      <Header>
        <QuestionCounter>Вопрос {questionNumber}</QuestionCounter>
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
              <FlagImage 
                src={currentQuestion.flag_url} 
                alt="Флаг страны"
                onError={(e) => {
                  (e.target as HTMLImageElement).src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjEwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMjAwIiBoZWlnaHQ9IjEwMCIgZmlsbD0iI2Y4ZjlmYSIvPjx0ZXh0IHg9IjUwJSIgeT0iNTAlIiBkb21pbmFudC1iYXNlbGluZT0ibWlkZGxlIiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBmaWxsPSIjNjc2Nzc3Ij7QpNC70LDQsyDQvdC1INC90LDQudC00LXQvTwvdGV4dD48L3N2Zz4=';
                }}
              />
              <QuestionText>{currentQuestion.question_text}</QuestionText>
            </FlagContainer>

            <AnswerSection>
              {error && <ErrorMessage>{error}</ErrorMessage>}
              
              {lastResult && (
                <ResultMessage isCorrect={lastResult.is_correct}>
                  {wasSkipped ? '⏭️ Вопрос пропущен' : lastResult.is_correct ? '✅ Правильно!' : '❌ Неправильно!'}
                </ResultMessage>
              )}

              {!lastResult && (
                <>
                  <AnswerInput
                    type="text"
                    placeholder="Введите название страны..."
                    value={answer}
                    onChange={(e) => setAnswer(e.target.value)}
                    onKeyPress={handleKeyPress}
                    disabled={loading}
                  />
                  <ButtonGroup>
                    <Button
                      onClick={handleSubmitAnswer}
                      disabled={loading || !answer.trim()}
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
                  </ButtonGroup>
                </>
              )}

              {lastResult && (
                <ButtonGroup>
                  <Button onClick={handleNextQuestion}>
                    Следующий вопрос
                  </Button>
                  <Button variant="secondary" onClick={handleEndGame}>
                    Завершить игру
                  </Button>
                </ButtonGroup>
              )}
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