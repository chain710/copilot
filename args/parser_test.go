package args

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseAssignStatementList(t *testing.T) {
	type expectValue struct {
		left  string
		right string
	}
	tests := []struct {
		name        string
		input       string
		expect      []expectValue
		expectError bool
	}{
		{
			name:  "normal",
			input: "CodeContent=file() ReviewFunction=func1",
			expect: []expectValue{
				{left: "CodeContent", right: "file"},
				{left: "ReviewFunction", right: ".func1"},
			},
		},
		{
			name:  "with arguments",
			input: "ReviewFunction=func1(arg1, `arg2`)",
			expect: []expectValue{
				{left: "ReviewFunction", right: "func1 .arg1 \"arg2\""},
			},
		},
		{
			name:        "begin space",
			input:       " ReviewFunction=func1(arg1, `arg2`)",
			expectError: true,
		},
		{
			name:        "invalid function",
			input:       "ReviewFunction=func1(arg1, `arg2`)()",
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ParseAssignStatementList(tt.input)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, len(tt.expect), len(actual))
				for i, actualStmt := range actual {
					require.Equal(t, tt.expect[i].left, string(actualStmt.Left))
					require.Equal(t, tt.expect[i].right, actualStmt.Right.TemplateText())
				}
			}
		})
	}
}
