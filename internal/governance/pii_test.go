package governance

import (
	"testing"
)

func TestDetector_Mask(t *testing.T) {
	d := NewDetector()

	tests := []struct {
		name     string
		input    string
		wantMask string
		wantPII  []string
	}{
		{
			name:     "Email masking",
			input:    "Contact me at john.doe@example.com for details.",
			wantMask: "Contact me at [EMAIL_1] for details.",
			wantPII:  []string{"john.doe@example.com"},
		},
		{
			name:     "Phone number masking",
			input:    "Call 123-456-7890.",
			wantMask: "Call [PHONE_1].",
			wantPII:  []string{"123-456-7890"},
		},
		{
			name:     "Credit card masking",
			input:    "My card is 1234-5678-1234-5678.",
			wantMask: "My card is [CREDIT_CARD_1].",
			wantPII:  []string{"1234-5678-1234-5678"},
		},
		{
			name:     "SSN masking",
			input:    "My SSN is 123-45-6789.",
			wantMask: "My SSN is [SSN_1].",
			wantPII:  []string{"123-45-6789"},
		},
		{
			name:     "IPv4 masking",
			input:    "Server is at 192.168.1.1.",
			wantMask: "Server is at [IPV4_1].",
			wantPII:  []string{"192.168.1.1"},
		},
		{
			name:     "IPv6 masking",
			input:    "IPv6 is 2001:0db8:85a3:0000:0000:8a2e:0370:7334.",
			wantMask: "IPv6 is [IPV6_1].",
			wantPII:  []string{"2001:0db8:85a3:0000:0000:8a2e:0370:7334"},
		},
		{
			name:     "API Key masking",
			input:    "Key: sk-abcdefghijklmnopqrstuvwxyz123456",
			wantMask: "Key: [API_KEY_1]",
			wantPII:  []string{"sk-abcdefghijklmnopqrstuvwxyz123456"},
		},
		{
			name:     "Multiple PII",
			input:    "Email john@test.com or call 555-555-5555.",
			wantMask: "Email [EMAIL_1] or call [PHONE_1].",
			wantPII:  []string{"john@test.com", "555-555-5555"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMask, gotMap := d.Mask(tt.input)
			if gotMask != tt.wantMask {
				t.Errorf("Detector.Mask() gotMask = %v, want %v", gotMask, tt.wantMask)
			}
			if len(gotMap) != len(tt.wantPII) {
				t.Errorf("Detector.Mask() gotMap size = %v, want %v", len(gotMap), len(tt.wantPII))
			}
		})
	}
}

func TestDetector_Unmask(t *testing.T) {
	d := NewDetector()
	input := "Hello [EMAIL_1], your code is [PHONE_1]."
	unmaskMap := map[string]string{
		"[EMAIL_1]": "test@test.com",
		"[PHONE_1]": "123-456-7890",
	}
	want := "Hello test@test.com, your code is 123-456-7890."

	if got := d.Unmask(input, unmaskMap); got != want {
		t.Errorf("Detector.Unmask() = %v, want %v", got, want)
	}
}
