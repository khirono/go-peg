package peg

var _ Expr = &Rule{}

type Rule struct {
	name string
	expr Expr
}

func NewRule(name string) *Rule {
	r := new(Rule)
	r.name = name
	return r
}

func (r *Rule) Define(expr Expr) {
	r.expr = expr
}

func (r *Rule) Parse(scan *Scanner) (*Tree, bool) {
	pos := scan.Pos
	memo, ok := scan.Memo(pos, r.name)
	if ok {
		scan.Pos = memo.Pos
		return memo.Tree, true
	}
	t, ok := r.expr.Parse(scan)
	if !ok {
		return t, false
	}
	t.SetTag("rule:" + r.name)
	scan.SetMemo(pos, r.name, Memo{scan.Pos, t})
	return t, true
}
