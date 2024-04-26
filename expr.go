package peg

type Expr interface {
	Parse(s *Scanner) (*Tree, bool)
}
