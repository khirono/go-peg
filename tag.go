package peg

var _ Expr = &Tag{}

type Tag struct {
	name string
	expr Expr
}

func NewTag(name string, expr Expr) *Tag {
	t := new(Tag)
	t.name = name
	t.expr = expr
	return t
}

func (t *Tag) Parse(scan *Scanner) (*Tree, bool) {
	child, ok := t.expr.Parse(scan)
	if child == nil {
		child = NewTree(scan.Pos)
	}
	child.SetTag(t.name)
	return child, ok
}
