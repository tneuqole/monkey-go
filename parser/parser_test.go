package parser

import (
	"fmt"
	"testing"

	"github.com/tneuqole/monkey-go/ast"
	"github.com/tneuqole/monkey-go/lexer"
)

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("token literal not 'let', got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not a let statement, got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("value not %s, got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("name not %s, got=%s", name, letStmt.Name)
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, val int64) bool {
	il, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("exp not IntegerLiteral, got=%T", exp)
		return false
	}

	if il.Value != val {
		t.Errorf("Value not %d, got=%d", val, il.Value)
		return false
	}

	if il.TokenLiteral() != fmt.Sprintf("%d", val) {
		t.Errorf("TokenLiteral not %d, got=%s", val, il.TokenLiteral())
		return false
	}

	return true
}

func testBoolean(t *testing.T, exp ast.Expression, val bool) bool {
	b, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not Boolean, got=%T", exp)
		return false
	}

	if b.Value != val {
		t.Errorf("Value not %t, got=%t", val, b.Value)
		return false
	}

	if b.TokenLiteral() != fmt.Sprintf("%t", val) {
		t.Errorf("TokenLiteral not %t, got=%s", val, b.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, val string) bool {
	i, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not Identifier, got=%T", exp)
		return false
	}

	if i.Value != val {
		t.Errorf("Value not %s, got=%s", val, i.Value)
		return false
	}

	if i.TokenLiteral() != val {
		t.Errorf("TokenLiteral not %s, got=%s", val, i.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBoolean(t, exp, v)
	}

	t.Errorf("type of exp not handled, got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	ie, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp not InfixExpression, got=%T(%s)", exp, exp)
	}

	if !testLiteralExpression(t, ie.Left, left) {
		return false
	}

	if ie.Operator != operator {
		t.Errorf("Operator is not %s, got=%q", operator, ie.Operator)
	}

	if !testLiteralExpression(t, ie.Right, right) {
		return false
	}

	return false
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedVal   interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = 10;", "y", 10},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)
		if program == nil {
			t.Fatalf("program is nil")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program does not have 1 statements, got=%d", len(program.Statements))
		}
		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdent) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		testLiteralExpression(t, val, tt.expectedVal)
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input       string
		expectedVal interface{}
	}{
		{"return 5;", 5},
		{"return 10;", 10},
		{"return 993322;", 993322},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)
		if program == nil {
			t.Fatalf("program is nil")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program does not have 1 statements, got=%d", len(program.Statements))
		}

		val := program.Statements[0].(*ast.ReturnStatement).ReturnValue
		testLiteralExpression(t, val, tt.expectedVal)
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program doesn't have enough statements, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("not an expression statement, got=%T", program.Statements[0])
	}

	testIdentifier(t, stmt.Expression, "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program doesn't have enough statements, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("not an expression statement, got=%T", program.Statements[0])
	}

	testIntegerLiteral(t, stmt.Expression, 5)
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		val      interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program doesn't have 1 statement, got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("not an expression statement, got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("not a PrefixExpression, got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %q, got=%s", tt.operator, exp.Operator)
		}

		testLiteralExpression(t, exp.Right, tt.val)
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program doesn't have 1 statement, got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("not an expression statement, got=%T", program.Statements[0])
		}

		testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{" 2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
		{"a * [1, 2, 3, 4][b * c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
		{"add(a * b[2], b[1], 2 * [1, 2][1])", "add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}

	}
}

func TestBooleanExpression(t *testing.T) {
	input := "true"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program doesn't have enough statements, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("not an expression statement, got=%T", program.Statements[0])
	}

	testBoolean(t, stmt.Expression, true)
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program doesn't have 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("not an ExpressionStatement, got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("not an IfExpression, got=%T", exp)
	}

	testInfixExpression(t, exp.Condition, "x", "<", "y")

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement, got=%d", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("consequence stmt not an ExpressionStatement, got=%T", exp.Consequence.Statements[0])
	}

	testIdentifier(t, consequence.Expression, "x")

	if exp.Alternative != nil {
		t.Errorf("Alternative not nil, got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program doesn't have 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("not an ExpressionStatement, got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("not an IfExpression, got=%T", exp)
	}

	testInfixExpression(t, exp.Condition, "x", "<", "y")

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement, got=%d", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("consequence stmt not an ExpressionStatement, got=%T", exp.Consequence.Statements[0])
	}

	testIdentifier(t, consequence.Expression, "x")

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("Alternative is not 1 statement, got=%d", len(exp.Consequence.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Alternative stmt not an ExpressionStatement, got=%T", exp.Alternative.Statements[0])
	}

	testIdentifier(t, alternative.Expression, "y")
}

func TestFunctionLiteral(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program doesn't have 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("not an ExpressionStatement, got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("not an FunctionLiteral, got=%T", exp)
	}

	if len(exp.Parameters) != 2 {
		t.Errorf("want 2 parameters, got=%d", len(exp.Parameters))
	}

	testLiteralExpression(t, exp.Parameters[0], "x")
	testLiteralExpression(t, exp.Parameters[1], "y")

	if len(exp.Body.Statements) != 1 {
		t.Errorf("body should have 1 stmt, got=%d", len(exp.Body.Statements))
	}

	body, ok := exp.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Body stmt not an ExpressionStatement, got=%T", exp.Body.Statements[0])
	}

	testInfixExpression(t, body.Expression, "x", "+", "y")
}

func TestFunctionParameter(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)

		exp := stmt.Expression.(*ast.FunctionLiteral)

		if len(exp.Parameters) != len(tt.expectedParams) {
			t.Errorf("want %d parameters, got=%d", len(tt.expectedParams), len(exp.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, exp.Parameters[i], ident)
		}
	}
}

func TestCallExpression(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5)"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program doesn't have 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("not an ExpressionStatement, got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("not an CallExpression, got=%T", exp)
	}

	testIdentifier(t, exp.Function, "add")
	if len(exp.Arguments) != 3 {
		t.Errorf("want 3 arguments, got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestCallExpressionParameterParsing(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{
			input:         "add();",
			expectedIdent: "add",
			expectedArgs:  []string{},
		},
		{
			input:         "add(1);",
			expectedIdent: "add",
			expectedArgs:  []string{"1"},
		},
		{
			input:         "add(1, 2 * 3, 4 + 5);",
			expectedIdent: "add",
			expectedArgs:  []string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		exp, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
				stmt.Expression)
		}

		if !testIdentifier(t, exp.Function, tt.expectedIdent) {
			return
		}

		if len(exp.Arguments) != len(tt.expectedArgs) {
			t.Fatalf("wrong number of arguments. want=%d, got=%d",
				len(tt.expectedArgs), len(exp.Arguments))
		}

		for i, arg := range tt.expectedArgs {
			if exp.Arguments[i].String() != arg {
				t.Errorf("argument %d wrong. want=%q, got=%q", i,
					arg, exp.Arguments[i].String())
			}
		}
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not StringLiteral, got=%T", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q, got=%q", "hello world", literal.Value)
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ArrayLiteral, got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	ie, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not IndexExpression, got=%T", stmt.Expression)
	}

	if !testIdentifier(t, ie.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, ie.Index, 1, "+", 1) {
		return
	}
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp not HashLiteral, got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs wrong length, got=%d", len(hash.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	for k, v := range hash.Pairs {
		l, ok := k.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral, got=%T", k)
		}

		testIntegerLiteral(t, v, expected[l.String()])
	}
}

func TestParsingHashLiteralsExpVals(t *testing.T) {
	input := `{"one": 1 + 0, "two": 10 - 8, "three": 15 / 5}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp not HashLiteral, got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs wrong length, got=%d", len(hash.Pairs))
	}

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 1, "+", 0)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}
	for k, v := range hash.Pairs {
		l, ok := k.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral, got=%T", k)
			continue
		}

		testf, ok := tests[l.String()]
		if !ok {
			t.Errorf("No test func for key %q", l.String())
			continue
		}
		testf(v)
	}
}

func TestParsingHashLiteralsIntKeys(t *testing.T) {
	input := `{1: 1, 2: 2, 3: 3}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp not HashLiteral, got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs wrong length, got=%d", len(hash.Pairs))
	}

	expected := map[int64]int64{
		1: 1,
		2: 2,
		3: 3,
	}
	for k, v := range hash.Pairs {
		l, ok := k.(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("key is not ast.IntegerLiteral, got=%T", k)
		}

		testIntegerLiteral(t, v, expected[l.Value])
	}
}

func TestParsingHashLiteralsBoolKeys(t *testing.T) {
	input := `{true: true, false: false}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp not HashLiteral, got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 2 {
		t.Errorf("hash.Pairs wrong length, got=%d", len(hash.Pairs))
	}

	expected := map[bool]bool{
		true:  true,
		false: false,
	}
	for k, v := range hash.Pairs {
		kl, ok := k.(*ast.Boolean)
		if !ok {
			t.Errorf("key is not ast.Boolean, got=%T", k)
		}
		vl, ok := v.(*ast.Boolean)
		if !ok {
			t.Errorf("value is not ast.Boolean, got=%T", k)
		}

		if !(vl.Value == expected[kl.Value]) {
			t.Errorf("value not correct, expected=%t, got=%t", expected[kl.Value], vl.Value)
		}
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := "{}"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp not HashLiteral, got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 0 {
		t.Errorf("hash.Pairs wrong length, got=%d", len(hash.Pairs))
	}
}

func TestMacroLiteralParsing(t *testing.T) {
	input := "macro(x, y) { x + y; }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not have 1 statement")
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	macro, ok := stmt.Expression.(*ast.MacroLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.MacroLiteral. got=%T", stmt.Expression)
	}

	if len(macro.Parameters) != 2 {
		t.Fatalf("macro literal parameters wrong. want=2, got=%d", len(macro.Parameters))
	}

	testLiteralExpression(t, macro.Parameters[0], "x")
	testLiteralExpression(t, macro.Parameters[1], "y")

	if len(macro.Body.Statements) != 1 {
		t.Fatalf("macro body does not have 1 statement. got=%d", len(macro.Body.Statements))
	}

	bodyStmt, ok := macro.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("macro body stmt is not ast.ExpressionStatement. got=%T", macro.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionLiteralWithName(t *testing.T) {
	input := `let myFunction = fn() { };`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.LetStatement. got=%T",
			program.Statements[0])
	}
	function, ok := stmt.Value.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Value is not ast.FunctionLiteral. got=%T",
			stmt.Value)
	}
	if function.Name != "myFunction" {
		t.Fatalf("function literal name wrong. want 'myFunction', got=%q\n",
			function.Name)
	}
}
