import React, { useState, useEffect, useMemo } from 'react';
import styled from 'styled-components';
import { useNavigate } from 'react-router-dom';
import { QuestionWithAnswer } from '../../types/api';

const FLAG_MAX_WIDTH = 120;
const FLAG_MAX_HEIGHT = 80;

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

function normalizeFlagSvg(raw: string, idPrefix: string): string {
  if (!raw) return '';
  const withUniqueIds = uniquifySvgIds(raw, idPrefix);
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
}

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

const FlagImage = styled.div`
  width: ${FLAG_MAX_WIDTH}px;
  margin-right: 1rem;
  flex-shrink: 0;
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

const ResultFlag: React.FC<{ flagSvg: string; idKey: string | number }> = ({ flagSvg, idKey }) => {
  const normalized = useMemo(
    () => normalizeFlagSvg(flagSvg, `r${idKey}-`),
    [flagSvg, idKey]
  );
  return <FlagImage dangerouslySetInnerHTML={{ __html: normalized }} />;
};

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
              <ResultFlag flagSvg={result.flag_svg} idKey={index} />
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