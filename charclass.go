package peg

import (
	"unicode/utf8"
)

var _ Expr = &Charclass{}

type Charclass struct {
	set RuneSubset
}

func NewCharclass(set RuneSubset) *Charclass {
	c := new(Charclass)
	c.set = set
	return c
}

func (c *Charclass) Parse(scan *Scanner) (*Tree, bool) {
	t := NewTree(scan.Pos)
	ch, size := utf8.DecodeRuneInString(scan.Text[scan.Pos:])
	if size == 0 {
		return t, false
	}
	if !c.set.Within(ch) {
		return t, false
	}
	scan.Pos += size
	if scan.Pos > scan.LPos {
		scan.LPos = scan.Pos
	}
	t.End = scan.Pos
	return t, true
}
