package utils

import (
	"testing"
)

func TestGenerateID(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"length_5", 5},
		{"length_10", 10},
		{"length_15", 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := GenerateID(tt.length)
			if len(id) != tt.length {
				t.Errorf("GenerateID() generated an ID of length %d, expected %d", len(id), tt.length)
			}
		})
	}
}
