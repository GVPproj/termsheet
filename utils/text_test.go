package utils

import "testing"

func TestTruncateText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		maxLen   int
		expected string
	}{
		{
			name:     "text shorter than maxLen",
			text:     "Short",
			maxLen:   10,
			expected: "Short",
		},
		{
			name:     "text equal to maxLen",
			text:     "ExactlyTen",
			maxLen:   10,
			expected: "ExactlyTen",
		},
		{
			name:     "text longer than maxLen",
			text:     "This is a very long item name",
			maxLen:   18,
			expected: "This is a very ...",
		},
		{
			name:     "maxLen is 3",
			text:     "Test",
			maxLen:   3,
			expected: "Tes",
		},
		{
			name:     "maxLen is less than 3",
			text:     "Test",
			maxLen:   2,
			expected: "Te",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TruncateText(tt.text, tt.maxLen)
			if result != tt.expected {
				t.Errorf("TruncateText(%q, %d) = %q, want %q", tt.text, tt.maxLen, result, tt.expected)
			}
		})
	}
}
