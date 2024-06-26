package parse

import (
	"reflect"
	"testing"
)

func TestSplitSections(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name:     "no sections",
			content:  "This is a test.",
			expected: []string{"This is a test."},
		},
		{
			name:     "one section",
			content:  "This is a test.<!-- stop -->\nAnd this is another test.",
			expected: []string{"This is a test.", "And this is another test."},
		},
		{
			name:     "multiple sections",
			content:  "First.<!-- stop -->\nSecond.<!-- stop -->\nThird.",
			expected: []string{"First.", "Second.", "Third."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitSections(tt.content)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SplitSections() = %v, want %v", result, tt.expected)
			}
		})
	}
}
