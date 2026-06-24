CREATE SCHEMA IF NOT EXISTS visitor;

CREATE TABLE IF NOT EXISTS visitor.visitors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    data JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
