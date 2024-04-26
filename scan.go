package peg

type Memo struct {
	Pos  int
	Tree *Tree
}

type Scanner struct {
	Text string
	Pos  int
	LPos int
	memo map[int]map[string]Memo
}

func NewScanner(text string) *Scanner {
	s := new(Scanner)
	s.Text = text
	s.memo = make(map[int]map[string]Memo)
	return s
}

func (s *Scanner) Longest() string {
	return s.Text[:s.LPos]
}

func (s *Scanner) Memo(pos int, name string) (Memo, bool) {
	x, ok := s.memo[pos]
	if !ok {
		return Memo{}, false
	}
	memo, ok := x[name]
	return memo, ok
}

func (s *Scanner) SetMemo(pos int, name string, memo Memo) {
	x, ok := s.memo[pos]
	if !ok {
		x = make(map[string]Memo)
	}
	x[name] = memo
	s.memo[pos] = x
}
