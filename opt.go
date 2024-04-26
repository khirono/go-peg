package peg

var _ Expr = &Optional{}

type Optional struct {
	expr Expr
}

func NewOptional(expr Expr) *Optional {
	o := new(Optional)
	o.expr = expr
	return o
}

func (o *Optional) Parse(scan *Scanner) (*Tree, bool) {
	pos := scan.Pos
	t, ok := o.expr.Parse(scan)
	if !ok {
		scan.Pos = pos
		return nil, true
	}
	return t, true
}
