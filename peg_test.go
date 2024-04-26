package peg

import (
	"testing"
)

func newRangeGrammar() Expr {
	// expr -> factor ('|' factor)*
	// factor -> number ('..' number)*
	// number -> [1-9][0-9]* | '0'
	expr := NewRule("expr")
	factor := NewRule("factor")
	number := NewRule("number")

	expr.Define(NewSequence(
		factor,
		NewZeroOrMore(NewSequence(
			NewLiteral("|"),
			factor,
		)),
	))

	factor.Define(NewSequence(
		number,
		NewZeroOrMore(NewSequence(
			NewLiteral(".."),
			number,
		)),
	))

	number.Define(NewChoice(
		NewSequence(
			NewCharclass(RuneRange{'1', '9'}),
			NewZeroOrMore(NewCharclass(RuneRange{'0', '9'})),
		),
		NewLiteral("0"),
	))

	return NewSequence(expr, EOT)
}

func newIPv4PrefixGrammar() Expr {
	// RFC 6991
	// pattern
	//    '(([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])\.){3}'
	//  +  '([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])'
	//  + '/(([0-9])|([1-2][0-9])|(3[0-2]))';
	expr := NewRule("expr")
	addr := NewRule("addr")
	mask := NewRule("mask")
	oct := NewRule("oct")

	// expr -> addr "/" mask
	expr.Define(NewSequence(addr, NewLiteral("/"), mask))

	// addr -> (oct "."){3} oct
	addr.Define(NewSequence(
		NewRepeat(NewSequence(oct, NewLiteral(".")), NewLimit(3, 3)),
		oct,
	))

	// mask -> (
	//   "3" [0-2] /
	//   [1-2] [0-9] /
	//   [0-9]
	mask.Define(NewChoice(
		NewSequence(
			NewLiteral("3"),
			NewCharclass(RuneRange{'0', '2'}),
		),
		NewSequence(
			NewCharclass(RuneRange{'1', '2'}),
			NewCharclass(RuneRange{'0', '9'}),
		),
		NewCharclass(RuneRange{'0', '9'}),
	))

	// oct ->
	//   "25" [0-5] /
	//   "2" [0-4] [0-9] /
	//   "1" [0-9] [0-9] /
	//   [1-9] [0-9] /
	//   [0-9]
	oct.Define(NewChoice(
		NewSequence(
			NewLiteral("25"),
			NewCharclass(RuneRange{'0', '5'}),
		),
		NewSequence(
			NewLiteral("2"),
			NewCharclass(RuneRange{'0', '4'}),
			NewCharclass(RuneRange{'0', '9'}),
		),
		NewSequence(
			NewLiteral("1"),
			NewCharclass(RuneRange{'0', '9'}),
			NewCharclass(RuneRange{'0', '9'}),
		),
		NewSequence(
			NewCharclass(RuneRange{'1', '9'}),
			NewCharclass(RuneRange{'0', '9'}),
		),
		NewCharclass(RuneRange{'0', '9'}),
	))

	return NewSequence(expr, EOT)
}

func TestLongestMatch(t *testing.T) {
	tests := []struct {
		name     string
		g        Expr
		text     string
		matched  string
		accepted bool
	}{
		{
			name:     "exact match",
			g:        newRangeGrammar(),
			text:     "3..15|48..279|4094",
			matched:  "3..15|48..279|4094",
			accepted: true,
		},
		{
			name:     "end of literal",
			g:        newRangeGrammar(),
			text:     "3..15|48..",
			matched:  "3..15|48..",
			accepted: false,
		},
		{
			name:     "middle of literal",
			g:        newRangeGrammar(),
			text:     "3..15|48.",
			matched:  "3..15|48.",
			accepted: false,
		},
		{
			name:     "ipv4-prefix exact match",
			g:        newIPv4PrefixGrammar(),
			text:     "192.168.30.254/24",
			matched:  "192.168.30.254/24",
			accepted: true,
		},
		{
			name:     "ipv4-address",
			g:        newIPv4PrefixGrammar(),
			text:     "192.168.30.254",
			matched:  "192.168.30.254",
			accepted: false,
		},
		{
			name:     "ipv4-address octet 3",
			g:        newIPv4PrefixGrammar(),
			text:     "192.168.30.",
			matched:  "192.168.30.",
			accepted: false,
		},
		{
			name:     "ipv4-address octet 2",
			g:        newIPv4PrefixGrammar(),
			text:     "192.168.",
			matched:  "192.168.",
			accepted: false,
		},
		{
			name:     "ipv4-address octet 1",
			g:        newIPv4PrefixGrammar(),
			text:     "192.",
			matched:  "192.",
			accepted: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			scan := NewScanner(tc.text)
			_, accepted := tc.g.Parse(scan)
			if accepted != tc.accepted {
				t.Errorf("want %v; but got %v", tc.accepted, accepted)
			}
			matched := scan.Longest()
			if matched != tc.matched {
				t.Errorf("want %q; but got %q", tc.matched, matched)
			}
		})
	}
}
