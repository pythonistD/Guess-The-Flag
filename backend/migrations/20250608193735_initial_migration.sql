-- +goose Up
CREATE TABLE users (
    user_id UUID PRIMARY KEY,
    username VARCHAR NOT NULL UNIQUE,
    email VARCHAR NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE countries (
    country_id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE,
    code VARCHAR(2) NOT NULL UNIQUE,
    flag_url VARCHAR NOT NULL
);

CREATE TABLE games (
    game_id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(user_id),
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP
);

CREATE TABLE questions (
    question_id UUID PRIMARY KEY,
    game_id UUID REFERENCES games(game_id),
    country_id INT NOT NULL REFERENCES countries(country_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE answers (
    answer_id UUID PRIMARY KEY,
    question_id UUID REFERENCES questions(question_id),
    answer TEXT,
    answered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_correct BOOLEAN
);

CREATE TABLE leaderboard (
    leaderboard_id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(user_id),
    score INT DEFAULT 0,
    games_played INT DEFAULT 0,
    last_game TIMESTAMP
);

CREATE TABLE unknown_flags (
    unknown_flags_id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(user_id),
    country_id INT REFERENCES countries(country_id)
);

-- +goose Down
DROP TABLE IF EXISTS unknown_flags;
DROP TABLE IF EXISTS leaderboard;
DROP TABLE IF EXISTS answers;
DROP TABLE IF EXISTS questions;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS countries;
DROP TABLE IF EXISTS users;

