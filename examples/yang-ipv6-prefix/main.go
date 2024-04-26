package main

import (
	"fmt"

	"github.com/khirono/go-peg"
)

func main() {
	g := NewIPv6PrefixGrammar()
	Do(g, "2001:0db8:85a3:0000:0000:8a2e:0370:7334/64")
	Do(g, "2001:0db8:85a3:0000:0000:8a2e:0370:")
	Do(g, "2001:db8:85a3:0:0:8a2e:370:7334/64")
	Do(g, "::ffff:192.0.2.128")
	Do(g, "::ffff:192.0.2.128ghi")
}

func Do(g peg.Expr, text string) {
	scan := peg.NewScanner(text)
	_, accepted := g.Parse(scan)
	fmt.Printf("accepted: %v\n", accepted)
	fmt.Printf("pos: %v\n", scan.Pos)
	fmt.Printf("lpos: %v\n", scan.LPos)
	matched := scan.Longest()
	fmt.Printf("longest matched: %q\n", matched)
}

func NewIPv6PrefixGrammar() peg.Expr {
	// RFC 6991
	// pattern '((:|[0-9a-fA-F]{0,4}):)([0-9a-fA-F]{0,4}:){0,5}'
	//       + '((([0-9a-fA-F]{0,4}:)?(:|[0-9a-fA-F]{0,4}))|'
	//       + '(((25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])\.){3}'
	//       + '(25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])))'
	//       + '(/(([0-9])|([0-9]{2})|(1[0-1][0-9])|(12[0-8])))';
	// pattern '(([^:]+:){6}(([^:]+:[^:]+)|(.*\..*)))|'
	//       + '((([^:]+:)*[^:]+)?::(([^:]+:)*[^:]+)?)'
	//       + '(/.+)';
	expr := peg.NewRule("expr")
	pattern1 := peg.NewRule("pattern1")
	pattern2 := peg.NewRule("pattern2")
	mask := peg.NewRule("mask")
	hexnum := peg.NewRule("hexnum")
	oct := peg.NewRule("oct")

	// expr -> pattern1 / pattern2
	expr.Define(peg.NewChoice(pattern1, pattern2))

	// pattern1 ->
	//   (
	//     (
	//       ":" /
	//       hexnum
	//     )
	//     ":"
	//   )
	//   (
	//     hexnum
	//     ":"
	//   ){0,5}
	//   (
	//     (
	//       (
	//         oct
	//         "."
	//       ){3}
	//       oct
	//     ) /
	//     (
	//       (
	//         hexnum
	//         ":"
	//       )?
	//       (
	//         ":" /
	//         hexnum
	//       )
	//     )
	//   )
	//   "/"
	//   mask
	pattern1.Define(peg.NewSequence(
		peg.NewChoice(
			peg.NewLiteral(":"),
			hexnum,
		),
		peg.NewLiteral(":"),
		peg.NewRepeat(
			peg.NewSequence(
				hexnum,
				peg.NewLiteral(":"),
			),
			peg.NewLimit(0, 5),
		),
		peg.NewChoice(
			peg.NewSequence(
				peg.NewRepeat(
					peg.NewSequence(
						oct,
						peg.NewLiteral("."),
					),
					peg.NewLimit(3, 3),
				),
				oct,
			),
			peg.NewSequence(
				peg.NewOptional(peg.NewSequence(
					hexnum,
					peg.NewLiteral(":"),
				)),
				peg.NewChoice(
					peg.NewLiteral(":"),
					hexnum,
				),
			),
		),
		peg.NewLiteral("/"),
		mask,
	))

	// mask ->
	//   "12" [0-8] /
	//   "1" [0-1] [0-9] /
	//   [0-9]{2} /
	//   [0-9]
	mask.Define(peg.NewChoice(
		peg.NewSequence(
			peg.NewLiteral("12"),
			peg.NewCharclass(peg.RuneRange{'0', '8'}),
		),
		peg.NewSequence(
			peg.NewLiteral("1"),
			peg.NewCharclass(peg.RuneRange{'0', '1'}),
			peg.NewCharclass(peg.RuneRange{'0', '9'}),
		),
		peg.NewRepeat(
			peg.NewCharclass(peg.RuneRange{'0', '9'}),
			peg.NewLimit(2, 2),
		),
		peg.NewCharclass(peg.RuneRange{'0', '9'}),
	))

	// hexnum -> [0-9a-fA-F]{0,4}
	// XXX: hexnum -> [0-9a-fA-F]{0,4} !"."
	hexnum.Define(peg.NewSequence(
		peg.NewRepeat(
			peg.NewCharclass(peg.RuneUnion{
				peg.RuneRange{'0', '9'},
				peg.RuneRange{'a', 'f'},
				peg.RuneRange{'A', 'F'},
			}),
			peg.NewLimit(0, 4),
		),
		peg.NewNot(peg.NewLiteral(".")),
	))

	// oct ->
	//   "25" [0-5] /
	//   "2" [0-4] [0-9] /
	//   [01]? [0-9]? [0-9]
	oct.Define(peg.NewChoice(
		peg.NewSequence(
			peg.NewLiteral("25"),
			peg.NewCharclass(peg.RuneRange{'0', '5'}),
		),
		peg.NewSequence(
			peg.NewLiteral("2"),
			peg.NewCharclass(peg.RuneRange{'0', '4'}),
			peg.NewCharclass(peg.RuneRange{'0', '9'}),
		),
		peg.NewSequence(
			peg.NewOptional(
				peg.NewCharclass(peg.RuneUnion{
					peg.RuneValue('0'),
					peg.RuneValue('1'),
				}),
			),
			peg.NewOptional(
				peg.NewCharclass(peg.RuneRange{'0', '9'}),
			),
			peg.NewCharclass(peg.RuneRange{'0', '9'}),
		),
	))

	// pattern2 ->
	//   (
	//     [^:]+ ":"
	//   ){6}
	//   (
	//     ([^:]+ ":" [^:]+) /
	//     (.* "." .*)
	//   ) /
	//   (
	//     ([^:]+ ":")*
	//     [^:]+
	//   )?
	//   "::"
	//   (
	//     ([^:]+ ":")*
	//     [^:]+
	//   )?
	//   "/" .+
	pattern2.Define(peg.NewChoice(
		peg.NewSequence(
			peg.NewRepeat(peg.NewSequence(
				peg.NewOneOrMore(peg.NewCharclass(
					peg.RuneInvert{peg.RuneValue(':')},
				)),
				peg.NewLiteral(":"),
			), peg.NewLimit(6, 6)),
			peg.NewChoice(
				peg.NewSequence(
					peg.NewOneOrMore(peg.NewCharclass(
						peg.RuneInvert{peg.RuneValue(':')},
					)),
					peg.NewLiteral(":"),
					peg.NewOneOrMore(peg.NewCharclass(
						peg.RuneInvert{peg.RuneValue(':')},
					)),
				),
				peg.NewSequence(
					peg.NewZeroOrMore(peg.Any),
					peg.NewLiteral("."),
					peg.NewZeroOrMore(peg.Any),
				),
			),
		),
		peg.NewSequence(
			peg.NewOptional(peg.NewSequence(
				peg.NewZeroOrMore(peg.NewSequence(
					peg.NewOneOrMore(peg.NewCharclass(
						peg.RuneInvert{peg.RuneValue(':')},
					)),
					peg.NewLiteral(":"),
				)),
				peg.NewOneOrMore(peg.NewCharclass(
					peg.RuneInvert{peg.RuneValue(':')},
				)),
			)),
			peg.NewLiteral("::"),
			peg.NewOptional(peg.NewSequence(
				peg.NewZeroOrMore(peg.NewSequence(
					peg.NewOneOrMore(peg.NewCharclass(
						peg.RuneInvert{peg.RuneValue(':')},
					)),
					peg.NewLiteral(":"),
				)),
				peg.NewOneOrMore(peg.NewCharclass(
					peg.RuneInvert{peg.RuneValue(':')},
				)),
			)),
			peg.NewLiteral("/"),
			peg.NewOneOrMore(peg.Any),
		),
	))

	return expr
}
