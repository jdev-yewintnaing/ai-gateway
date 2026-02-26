package governance

import (
	"fmt"
	"regexp"
	"strings"
)

type PIIType string

const (
	PIIEmail      PIIType = "EMAIL"
	PIIPhone      PIIType = "PHONE"
	PIICreditCard PIIType = "CREDIT_CARD"
	PIISSN        PIIType = "SSN"
	PIIIPv4       PIIType = "IPV4"
	PIIIPv6       PIIType = "IPV6"
	PIIApiKey     PIIType = "API_KEY"
)

type Detector struct {
	patterns map[PIIType]*regexp.Regexp
}

func NewDetector() *Detector {
	return &Detector{
		patterns: map[PIIType]*regexp.Regexp{
			PIIEmail:      regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`),
			PIIPhone:      regexp.MustCompile(`(\+?\d{1,3}[-.\s]?)?\(?\d{3}\)?[-.\s]?\d{3}[-.\s]?\d{4}`),
			PIICreditCard: regexp.MustCompile(`\b\d{4}[-.\s]?\d{4}[-.\s]?\d{4}[-.\s]?\d{4}\b`),
			PIISSN:        regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`),
			PIIIPv4:       regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`),
			PIIIPv6:       regexp.MustCompile(`\b(?:[a-fA-F0-9]{1,4}:){7}[a-fA-F0-9]{1,4}\b`),
			PIIApiKey:     regexp.MustCompile(`\b(?:sk-[a-zA-Z0-9]{20,}|AIza[0-9A-Za-z-_]{35})\b`),
		},
	}
}

// Mask replaces PII in the text with tokens like [EMAIL_1], [PHONE_1], etc.
// It returns the masked text and a map of tokens to original values.
func (d *Detector) Mask(text string) (string, map[string]string) {
	masked := text
	unmaskMap := make(map[string]string)

	for piiType, re := range d.patterns {
		matches := re.FindAllString(masked, -1)
		for i, match := range matches {
			token := fmt.Sprintf("[%s_%d]", piiType, i+1)
			masked = strings.ReplaceAll(masked, match, token)
			unmaskMap[token] = match
		}
	}

	return masked, unmaskMap
}

// Unmask restores original values from a masked response.
func (d *Detector) Unmask(text string, unmaskMap map[string]string) string {
	unmasked := text
	for token, original := range unmaskMap {
		unmasked = strings.ReplaceAll(unmasked, token, original)
	}
	return unmasked
}
