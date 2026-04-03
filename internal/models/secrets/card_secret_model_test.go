package secrets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardUpdateValidate(t *testing.T) {
	tests := []struct {
		name      string
		number    string
		expiry    string
		holder    string
		cvv       string
		shouldErr bool
	}{
		{
			name:      "Update only number",
			number:    "4532015112830366",
			expiry:    "",
			holder:    "",
			cvv:       "",
			shouldErr: false,
		},
		{
			name:      "Update only expiry",
			number:    "",
			expiry:    "12/26",
			holder:    "",
			cvv:       "",
			shouldErr: false,
		},
		{
			name:      "Update all fields",
			number:    "5425233010103441",
			expiry:    "11/27",
			holder:    "Jane Doe",
			cvv:       "456",
			shouldErr: false,
		},
		{
			name:      "No update fields",
			number:    "",
			expiry:    "",
			holder:    "",
			cvv:       "",
			shouldErr: true,
		},
		{
			name:      "Update with invalid number",
			number:    "invalid",
			expiry:    "",
			holder:    "",
			cvv:       "",
			shouldErr: true,
		},
		{
			name:      "Update with invalid expiry",
			number:    "",
			expiry:    "13/26",
			holder:    "",
			cvv:       "",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			card := &Card{
				Number: tt.number,
				Expiry: tt.expiry,
				Holder: tt.holder,
				CVV:    tt.cvv,
			}

			err := card.UpdateValidate()
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCardNumberLuhnAlgorithm(t *testing.T) {
	tests := []struct {
		name     string
		number   string
		expected bool
	}{
		{
			name:     "Valid Visa 4532",
			number:   "4532015112830366",
			expected: true,
		},
		{
			name:     "Valid MasterCard 5425",
			number:   "5425233010103441",
			expected: true,
		},
		{
			name:     "Invalid number fails Luhn",
			number:   "4532015112830367",
			expected: false,
		},
		{
			name:     "Too short",
			number:   "123",
			expected: false,
		},
		{
			name:     "Too long",
			number:   "12345678901234567890",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateNumber(tt.number)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCardFormattingEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Already formatted with spaces",
			input:    "4532 0151 1283 0366",
			expected: "4532 0151 1283 0366",
		},
		{
			name:     "Mixed formatting",
			input:    "4532-0151 1283 0366",
			expected: "4532 0151 1283 0366",
		},
		{
			name:     "No formatting",
			input:    "4532015112830366",
			expected: "4532 0151 1283 0366",
		},
		{
			name:     "Extra spaces",
			input:    "4532  0151  1283  0366",
			expected: "4532 0151 1283 0366",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatCardNumber(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
