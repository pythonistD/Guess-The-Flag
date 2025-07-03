import React, { useState } from 'react';
import styled from 'styled-components';
import { ApiService, TokenUtils } from '../../services/api';
import { useNavigate } from 'react-router-dom';

const Container = styled.div`
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
`;

const Header = styled.div`
  position: absolute;
  top: 20px;
  right: 20px;
`;

const LogoutButton = styled.button`
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
  padding: 3rem;
  border-radius: 15px;
  box-shadow: 0 15px 35px rgba(0, 0, 0, 0.2);
  text-align: center;
  max-width: 500px;
  width: 100%;
`;

const Title = styled.h1`
  color: #333;
  margin-bottom: 1rem;
  font-size: 2.5rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
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
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 1.2rem;
  cursor: pointer;
  transition: transform 0.2s;
  min-width: 200px;

  &:hover {
    transform: translateY(-3px);
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
  font-size: 4rem;
  margin-bottom: 1rem;
`;

const GameStart: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  const handleStartGame = async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await ApiService.startGame();
      // Сохраняем ID игры в localStorage для использования в игровых компонентах
      localStorage.setItem('current_game_id', response.game_id);
      navigate('/game/play');
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
        <FlagIcon>🏁</FlagIcon>
        <Title>Угадай Флаг</Title>
        <Subtitle>Проверь свои знания флагов стран мира!</Subtitle>
        <Description>
          В этой игре тебе предстоит угадывать страны по их флагам. 
          Будут показаны флаги различных стран, а твоя задача - правильно 
          назвать страну. Готов к вызову?
        </Description>

        {error && <ErrorMessage>{error}</ErrorMessage>}

        <StartButton onClick={handleStartGame} disabled={loading}>
          {loading ? 'Создание игры...' : 'Начать игру'}
        </StartButton>
      </GameCard>
    </Container>
  );
};

export default GameStart; 