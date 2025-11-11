package jet

import (
	"encoding/json"
	"fmt"

	"github.com/tiktoken-go/tokenizer"
)

// TokenComparison holds the comparison results between Jet and JSON
type TokenComparison struct {
	JetBytes     int
	JSONBytes    int
	JetTokens    int
	JSONTokens   int
	ByteSavings  float64 // percentage
	TokenSavings float64 // percentage
}

// CompareTokens compares the token count and byte size between Jet and JSON formats
func CompareTokens(data interface{}) (*TokenComparison, error) {
	// Marshal to Jet format
	jetBytes, err := Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to Jet: %w", err)
	}

	// Marshal to JSON format
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	// Initialize tokenizer (using cl100k_base encoding, commonly used by GPT-4)
	enc, err := tokenizer.Get(tokenizer.Cl100kBase)
	if err != nil {
		return nil, fmt.Errorf("failed to get tokenizer: %w", err)
	}

	// Count tokens
	jetTokens, _, err := enc.Encode(string(jetBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to tokenize Jet output: %w", err)
	}

	jsonTokens, _, err := enc.Encode(string(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to tokenize JSON output: %w", err)
	}

	// Calculate savings
	byteSavings := float64(len(jsonBytes)-len(jetBytes)) / float64(len(jsonBytes)) * 100
	tokenSavings := float64(len(jsonTokens)-len(jetTokens)) / float64(len(jsonTokens)) * 100

	return &TokenComparison{
		JetBytes:     len(jetBytes),
		JSONBytes:    len(jsonBytes),
		JetTokens:    len(jetTokens),
		JSONTokens:   len(jsonTokens),
		ByteSavings:  byteSavings,
		TokenSavings: tokenSavings,
	}, nil
}

// CompareTokensFlattened compares the token count and byte size between Jet (flattened) and JSON formats
func CompareTokensFlattened(data interface{}) (*TokenComparison, error) {
	// Marshal to Jet flattened format
	jetBytes, err := MarshalFlattened(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to Jet (flattened): %w", err)
	}

	// Marshal to JSON format
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	// Initialize tokenizer (using cl100k_base encoding, commonly used by GPT-4)
	enc, err := tokenizer.Get(tokenizer.Cl100kBase)
	if err != nil {
		return nil, fmt.Errorf("failed to get tokenizer: %w", err)
	}

	// Count tokens
	jetTokens, _, err := enc.Encode(string(jetBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to tokenize Jet output: %w", err)
	}

	jsonTokens, _, err := enc.Encode(string(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to tokenize JSON output: %w", err)
	}

	// Calculate savings
	byteSavings := float64(len(jsonBytes)-len(jetBytes)) / float64(len(jsonBytes)) * 100
	tokenSavings := float64(len(jsonTokens)-len(jetTokens)) / float64(len(jsonTokens)) * 100

	return &TokenComparison{
		JetBytes:     len(jetBytes),
		JSONBytes:    len(jsonBytes),
		JetTokens:    len(jetTokens),
		JSONTokens:   len(jsonTokens),
		ByteSavings:  byteSavings,
		TokenSavings: tokenSavings,
	}, nil
}

func CompareTokensNormalized(data interface{}) (*TokenComparison, error) {
	// Marshal to Jet normalized format
	jetBytes, err := MarshalNormalized(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to Jet (normalized): %w", err)
	}

	// Marshal to JSON format
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	// Initialize tokenizer (using cl100k_base encoding, commonly used by GPT-4)
	enc, err := tokenizer.Get(tokenizer.Cl100kBase)
	if err != nil {
		return nil, fmt.Errorf("failed to get tokenizer: %w", err)
	}

	// Count tokens
	jetTokens, _, err := enc.Encode(string(jetBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to tokenize Jet output: %w", err)
	}

	jsonTokens, _, err := enc.Encode(string(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to tokenize JSON output: %w", err)
	}

	// Calculate savings
	byteSavings := float64(len(jsonBytes)-len(jetBytes)) / float64(len(jsonBytes)) * 100
	tokenSavings := float64(len(jsonTokens)-len(jetTokens)) / float64(len(jsonTokens)) * 100

	return &TokenComparison{
		JetBytes:     len(jetBytes),
		JSONBytes:    len(jsonBytes),
		JetTokens:    len(jetTokens),
		JSONTokens:   len(jsonTokens),
		ByteSavings:  byteSavings,
		TokenSavings: tokenSavings,
	}, nil
}
