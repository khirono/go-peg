package peg

var _ Expr = &Literal{}

type Literal struct {
	text string
}

func NewLiteral(text string) *Literal {
	l := new(Literal)
	l.text = text
	return l
}

func (l *Literal) Parse(scan *Scanner) (*Tree, bool) {
	t := NewTree(scan.Pos)
	m := len(l.text)
	n := len(scan.Text)
	for i := 0; i < m; i++ {
		if scan.Pos >= n {
			return t, false
		}
		if scan.Text[scan.Pos] != l.text[i] {
			return t, false
		}
		scan.Pos++
		if scan.Pos > scan.LPos {
			scan.LPos = scan.Pos
		}
		t.End = scan.Pos
	}
	return t, true
}
