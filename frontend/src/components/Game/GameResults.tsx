import React, { useState, useEffect } from 'react';
import styled from 'styled-components';
import { useNavigate } from 'react-router-dom';
import { QuestionWithAnswer } from '../../types/api';

const Container = styled.div`
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
`;

const Header = styled.div`
  text-align: center;
  margin-bottom: 2rem;
`;

const Title = styled.h1`
  color: white;
  font-size: 2.5rem;
  margin-bottom: 1rem;
`;

const ScoreCard = styled.div`
  background: white;
  padding: 2rem;
  border-radius: 15px;
  text-align: center;
  margin-bottom: 2rem;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
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
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
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

const FlagImage = styled.img`
  width: 60px;
  height: 40px;
  object-fit: cover;
  border-radius: 5px;
  margin-right: 1rem;
  border: 1px solid #e9ecef;
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
  font-size: 1.5rem;
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
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: white;
        &:hover {
          transform: translateY(-2px);
        }
      `
      : `
        background: rgba(255, 255, 255, 0.9);
        color: #333;
        border: 2px solid white;
        &:hover {
          background: white;
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

const GameResults: React.FC = () => {
  const [results, setResults] = useState<QuestionWithAnswer[]>([]);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    const storedResults = localStorage.getItem('game_results');
    if (storedResults) {
      try {
        const parsedResults = JSON.parse(storedResults);
        setResults(parsedResults);
      } catch (error) {
        console.error('Error parsing game results:', error);
      }
    }
    setLoading(false);
  }, []);

  const correctAnswers = results.filter(result => result.is_correct).length;
  const totalQuestions = results.length;
  const scorePercentage = totalQuestions > 0 ? Math.round((correctAnswers / totalQuestions) * 100) : 0;

  const handleNewGame = () => {
    localStorage.removeItem('game_results');
    navigate('/game');
  };

  const handleBackToMenu = () => {
    localStorage.removeItem('game_results');
    navigate('/game');
  };

  if (loading) {
    return (
      <Container>
        <EmptyState>Загрузка результатов...</EmptyState>
      </Container>
    );
  }

  if (results.length === 0) {
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
        <Title>🏆 Результаты игры</Title>
      </Header>

      <ScoreCard>
        <ScoreText>
          {correctAnswers} из {totalQuestions} ({scorePercentage}%)
        </ScoreText>
        <ScoreSubtext>
          {scorePercentage >= 80 && 'Отличный результат! 🎉'}
          {scorePercentage >= 60 && scorePercentage < 80 && 'Хороший результат! 👍'}
          {scorePercentage >= 40 && scorePercentage < 60 && 'Неплохо, но можно лучше! 💪'}
          {scorePercentage < 40 && 'Не расстраивайтесь, попробуйте еще раз! 🔄'}
        </ScoreSubtext>
      </ScoreCard>

      <ResultsContainer>
        <ResultsHeader>
          <ResultsTitle>Подробные результаты</ResultsTitle>
        </ResultsHeader>
        
        <ResultsList>
          {results.map((result, index) => (
            <ResultItem key={index} isCorrect={result.is_correct}>
              <FlagImage 
                src={result.flag_url} 
                alt={`Флаг ${result.name}`}
                onError={(e) => {
                  (e.target as HTMLImageElement).src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNjAiIGhlaWdodD0iNDAiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PHJlY3Qgd2lkdGg9IjYwIiBoZWlnaHQ9IjQwIiBmaWxsPSIjZjhmOWZhIi8+PHRleHQgeD0iNTAlIiB5PSI1MCUiIGRvbWluYW50LWJhc2VsaW5lPSJtaWRkbGUiIHRleHQtYW5jaG9yPSJtaWRkbGUiIGZpbGw9IiM2NzY3NzciIGZvbnQtc2l6ZT0iOHB4Ij7QpNC70LDQszwvdGV4dD48L3N2Zz4=';
                }}
              />
              <ResultInfo>
                <CountryName>{result.name}</CountryName>
                <UserAnswer isCorrect={result.is_correct}>
                  Ваш ответ: {result.answer}
                </UserAnswer>
              </ResultInfo>
              <ResultIcon isCorrect={result.is_correct}>
                {result.is_correct ? '✅' : '❌'}
              </ResultIcon>
            </ResultItem>
          ))}
        </ResultsList>
      </ResultsContainer>

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