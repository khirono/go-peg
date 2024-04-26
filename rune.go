package peg

type RuneSubset interface {
	Within(rune) bool
}

type RuneAny struct{}

func (RuneAny) Within(x rune) bool {
	return true
}

type RuneValue rune

func (v RuneValue) Within(x rune) bool {
	return x == rune(v)
}

// ^A
type RuneInvert struct {
	S RuneSubset
}

func (i RuneInvert) Within(x rune) bool {
	return !i.S.Within(x)
}

// LowerBound - UpperBound
type RuneRange [2]rune

func (r RuneRange) Within(x rune) bool {
	if x < r[0] {
		return false
	}
	if x > r[1] {
		return false
	}
	return true
}

// A B C ...
type RuneUnion []RuneSubset

func (u RuneUnion) Within(x rune) bool {
	for _, e := range u {
		if e.Within(x) {
			return true
		}
	}
	return false
}
