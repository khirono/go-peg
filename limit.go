package peg

type Limit struct {
	lower      int
	upper      int
	lowervalid bool
	uppervalid bool
}

func NewLimitLower(lower int) *Limit {
	l := new(Limit)
	l.lower = lower
	l.lowervalid = true
	return l
}

func NewLimitUpper(upper int) *Limit {
	l := new(Limit)
	l.upper = upper
	l.uppervalid = true
	return l
}

func NewLimit(lower, upper int) *Limit {
	l := new(Limit)
	l.lower = lower
	l.upper = upper
	l.lowervalid = true
	l.uppervalid = true
	return l
}

func (l *Limit) Over(val int) bool {
	return l.uppervalid && val >= l.upper
}

func (l *Limit) Under(val int) bool {
	return l.lowervalid && val < l.lower
}
