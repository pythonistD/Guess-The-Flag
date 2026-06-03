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
  background: #e9edf2;
  padding: 20px;
`;

const AuthCard = styled.div`
  background: #ffffff;
  padding: 3rem;
  border-radius: 15px;
  border: 1px solid #d5dce6;
  box-shadow: 0 8px 24px rgba(31, 41, 55, 0.12);
  width: 100%;
  max-width: 440px;
`;

const FlagIcon = styled.div`
  color: #4b5563;
  font-size: 1rem;
  font-weight: 600;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  margin-bottom: 1rem;
  text-align: center;
`;

const Title = styled.h1`
  color: #1f2937;
  margin-bottom: 0.5rem;
  font-size: 2.5rem;
  text-align: center;
`;

const Subtitle = styled.h2`
  color: #666;
  margin-bottom: 2rem;
  font-weight: 300;
  font-size: 1.1rem;
  text-align: center;
`;

const Form = styled.form`
  display: flex;
  flex-direction: column;
  gap: 1rem;
`;

const Input = styled.input`
  width: 100%;
  padding: 1rem;
  border: 2px solid #e1e5e9;
  border-radius: 8px;
  font-size: 1rem;
  transition: border-color 0.3s;
  background: #ffffff;
  color: #1f2937;

  &::placeholder {
    color: #9ca3af;
  }

  &:focus {
    outline: none;
    border-color: #1f4b99;
  }
`;

const SubmitButton = styled.button`
  padding: 1rem;
  background: #1f4b99;
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 1.1rem;
  cursor: pointer;
  transition: all 0.2s;
  margin-top: 0.5rem;

  &:hover:not(:disabled) {
    background: #183d80;
    transform: translateY(-2px);
  }

  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
    transform: none;
  }
`;

const ToggleText = styled.p`
  text-align: center;
  color: #888;
  margin-top: 1.5rem;
  font-size: 0.95rem;
`;

const ToggleLink = styled.button`
  background: none;
  border: none;
  padding: 0;
  color: #1f4b99;
  cursor: pointer;
  font-size: inherit;
  font-family: inherit;
  text-decoration: underline;
  text-underline-offset: 2px;

  &:hover {
    color: #183d80;
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

interface AuthFormProps {}

const AuthForm: React.FC<AuthFormProps> = () => {
  const [isLogin, setIsLogin] = useState(true);
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      let response;
      if (isLogin) {
        response = await ApiService.login({
          username: formData.username,
          password: formData.password,
        });
      } else {
        response = await ApiService.register({
          username: formData.username,
          email: formData.email,
          password: formData.password,
        });
      }

      TokenUtils.setToken(response.token);
      navigate('/');
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || err.message || 'Произошла ошибка при аутентификации';
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const toggleMode = () => {
    setIsLogin(!isLogin);
    setError(null);
    setFormData({
      username: '',
      email: '',
      password: '',
    });
  };

  return (
    <Container>
      <AuthCard>
        <FlagIcon>Угадай Флаг</FlagIcon>
        <Title>{isLogin ? 'Вход' : 'Регистрация'}</Title>
        <Subtitle>
          {isLogin
            ? 'Войди, чтобы начать угадывать флаги'
            : 'Создай аккаунт и проверь свои знания'}
        </Subtitle>

        {error && <ErrorMessage>{error}</ErrorMessage>}

        <Form onSubmit={handleSubmit}>
          <Input
            type="text"
            name="username"
            placeholder="Имя пользователя"
            value={formData.username}
            onChange={handleInputChange}
            required
          />

          {!isLogin && (
            <Input
              type="email"
              name="email"
              placeholder="Email"
              value={formData.email}
              onChange={handleInputChange}
              required
            />
          )}

          <Input
            type="password"
            name="password"
            placeholder="Пароль"
            value={formData.password}
            onChange={handleInputChange}
            required
          />

          <SubmitButton type="submit" disabled={loading}>
            {loading ? 'Загрузка...' : isLogin ? 'Войти' : 'Зарегистрироваться'}
          </SubmitButton>
        </Form>

        <ToggleText>
          {isLogin ? 'Нет аккаунта? ' : 'Уже есть аккаунт? '}
          <ToggleLink type="button" onClick={toggleMode}>
            {isLogin ? 'Зарегистрироваться' : 'Войти'}
          </ToggleLink>
        </ToggleText>
      </AuthCard>
    </Container>
  );
};

export default AuthForm;
