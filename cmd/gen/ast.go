package main

type AST interface {
}

type Expr interface {
}

type Program struct {
	Stmts []DefineStmt
}

type DefineStmt struct {
	Ident Ident
	Expr  Expr
}

type SequenceExpr struct {
	Exprs []Expr
}

type ChoiceExpr struct {
	Exprs []Expr
}

type ZeroOrMoreExpr struct {
	Expr Expr
}

type OneOrMoreExpr struct {
	Expr Expr
}

type RepeatExpr struct {
	Expr  Expr
	Limit Limit
}

type OptionalExpr struct {
	Expr Expr
}

type AndExpr struct {
	Expr Expr
}

type NotExpr struct {
	Expr Expr
}

type Charclass struct {
	Invert bool
	Set    []CharRange
}

type CharRange struct {
	Lower rune
	Upper rune
}

type Literal struct {
	Text string
}

type Ident struct {
	Name string
}

type Limit struct {
	Lower      int
	Upper      int
	LowerValid bool
	UpperValid bool
}

type Any struct{}

type EOT struct{}
