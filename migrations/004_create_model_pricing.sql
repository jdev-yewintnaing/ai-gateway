CREATE TABLE IF NOT EXISTS model_pricing (
    model VARCHAR(255) PRIMARY KEY,
    provider VARCHAR(255) NOT NULL,
    input_rate_1m DECIMAL(10, 4) NOT NULL,
    output_rate_1m DECIMAL(10, 4) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Seed with initial pricing data
INSERT INTO model_pricing (model, provider, input_rate_1m, output_rate_1m) VALUES
('gpt-4o-mini', 'openai', 0.15, 0.60),
('gpt-4o', 'openai', 5.00, 15.00),
('claude-3-5-sonnet', 'anthropic', 3.00, 15.00),
('claude-3-opus-20240229', 'anthropic', 15.00, 75.00)
ON CONFLICT (model) DO UPDATE SET
    input_rate_1m = EXCLUDED.input_rate_1m,
    output_rate_1m = EXCLUDED.output_rate_1m,
    updated_at = CURRENT_TIMESTAMP;
