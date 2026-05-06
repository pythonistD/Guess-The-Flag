import React, { useState, useEffect, useRef, useMemo } from 'react';
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
  background:rgb(234, 234, 234);
`;

const FLAG_MAX_WIDTH = 420;
const FLAG_MAX_HEIGHT = 260;

const FlagImage = styled.div`
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 0;

  svg {
    display: block;
    max-width: ${FLAG_MAX_WIDTH}px;
    max-height: ${FLAG_MAX_HEIGHT}px;
    width: auto;
    height: auto;
  }
`;

const QuestionText = styled.h2`
  color: #333;
  margin: 0 0 1.25rem 0;
  font-size: 1.5rem;
  text-align: center;
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

function uniquifySvgIds(svg: string, prefix: string): string {
  const ids = new Set<string>();
  const idRegex = /\sid=["']([^"']+)["']/gi;
  let m: RegExpExecArray | null;
  while ((m = idRegex.exec(svg)) !== null) {
    ids.add(m[1]);
  }
  if (ids.size === 0) return svg;

  let result = svg;
  ids.forEach((id) => {
    const newId = `${prefix}${id}`;
    const escId = id.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
    result = result.replace(new RegExp(`(\\sid=["'])${escId}(["'])`, 'g'), `$1${newId}$2`);
    result = result.replace(new RegExp(`(href=["'])#${escId}(["'])`, 'g'), `$1#${newId}$2`);
    result = result.replace(new RegExp(`(xlink:href=["'])#${escId}(["'])`, 'g'), `$1#${newId}$2`);
    result = result.replace(new RegExp(`url\\(#${escId}\\)`, 'g'), `url(#${newId})`);
  });
  return result;
}

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
  const didInitRef = useRef(false);

  const normalizedFlagSvg = useMemo(() => {
    if (!currentQuestion?.flag_svg) return '';
    const idPrefix = `q${currentQuestion.question_id ?? 'x'}-`;
    const withUniqueIds = uniquifySvgIds(currentQuestion.flag_svg, idPrefix);
    return withUniqueIds.replace(/<svg\b([^>]*)>/i, (_match, attrs: string) => {
      const widthMatch = attrs.match(/\swidth=["']([^"']+)["']/i);
      const heightMatch = attrs.match(/\sheight=["']([^"']+)["']/i);
      const viewBoxMatch = attrs.match(/\sviewBox=["']([^"']+)["']/i);

      let newAttrs = attrs
        .replace(/\swidth=["'][^"']*["']/gi, '')
        .replace(/\sheight=["'][^"']*["']/gi, '')
        .replace(/\spreserveAspectRatio=["'][^"']*["']/gi, '');

      let vbW = 0;
      let vbH = 0;
      if (viewBoxMatch) {
        const parts = viewBoxMatch[1].trim().split(/[\s,]+/).map(parseFloat);
        if (parts.length === 4 && parts.every(Number.isFinite)) {
          vbW = parts[2];
          vbH = parts[3];
        }
      } else if (widthMatch && heightMatch) {
        const w = parseFloat(widthMatch[1]);
        const h = parseFloat(heightMatch[1]);
        if (Number.isFinite(w) && Number.isFinite(h) && w > 0 && h > 0) {
          newAttrs += ` viewBox="0 0 ${w} ${h}"`;
          vbW = w;
          vbH = h;
        }
      }

      newAttrs += ' preserveAspectRatio="xMidYMid meet"';

      if (vbW > 0 && vbH > 0) {
        const ratio = vbW / vbH;
        const targetWidth = Math.min(FLAG_MAX_WIDTH, FLAG_MAX_HEIGHT * ratio);
        const targetHeight = targetWidth / ratio;
        return `<svg${newAttrs} width="${targetWidth}" height="${targetHeight}">`;
      }
      return `<svg${newAttrs} width="${FLAG_MAX_WIDTH}" height="${FLAG_MAX_HEIGHT}">`;
    });
  }, [currentQuestion?.flag_svg, currentQuestion?.question_id]);

  useEffect(() => {
    if (didInitRef.current) return;
    didInitRef.current = true;
    const storedGameId = localStorage.getItem('current_game_id');
    if (!storedGameId) {
      navigate('/game');
      return;
    }
    setGameId(storedGameId);
    loadQuestion(storedGameId);
  }, [navigate]);

  const loadQuestion = async (gameId: string) => {
    setLoading(true);
    setError(null);
    setLastResult(null);
    setWasSkipped(false);

    try {
      const response = await ApiService.getQuestion({ gameId });
      setCurrentQuestion(response);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Не удалось загрузить вопрос');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmitAnswer = async () => {
    if (!gameId || !currentQuestion || !answer.trim()) return;

    setLoading(true);
    setError(null);

    try {
      const response = await ApiService.answerQuestion({
        gameId,
        questionId: currentQuestion.question_id,
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
    setQuestionNumber((n) => n + 1);
    loadQuestion(gameId);
  };

  const handleSkipQuestion = async () => {
    if (!gameId || !currentQuestion) return;

    setLoading(true);
    setError(null);

    try {
      const response = await ApiService.answerQuestion({
        gameId,
        questionId: currentQuestion.question_id,
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
              <QuestionText>{currentQuestion.question_text}</QuestionText>
              <FlagImage dangerouslySetInnerHTML={{ __html: normalizedFlagSvg }} />
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