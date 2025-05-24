-- CREATE TABLE IF NOT EXISTS games (
--     id UUID PRIMARY KEY,
--     user_id UUID NOT NULL REFERENCES users(id),
--     opponent_id UUID REFERENCES users(id),
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );

CREATE TABLE IF NOT EXISTS games (
    id UUID PRIMARY KEY,
    mode VARCHAR NOT NULL CHECK (mode IN ('pvp', 'bot')), --pvp/bot
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);