package peg

var _ Expr = &And{}

type And struct {
	expr Expr
}

func NewAnd(expr Expr) *And {
	a := new(And)
	a.expr = expr
	return a
}

func (a *And) Parse(scan *Scanner) (*Tree, bool) {
	drop := NewScanner(scan.Text[scan.Pos:])
	_, ok := a.expr.Parse(drop)
	return nil, ok
}
