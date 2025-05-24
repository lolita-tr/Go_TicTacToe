--
-- CREATE TABLE IF NOT EXISTS game_results (
--     game_id UUID NOT NULL REFERENCES games(id),
--     is_bot_win BOOLEAN NOT NULL DEFAULT false,
--     winner_id UUID REFERENCES users(id),
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );

CREATE TABLE IF NOT EXISTS game_results (
    user_id UUID NOT NULL REFERENCES users(id),
    game_id UUID NOT NULL REFERENCES games(id),
    result VARCHAR NOT NULL CHECK (result IN ('win', 'lose', 'draw'))--win/lose/draw
);


