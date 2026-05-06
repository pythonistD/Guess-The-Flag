-- +goose Up
CREATE TABLE users (
    user_id UUID PRIMARY KEY,
    username VARCHAR NOT NULL UNIQUE,
    email VARCHAR NOT NULL UNIQUE,
    password_hash VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
;
CREATE TABLE images (
    image_id SERIAL PRIMARY KEY,
    svg_data TEXT NOT NULL,
    image_hash VARCHAR(64) UNIQUE NOT NULL, -- SHA-256 хэш для дедупликации
    file_size INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE countries (
    country_id SERIAL PRIMARY KEY,
    code VARCHAR NOT NULL UNIQUE,
    flag_image_id INTEGER REFERENCES images(image_id) NOT NULL
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

CREATE TABLE country_names (
    country_names_id SERIAL PRIMARY KEY,
    language_code VARCHAR(3) NOT NULL,
    country_id INTEGER NOT NULL REFERENCES countries(country_id),
    name TEXT NOT NULL,
    normalized_name TEXT NOT NULL,
    threshold INTEGER NOT NULL,
    is_display_name BOOLEAN NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS unknown_flags;
DROP TABLE IF EXISTS leaderboard;
DROP TABLE IF EXISTS answers;
DROP TABLE IF EXISTS questions;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS countries;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS images;
DROP TABLE IF EXISTS country_names;

