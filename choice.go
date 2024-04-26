package peg

var _ Expr = &Choice{}

type Choice struct {
	exprs []Expr
}

func NewChoice(exprs ...Expr) *Choice {
	c := new(Choice)
	c.exprs = exprs
	return c
}

func (c *Choice) Parse(scan *Scanner) (*Tree, bool) {
	pos := scan.Pos
	for i, expr := range c.exprs {
		t, ok := expr.Parse(scan)
		if ok {
			t.Index = i
			return t, true
		}
		scan.Pos = pos
	}
	return nil, false
}
