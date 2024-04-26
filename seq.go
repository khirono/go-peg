package peg

var _ Expr = &Sequence{}

type Sequence struct {
	exprs []Expr
}

func NewSequence(exprs ...Expr) *Sequence {
	s := new(Sequence)
	s.exprs = exprs
	return s
}

func (s *Sequence) Parse(scan *Scanner) (*Tree, bool) {
	t := NewTree(scan.Pos)
	pos := scan.Pos
	for _, expr := range s.exprs {
		child, ok := expr.Parse(scan)
		if !ok {
			scan.Pos = pos
			return t, false
		}
		t.Append(child)
	}
	return t, true
}
