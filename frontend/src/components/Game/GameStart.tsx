import React, { useState } from 'react';
import styled from 'styled-components';
import { ApiService, TokenUtils } from '../../services/api';
import { GameVariant } from '../../types/api';
import { useNavigate } from 'react-router-dom';

const Container = styled.div`
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: #e9edf2;
  padding: 20px;
`;

const Header = styled.div`
  position: absolute;
  top: 20px;
  right: 20px;
`;

const LogoutButton = styled.button`
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
  background: #ffffff;
  padding: 3rem;
  border-radius: 15px;
  border: 1px solid #d5dce6;
  box-shadow: 0 8px 24px rgba(31, 41, 55, 0.12);
  text-align: center;
  max-width: 500px;
  width: 100%;
`;

const Title = styled.h1`
  color: #1f2937;
  margin-bottom: 1rem;
  font-size: 2.5rem;
`;

const Subtitle = styled.h2`
  color: #666;
  margin-bottom: 2rem;
  font-weight: 300;
  font-size: 1.2rem;
`;

const Description = styled.p`
  color: #888;
  margin-bottom: 2rem;
  line-height: 1.6;
  font-size: 1rem;
`;

const StartButton = styled.button`
  padding: 1.2rem 2rem;
  background: #1f4b99;
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 1.2rem;
  cursor: pointer;
  transition: transform 0.2s;
  min-width: 200px;

  &:hover {
    transform: translateY(-3px);
    background: #183d80;
  }

  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
    transform: none;
  }
`;

const ErrorMessage = styled.div`
  background: #fee;
  color: #c33;
  padding: 1rem;
  border-radius: 5px;
  border: 1px solid #fcc;
  margin-bottom: 1rem;
`;

const FlagIcon = styled.div`
  color: #4b5563;
  font-size: 1rem;
  font-weight: 600;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  margin-bottom: 1rem;
`;

const VariantBadge = styled.div`
  display: inline-block;
  margin-bottom: 1.5rem;
  padding: 0.4rem 0.75rem;
  border-radius: 999px;
  background: #eef2ff;
  color: #1f4b99;
  font-size: 0.85rem;
  font-weight: 600;
`;

function variantLabel(variant: GameVariant): string {
  return variant === 'multiple_choice' ? 'Режим: выбор из вариантов' : 'Режим: ввод с клавиатуры';
}

const GameStart: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [lastVariant, setLastVariant] = useState<GameVariant | null>(null);
  const navigate = useNavigate();

  const handleStartGame = async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await ApiService.startGame();
      localStorage.setItem('current_game_id', response.game_id);
      localStorage.setItem('current_game_variant', response.variant);
      setLastVariant(response.variant);
      navigate('/play');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Не удалось начать игру');
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    TokenUtils.removeToken();
    localStorage.removeItem('current_game_id');
    navigate('/login');
  };

  return (
    <Container>
      <Header>
        <LogoutButton onClick={handleLogout}>
          Выйти
        </LogoutButton>
      </Header>

      <GameCard>
        <FlagIcon>Тренировка</FlagIcon>
        <Title>Угадай Флаг</Title>
        <Subtitle>Проверь свои знания флагов стран мира!</Subtitle>
        <Description>
          В этой игре тебе предстоит угадывать страны по их флагам.
          Режим игры выбирается случайно: ввод названия с клавиатуры или выбор из четырёх вариантов.
        </Description>

        {lastVariant && <VariantBadge>{variantLabel(lastVariant)}</VariantBadge>}

        {error && <ErrorMessage>{error}</ErrorMessage>}

        <StartButton onClick={handleStartGame} disabled={loading}>
          {loading ? 'Создание игры...' : 'Начать игру'}
        </StartButton>
      </GameCard>
    </Container>
  );
};

export default GameStart; 