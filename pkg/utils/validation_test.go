package utils_test

import (
	"testing"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/utils"
)

type stringValidationTestCase struct {
	name     string
	in       string
	expected bool
}

func TestValidation_IsAlphaNumeric(t *testing.T) {
	testcases := []stringValidationTestCase{
		{
			name:     "only characters",
			in:       "AStringWithOnlyAlphabeticalCharacters",
			expected: true,
		},
		{
			name:     "only numbers",
			in:       "12345678910",
			expected: true,
		},
		{
			name:     "both alphabetical and number characters",
			in:       "AMix123OfNumbers456AndAlphabetical789Characters",
			expected: true,
		},
		{
			name:     "contains spaces",
			in:       "a string with spaces",
			expected: false,
		},
		{
			name:     "contains symbols",
			in:       "a-string_with/symbols!",
			expected: false,
		},
		{
			name:     "other whitespace",
			in:       "a\ncontains\r\nother\nwhitespace",
			expected: false,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			actual := utils.IsAlphaNumeric(testcase.in)
			if actual != testcase.expected {
				t.Errorf("expected '%v' but was '%v'", testcase.expected, actual)
			}
		})
	}
}
