-- Create queue table
CREATE TABLE IF NOT EXISTS queue (
    id BIGSERIAL PRIMARY KEY,
    author VARCHAR(255) NOT NULL,
    album_id VARCHAR(255),
    chat_id BIGINT NOT NULL,
    msg_id VARCHAR(255) NOT NULL
);

-- Create history table
CREATE TABLE IF NOT EXISTS history (
    id BIGSERIAL PRIMARY KEY,
    album_id VARCHAR(255),
    chat_id BIGINT NOT NULL,
    msg_id VARCHAR(255) NOT NULL,
    posted_at TIMESTAMP NOT NULL
);

-- Insert mock post to history
INSERT INTO history VALUES (0, '', 0, '', '2025-01-01');
