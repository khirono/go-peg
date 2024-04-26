package peg

var _ Expr = &Not{}

type Not struct {
	expr Expr
}

func NewNot(expr Expr) *Not {
	n := new(Not)
	n.expr = expr
	return n
}

func (n *Not) Parse(scan *Scanner) (*Tree, bool) {
	drop := NewScanner(scan.Text[scan.Pos:])
	_, ok := n.expr.Parse(drop)
	return nil, !ok
}
