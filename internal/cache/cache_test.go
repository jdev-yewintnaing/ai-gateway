package cache

import (
	"testing"
)

func TestGenerateKey(t *testing.T) {
	model := "gpt-4"
	messages := []map[string]string{
		{"role": "user", "content": "hello"},
	}

	key1, err := GenerateKey(model, messages)
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}

	key2, err := GenerateKey(model, messages)
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}

	if key1 != key2 {
		t.Errorf("Expected same key for same input, got %s and %s", key1, key2)
	}

	messages2 := []map[string]string{
		{"role": "user", "content": "hello world"},
	}
	key3, err := GenerateKey(model, messages2)
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}

	if key1 == key3 {
		t.Errorf("Expected different keys for different messages, got same key %s", key1)
	}
}
