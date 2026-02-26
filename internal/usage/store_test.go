package usage

import (
	"testing"
)

func TestStore_EstimateCost(t *testing.T) {
	s := &Store{}

	tests := []struct {
		name             string
		pricing          Pricing
		promptTokens     int
		completionTokens int
		want             float64
	}{
		{
			name: "GPT-4o Mini pricing",
			pricing: Pricing{
				InputRate1M:  0.15,
				OutputRate1M: 0.60,
			},
			promptTokens:     1000000,
			completionTokens: 1000000,
			want:             0.75,
		},
		{
			name: "Claude 3.5 Sonnet pricing",
			pricing: Pricing{
				InputRate1M:  3.00,
				OutputRate1M: 15.00,
			},
			promptTokens:     1000,
			completionTokens: 1000,
			want:             0.018, // (1000/1M * 3) + (1000/1M * 15) = 0.003 + 0.015 = 0.018
		},
		{
			name: "Cheap model pricing",
			pricing: Pricing{
				InputRate1M:  0.01,
				OutputRate1M: 0.02,
			},
			promptTokens:     100,
			completionTokens: 100,
			want:             0.000003, // (100/1M * 0.01) + (100/1M * 0.02) = 0.000001 + 0.000002 = 0.000003
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.EstimateCost(tt.pricing, tt.promptTokens, tt.completionTokens); got != tt.want {
				t.Errorf("Store.EstimateCost() = %v, want %v", got, tt.want)
			}
		})
	}
}
