# AI Gateway

Production-grade AI Gateway in Go with streaming, fallback routing, rate limiting, and observability.

## Features
- **Streaming**: Support for Server-Sent Events (SSE).
- **Fallback Routing**: Automatic fallback to secondary models if primary fails.
- **Retries**: Configurable retries for primary and fallback targets.
- **Semantic Caching**: Redis-based exact-match caching to reduce latency and costs.
- **Governance (PII Masking)**: Automatic detection and masking of sensitive data (Emails, SSNs, CCs, IPs) before they leave the network.
- **Dynamic Cost Management**: Database-backed model pricing with in-memory caching for accurate cost estimation.
- **Rate Limiting**: Redis-based token-per-minute limiting per tenant.
- **Observability**: OpenTelemetry tracing and metrics.

## Implementation Status
- **OpenAI**: ✅ Fully implemented (including streaming)
- **Anthropic**: ✅ Fully implemented (including streaming)

## Pending Features
- **Extensible**: Plugin system for custom providers and middleware.
- **Multi-tenant Production-ready**: Enhanced isolation, billing integration, and high-availability deployment patterns.
- **Google Gemini**: Support for Google's Gemini models.

## Setup & Running

### 1. Configure Environment Variables

**Option A: Interactive Setup (Recommended)**
```bash
./setup-env.sh
```
This will:
- Create a `.env` file from `.env.example`
- Prompt you for your API keys
- Set up default configurations

**Option B: Manual Setup**
```bash
# Copy the example file
cp .env.example .env

# Edit .env and add your API keys
nano .env
```

**Option C: Export directly (for testing)**
```bash
export OPENAI_API_KEY=your-key-here
export ANTHROPIC_API_KEY=your-key-here
```

### 2. Run Infrastructure
```bash
make docker-up
```
Starts:
- AI Gateway (Go)
- PostgreSQL (Usage & attempts)
- Redis (Rate limiting)

### 3. Rate Limiting Configuration
Set `TOKENS_PER_MINUTE` (default 50,000) in `docker-compose.yml` or via env.

## Usage Examples

### Non-Streaming Request
```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [{"role": "user", "content": "Hello!"}],
    "metadata": {"tenant": "test", "use_case": "support_summary"}
  }'
```

### Streaming Request
```bash
curl -N -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "stream": true,
    "messages": [{"role": "user", "content": "Tell me a long story."}],
    "metadata": {"tenant": "test", "use_case": "support_summary"}
  }'
```

## Observability
By default, traces are exported to stdout. To use an OTLP collector:
```bash
export OTEL_EXPORTER_OTLP_ENDPOINT=http://otel-collector:4318
```

## Database Schema
- `requests`: Final status of each request.
- `provider_attempts`: Detailed log of every primary, retry, and fallback attempt.
- `model_pricing`: Dynamic pricing data for cost estimation.
