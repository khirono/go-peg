package peg

var _ Expr = &Repeat{}

type Repeat struct {
	expr  Expr
	limit *Limit
}

func NewZeroOrMore(expr Expr) *Repeat {
	return NewRepeat(expr, nil)
}

func NewOneOrMore(expr Expr) *Repeat {
	return NewRepeat(expr, NewLimitLower(1))
}

func NewRepeat(expr Expr, limit *Limit) *Repeat {
	r := new(Repeat)
	r.expr = expr
	if limit == nil {
		limit = &Limit{}
	}
	r.limit = limit
	return r
}

func (r *Repeat) Parse(scan *Scanner) (*Tree, bool) {
	t := NewTree(scan.Pos)
	for !r.limit.Over(len(t.Child)) {
		pos := scan.Pos
		child, ok := r.expr.Parse(scan)
		if !ok {
			scan.Pos = pos
			break
		}
		if scan.Pos == pos {
			break
		}
		t.Append(child)
	}
	if r.limit.Under(len(t.Child)) {
		return t, false
	}
	return t, true
}
