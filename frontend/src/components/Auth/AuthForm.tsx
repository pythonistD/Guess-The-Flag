import React, { useState } from 'react';
import styled from 'styled-components';
import { ApiService, TokenUtils } from '../../services/api';
import { useNavigate } from 'react-router-dom';

const Container = styled.div`
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
`;

const FormWrapper = styled.div`
  background: white;
  padding: 2rem;
  border-radius: 10px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
  width: 100%;
  max-width: 400px;
`;

const Title = styled.h1`
  text-align: center;
  color: #333;
  margin-bottom: 2rem;
  font-size: 2rem;
`;

const Form = styled.form`
  display: flex;
  flex-direction: column;
  gap: 1rem;
`;

const Input = styled.input`
  padding: 1rem;
  border: 2px solid #e1e5e9;
  border-radius: 5px;
  font-size: 1rem;
  transition: border-color 0.3s;

  &:focus {
    outline: none;
    border-color: #667eea;
  }
`;

const Button = styled.button`
  padding: 1rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 5px;
  font-size: 1rem;
  cursor: pointer;
  transition: transform 0.2s;

  &:hover {
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
  color: #666;
  margin-top: 1rem;
`;

const ToggleLink = styled.span`
  color: #667eea;
  cursor: pointer;
  text-decoration: underline;

  &:hover {
    color: #764ba2;
  }
`;

const ErrorMessage = styled.div`
  background: #fee;
  color: #c33;
  padding: 0.75rem;
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
    console.log('Form submitted:', { isLogin, formData });
    setLoading(true);
    setError(null);

    try {
      let response;
      if (isLogin) {
        console.log('Attempting login...');
        response = await ApiService.login({
          username: formData.username,
          password: formData.password,
        });
      } else {
        console.log('Attempting registration...');
        response = await ApiService.register({
          username: formData.username,
          email: formData.email,
          password: formData.password,
        });
      }

      console.log('Authentication successful:', response);
      TokenUtils.setToken(response.token);
      navigate('/game');
    } catch (err: any) {
      console.error('Authentication failed:', err);
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
      <FormWrapper>
        <Title>{isLogin ? 'Вход' : 'Регистрация'}</Title>
        
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
          
          <Button type="submit" disabled={loading}>
            {loading ? 'Загрузка...' : isLogin ? 'Войти' : 'Зарегистрироваться'}
          </Button>
        </Form>
        
        <ToggleText>
          {isLogin ? 'Нет аккаунта? ' : 'Уже есть аккаунт? '}
          <ToggleLink onClick={toggleMode}>
            {isLogin ? 'Зарегистрироваться' : 'Войти'}
          </ToggleLink>
        </ToggleText>
      </FormWrapper>
    </Container>
  );
};

export default AuthForm; 