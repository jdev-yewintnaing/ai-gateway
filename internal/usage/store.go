package usage

import (
	"context"
	"math"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Pricing struct {
	Model        string
	InputRate1M  float64
	OutputRate1M float64
}

type Record struct {
	RequestID        string
	Tenant           string
	UseCase          string
	RouteName        string
	Provider         string
	Model            string
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
	CostEstimate     float64
	LatencyMS        int
	StatusCode       int
	ErrorMessage     string
}

type Attempt struct {
	RequestID    string
	AttemptNo    int
	Provider     string
	Model        string
	LatencyMS    int
	StatusCode   int
	ErrorMessage string
}

type Store struct {
	db           *pgxpool.Pool
	pricingCache sync.Map // map[string]Pricing
}

func NewStore(connString string) (*Store, error) {
	db, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) getPricing(ctx context.Context, model string) Pricing {
	if val, ok := s.pricingCache.Load(model); ok {
		return val.(Pricing)
	}

	var p Pricing
	err := s.db.QueryRow(ctx, `
		SELECT model, input_rate_1m, output_rate_1m 
		FROM model_pricing 
		WHERE model = $1
	`, model).Scan(&p.Model, &p.InputRate1M, &p.OutputRate1M)

	if err != nil {
		// Fallback to default if not found in DB
		return Pricing{
			Model:        model,
			InputRate1M:  0.15,
			OutputRate1M: 0.60,
		}
	}

	s.pricingCache.Store(model, p)
	return p
}

func (s *Store) Log(ctx context.Context, r Record) error {
	p := s.getPricing(ctx, r.Model)
	cost := s.EstimateCost(p, r.PromptTokens, r.CompletionTokens)

	_, err := s.db.Exec(ctx, `
		INSERT INTO requests (request_id, tenant, use_case, route_name, provider, model, prompt_tokens, completion_tokens, total_tokens, cost_estimate_usd, latency_ms, status_code, error_message)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (request_id) DO UPDATE SET
			tenant = EXCLUDED.tenant,
			use_case = EXCLUDED.use_case,
			route_name = EXCLUDED.route_name,
			provider = EXCLUDED.provider,
			model = EXCLUDED.model,
			prompt_tokens = EXCLUDED.prompt_tokens,
			completion_tokens = EXCLUDED.completion_tokens,
			total_tokens = EXCLUDED.total_tokens,
			cost_estimate_usd = EXCLUDED.cost_estimate_usd,
			latency_ms = EXCLUDED.latency_ms,
			status_code = EXCLUDED.status_code,
			error_message = EXCLUDED.error_message
	`, r.RequestID, r.Tenant, r.UseCase, r.RouteName, r.Provider, r.Model, r.PromptTokens, r.CompletionTokens, r.TotalTokens, cost, r.LatencyMS, r.StatusCode, r.ErrorMessage)
	return err
}

func (s *Store) LogAttempt(ctx context.Context, reqCorrelationID string, a Attempt) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO provider_attempts (request_id, attempt_no, provider, model, latency_ms, status_code, error_message)
		SELECT id, $2, $3, $4, $5, $6, $7 FROM requests WHERE request_id = $1 LIMIT 1
	`, reqCorrelationID, a.AttemptNo, a.Provider, a.Model, a.LatencyMS, a.StatusCode, a.ErrorMessage)
	return err
}

func (s *Store) Close() {
	s.db.Close()
}

// EstimateCost calculates approximate cost based on provided pricing
func (s *Store) EstimateCost(p Pricing, promptTokens, completionTokens int) float64 {
	cost := (float64(promptTokens) / 1000000.0 * p.InputRate1M) + (float64(completionTokens) / 1000000.0 * p.OutputRate1M)
	return math.Round(cost*1000000) / 1000000 // Round to 6 decimal places
}

// ApproximateTokens is useful for cases where usage isn't returned (e.g. errors before downstream call)
func ApproximateTokens(text string) int {
	return len(text) / 4
}
