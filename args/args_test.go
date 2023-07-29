package args

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"text/template"
)

func TestAssignStatementList_Evaluate(t *testing.T) {
	tests := []struct {
		name        string
		exprString  string
		functions   template.FuncMap
		expect      map[string]string
		expectError bool
	}{
		{
			name:       "normal",
			exprString: "CodeContent=ReadFile(`source.go`) ReviewFunction=`func1`",
			functions: template.FuncMap{
				"ReadFile": func(path string) string {
					return fmt.Sprintf("content of %s", path)
				},
			},
			expect: map[string]string{
				"CodeContent":    "content of source.go",
				"ReviewFunction": "func1",
			},
		},
		{
			name:        "ReadFile missing",
			exprString:  "CodeContent=ReadFile(`source.go`) ReviewFunction=`func1`",
			expectError: true,
		},
		{
			name:       "func1 not string literal",
			exprString: "CodeContent=ReadFile(`source.go`) ReviewFunction=func1",
			functions: template.FuncMap{
				"ReadFile": func(path string) string {
					return fmt.Sprintf("content of %s", path)
				},
			},
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse("", []byte(tt.exprString))
			require.NoError(t, err)
			a := result.(AssignStatementList)
			actual, err := a.Evaluate(tt.functions)
			if tt.expectError {
				require.Error(t, err, actual)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expect, actual)
			}
		})
	}
}
