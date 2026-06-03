import React, { useState, useEffect, useMemo } from 'react';
import styled from 'styled-components';
import { useNavigate } from 'react-router-dom';
import { QuestionWithAnswer } from '../../types/api';

const FLAG_MAX_WIDTH = 120;
const FLAG_MAX_HEIGHT = 80;

function svgToDataUrl(svg: string): string {
  const encoded = encodeURIComponent(svg)
    .replace(/'/g, '%27')
    .replace(/"/g, '%22');
  return `data:image/svg+xml;charset=utf-8,${encoded}`;
}

const Container = styled.div`
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: #e9edf2;
  padding: 20px;
`;

const Header = styled.div`
  text-align: center;
  margin-bottom: 2rem;
`;

const Title = styled.h1`
  color: #1f2937;
  font-size: 2.5rem;
  margin-bottom: 1rem;
`;

const ScoreCard = styled.div`
  background: white;
  padding: 2rem;
  border-radius: 15px;
  text-align: center;
  margin-bottom: 2rem;
  border: 1px solid #d5dce6;
  box-shadow: 0 8px 24px rgba(31, 41, 55, 0.12);
`;

const ScoreText = styled.h2`
  color: #333;
  font-size: 2rem;
  margin-bottom: 0.5rem;
`;

const ScoreSubtext = styled.p`
  color: #666;
  font-size: 1.1rem;
`;

const ResultsContainer = styled.div`
  background: white;
  border-radius: 15px;
  border: 1px solid #d5dce6;
  box-shadow: 0 8px 24px rgba(31, 41, 55, 0.12);
  overflow: hidden;
  max-width: 800px;
  width: 100%;
  margin: 0 auto;
`;

const ResultsHeader = styled.div`
  background: #f8f9fa;
  padding: 1.5rem;
  border-bottom: 1px solid #e9ecef;
`;

const ResultsTitle = styled.h3`
  color: #333;
  margin: 0;
  font-size: 1.5rem;
`;

const ResultsList = styled.div`
  max-height: 400px;
  overflow-y: auto;
`;

const ResultItem = styled.div<{ isCorrect: boolean }>`
  display: flex;
  align-items: center;
  padding: 1rem 1.5rem;
  border-bottom: 1px solid #e9ecef;
  
  ${({ isCorrect }) =>
    isCorrect
      ? `
        background: rgba(212, 237, 218, 0.3);
        border-left: 4px solid #28a745;
      `
      : `
        background: rgba(248, 215, 218, 0.3);
        border-left: 4px solid #dc3545;
      `}

  &:last-child {
    border-bottom: none;
  }
`;

const FlagSlot = styled.div`
  width: ${FLAG_MAX_WIDTH}px;
  height: ${FLAG_MAX_HEIGHT}px;
  margin-right: 1rem;
  flex-shrink: 0;
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
  box-shadow: 0 0 0 1px #c5ccd6;
`;

const ResultInfo = styled.div`
  flex: 1;
`;

const CountryName = styled.div`
  font-weight: bold;
  color: #333;
  margin-bottom: 0.25rem;
`;

const UserAnswer = styled.div<{ isCorrect: boolean }>`
  color: ${({ isCorrect }) => (isCorrect ? '#28a745' : '#dc3545')};
  font-size: 0.9rem;
`;

const ResultIcon = styled.div<{ isCorrect: boolean }>`
  color: ${({ isCorrect }) => (isCorrect ? '#155724' : '#721c24')};
  background: ${({ isCorrect }) => (isCorrect ? '#d4edda' : '#f8d7da')};
  border: 1px solid ${({ isCorrect }) => (isCorrect ? '#c3e6cb' : '#f5c6cb')};
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 700;
  letter-spacing: 0.04em;
  text-transform: uppercase;
  padding: 0.3rem 0.6rem;
  margin-left: 1rem;
`;

const ButtonGroup = styled.div`
  display: flex;
  justify-content: center;
  gap: 1rem;
  margin-top: 2rem;
`;

const Button = styled.button<{ variant?: 'primary' | 'secondary' }>`
  padding: 1rem 2rem;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.3s;
  min-width: 150px;

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
        background: #ffffff;
        color: #333;
        border: 2px solid #d5dce6;
        &:hover {
          background: #f4f6f9;
        }
      `}

  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
    transform: none;
  }
`;

const EmptyState = styled.div`
  text-align: center;
  padding: 3rem 2rem;
  color: #666;
`;

const ResultFlag: React.FC<{ flagSvg: string; idKey: string | number }> = ({ flagSvg }) => {
  const url = useMemo(() => (flagSvg ? svgToDataUrl(flagSvg) : ''), [flagSvg]);
  return (
    <FlagSlot>
      {url && <FlagImage src={url} alt="Флаг страны" />}
    </FlagSlot>
  );
};

const GameResults: React.FC = () => {
  const [results, setResults] = useState<QuestionWithAnswer[]>([]);
  const [hasStoredResults, setHasStoredResults] = useState(false);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    const storedResults = localStorage.getItem('game_results');
    if (storedResults !== null) {
      setHasStoredResults(true);
      try {
        const parsedResults = JSON.parse(storedResults);
        setResults(Array.isArray(parsedResults) ? parsedResults : []);
      } catch (error) {
        console.error('Error parsing game results:', error);
        setResults([]);
      }
    }
    setLoading(false);
  }, []);

  const correctAnswers = results.filter(result => result.is_correct).length;
  const totalQuestions = results.length;
  const scorePercentage = totalQuestions > 0 ? Math.round((correctAnswers / totalQuestions) * 100) : 0;

  const handleNewGame = () => {
    localStorage.removeItem('game_results');
    navigate('/');
  };

  const handleBackToMenu = () => {
    localStorage.removeItem('game_results');
    navigate('/');
  };

  if (loading) {
    return (
      <Container>
        <EmptyState>Загрузка результатов...</EmptyState>
      </Container>
    );
  }

  if (!hasStoredResults) {
    return (
      <Container>
        <EmptyState>
          <h2>Результаты не найдены</h2>
          <p>Похоже, что результаты игры не сохранились.</p>
          <Button onClick={handleBackToMenu}>
            Вернуться к игре
          </Button>
        </EmptyState>
      </Container>
    );
  }

  return (
    <Container>
      <Header>
        <Title>Результаты игры</Title>
      </Header>

      <ScoreCard>
        <ScoreText>
          {correctAnswers} из {totalQuestions} ({scorePercentage}%)
        </ScoreText>
        <ScoreSubtext>
          {totalQuestions === 0 && 'Вы не ответили ни на один вопрос.'}
          {totalQuestions > 0 && scorePercentage >= 80 && 'Отличный результат!'}
          {totalQuestions > 0 && scorePercentage >= 60 && scorePercentage < 80 && 'Хороший результат!'}
          {totalQuestions > 0 && scorePercentage >= 40 && scorePercentage < 60 && 'Неплохо, но можно лучше.'}
          {totalQuestions > 0 && scorePercentage < 40 && 'Попробуйте еще раз.'}
        </ScoreSubtext>
      </ScoreCard>

      {results.length > 0 && (
      <ResultsContainer>
        <ResultsHeader>
          <ResultsTitle>Подробные результаты</ResultsTitle>
        </ResultsHeader>
        
        <ResultsList>
          {results.map((result, index) => (
            <ResultItem key={index} isCorrect={result.is_correct}>
              <ResultFlag flagSvg={result.flag_svg} idKey={index} />
              <ResultInfo>
                <CountryName>{result.name}</CountryName>
                <UserAnswer isCorrect={result.is_correct}>
                  Ваш ответ: {result.answer}
                </UserAnswer>
              </ResultInfo>
              <ResultIcon isCorrect={result.is_correct}>
                {result.is_correct ? 'Верно' : 'Ошибка'}
              </ResultIcon>
            </ResultItem>
          ))}
        </ResultsList>
      </ResultsContainer>
      )}

      <ButtonGroup>
        <Button onClick={handleNewGame}>
          Новая игра
        </Button>
        <Button variant="secondary" onClick={handleBackToMenu}>
          Главное меню
        </Button>
      </ButtonGroup>
    </Container>
  );
};

export default GameResults; 