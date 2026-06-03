-- +goose Up
-- +goose StatementBegin
ALTER TABLE games
    ADD COLUMN game_variant VARCHAR(32) NOT NULL DEFAULT 'text_input';

ALTER TABLE answers
    ADD COLUMN selected_country_id INT REFERENCES countries(country_id);

ALTER TABLE games
    ADD CONSTRAINT games_game_variant_check
        CHECK (game_variant IN ('text_input', 'multiple_choice'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE games
    DROP CONSTRAINT IF EXISTS games_game_variant_check;

ALTER TABLE answers
    DROP COLUMN IF EXISTS selected_country_id;

ALTER TABLE games
    DROP COLUMN IF EXISTS game_variant;
-- +goose StatementEnd
