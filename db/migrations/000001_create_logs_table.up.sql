CREATE TABLE logs (
    id SERIAL PRIMARY KEY,
    service TEXT NOT NULL,
    level TEXT NOT NULL,
    message TEXT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT now(),
    meta JSONB
);
