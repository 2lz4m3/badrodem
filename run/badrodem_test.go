package run

import (
	"testing"
)

func TestIsProbablyText(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected bool
	}{
		{
			name:     "string a",
			input:    []byte("a"),
			expected: true,
		},
		{
			name:     "binary 0x01",
			input:    []byte{0x01},
			expected: false,
		},
	}

	for _, testCase := range testCases {
		got := isProbablyText(testCase.input)
		if got != testCase.expected {
			t.Errorf("got:%v, %+v", got, testCase)
		}
	}
}
