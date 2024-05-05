CREATE TABLE IF NOT EXISTS library (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    game_id bigint NOT NULL REFERENCES games ON DELETE CASCADE,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS wallet (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    balance double precision NOT NULL DEFAULT 0
);
