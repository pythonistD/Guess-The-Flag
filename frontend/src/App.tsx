import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { createGlobalStyle } from 'styled-components';
import AuthForm from './components/Auth/AuthForm';
import GameStart from './components/Game/GameStart';
import GamePlay from './components/Game/GamePlay';
import GameResults from './components/Game/GameResults';
import FlagsDebug from './components/Debug/FlagsDebug';
import ProtectedRoute from './components/ProtectedRoute';
import { TokenUtils } from './services/api';

const GlobalStyle = createGlobalStyle`
  * {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }

  body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
      'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
      sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    background: #f5f5f5;
  }

  html, body, #root {
    height: 100%;
  }

  input, button {
    font-family: inherit;
  }

  button {
    cursor: pointer;
  }

  img {
    max-width: 100%;
    height: auto;
  }
`;

function App() {
  const isAuthenticated = TokenUtils.isAuthenticated();

  return (
    <>
      <GlobalStyle />
      <Router>
        <Routes>
          {/* Публичные маршруты */}
          <Route 
            path="/login" 
            element={
              isAuthenticated ? (
                <Navigate to="/game" replace />
              ) : (
                <AuthForm />
              )
            } 
          />

          {/* Открытый отладочный маршрут: рендерит все флаги в FlagContainer */}
          <Route path="/debug/flags" element={<FlagsDebug />} />
          
          {/* Защищенные маршруты */}
          <Route
            path="/game"
            element={
              <ProtectedRoute>
                <GameStart />
              </ProtectedRoute>
            }
          />
          
          <Route
            path="/game/play"
            element={
              <ProtectedRoute>
                <GamePlay />
              </ProtectedRoute>
            }
          />
          
          <Route
            path="/game/results"
            element={
              <ProtectedRoute>
                <GameResults />
              </ProtectedRoute>
            }
          />
          
          {/* Перенаправления */}
          <Route
            path="/"
            element={
              <Navigate 
                to={isAuthenticated ? "/game" : "/login"} 
                replace 
              />
            }
          />
          
          {/* 404 - страница не найдена */}
          <Route
            path="*"
            element={
              <Navigate 
                to={isAuthenticated ? "/game" : "/login"} 
                replace 
              />
            }
          />
        </Routes>
      </Router>
    </>
  );
}

export default App;
