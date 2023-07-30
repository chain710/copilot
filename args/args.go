package args

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/chain710/dev_agent/log"
	"strings"
	"text/template"
)

type Expression interface {
	TemplateText() string // to text/template text without {{ }}
	String() string
}

type AssignStatement struct {
	Left  simpleVariable
	Right Expression
}

func (a AssignStatement) String() string {
	return fmt.Sprintf("%s = %s", a.Left.String(), a.Right.String())
}

type AssignStatementList []AssignStatement

func (a AssignStatementList) String() string {
	parts := make([]string, len(a))
	for i, assign := range a {
		parts[i] = assign.String()
	}
	return strings.Join(parts, "\n")
}

func (a AssignStatementList) Evaluate(functions template.FuncMap) (map[string]string, error) {
	tpl := template.New("assign").Funcs(functions).Option("missingkey=error")
	result := make(map[string]string)
	for _, assign := range a {
		var buf bytes.Buffer
		tplText := fmt.Sprintf(`{{ %s }}`, assign.Right.TemplateText())
		tpl, err := tpl.Parse(tplText)
		if err != nil {
			log.Errorf("parse template `%s` error: %s", tplText, err)
			return nil, err
		}
		if err := tpl.Execute(&buf, nil); err != nil {
			log.Errorf("execute assign statement `%s` error: %s", tplText, err)
			return nil, err
		}
		result[string(assign.Left)] = buf.String()
	}

	return result, nil
}

type stringLiteral string

func (s stringLiteral) TemplateText() string {
	return fmt.Sprintf(`"%s"`, string(s))
}

func (s stringLiteral) String() string {
	return fmt.Sprintf("stringLiteral: `%s`", string(s))
}

type simpleVariable string

func (s simpleVariable) TemplateText() string {
	return "." + string(s) // variable in dot
}

func (s simpleVariable) String() string {
	return fmt.Sprintf("simpleVariable: %s", string(s))
}

type functionExpression struct {
	function string
	args     []Expression
}

func (f functionExpression) TemplateText() string {
	parts := make([]string, 0, len(f.args)+1)
	parts = append(parts, f.function)
	for _, arg := range f.args {
		parts = append(parts, arg.TemplateText())
	}

	return strings.Join(parts, " ")
}

func (f functionExpression) String() string {
	parts := make([]string, len(f.args))
	for i, arg := range f.args {
		parts[i] = arg.String()
	}
	return fmt.Sprintf("%s(%s)", f.function, strings.Join(parts, ", "))
}

func newExpression(e any) (Expression, error) {
	expr, ok := e.(Expression)
	if !ok {
		return nil, fmt.Errorf("type not supported: %T", e)
	}

	return expr, nil
}
func newMultiExpression(e1, e2 any) ([]Expression, error) {
	exprList := e2.([]any)
	result := make([]Expression, 0, len(exprList)+1)
	result = append(result, e1.(Expression))
	for _, a := range exprList {
		result = append(result, a.(Expression))
	}
	return result, nil
}

func newExpressionList(e any) []Expression {
	switch val := e.(type) {
	case []any:
		if len(val) != 0 {
			panic(fmt.Errorf("[]any should be empty, actual: %v", val))
		}
		return nil
	case []Expression:
		return val
	default:
		panic(fmt.Errorf("type not supported: %T", e))
	}
}

func newQuotedString(raw []byte) (stringLiteral, error) {
	if len(raw) < 2 {
		return "", errors.New("invalid quoted string: insufficient length")
	}

	if raw[0] != raw[len(raw)-1] || (raw[0] != '`' && raw[0] != '"' && raw[0] != '\'') {
		return "", errors.New("invalid quoted string: begin char should match end char")
	}

	return stringLiteral(raw[1 : len(raw)-1]), nil
}

func newAssignStatement(v, e any) AssignStatement {
	return AssignStatement{Left: v.(simpleVariable), Right: e.(Expression)}
}

func newAssignStatementList(a1, a2 any) AssignStatementList {
	stmtList := a2.([]any)
	result := make([]AssignStatement, 0, len(stmtList)+1)
	result = append(result, a1.(AssignStatement))
	for _, a := range stmtList {
		result = append(result, a.(AssignStatement))
	}
	return result
}
