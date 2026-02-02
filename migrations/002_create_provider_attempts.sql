CREATE TABLE IF NOT EXISTS provider_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id UUID REFERENCES requests(id),
    attempt_no INT NOT NULL,
    provider TEXT NOT NULL,
    model TEXT NOT NULL,
    latency_ms INT,
    status_code INT,
    error_message TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
