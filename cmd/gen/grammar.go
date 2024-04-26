package main

import (
	"github.com/khirono/go-peg"
)

func NewPEGGrammar() peg.Expr {
	program := peg.NewRule("program")
	statement := peg.NewRule("statement")
	expression := peg.NewRule("expression")
	sequence := peg.NewRule("sequence")
	term := peg.NewRule("term")
	andpred := peg.NewRule("andpred")
	notpred := peg.NewRule("notpred")
	factor := peg.NewRule("factor")
	repeat := peg.NewRule("repeat")
	primary := peg.NewRule("primary")
	charclass := peg.NewRule("charclass")
	refident := peg.NewRule("refident")
	ident := peg.NewRule("ident")
	literal := peg.NewRule("literal")
	Range := peg.NewRule("Range")
	Char := peg.NewRule("Char")
	digits := peg.NewRule("digits")
	S0 := peg.NewRule("S0")
	space := peg.NewRule("space")

	// program <- S0 (statement S0)* EOT
	program.Define(peg.NewSequence(
		S0,
		peg.NewZeroOrMore(peg.NewSequence(
			statement,
			S0,
		)),
		peg.EOT,
	))

	// statement <- ident S0 "<-" S0 expression
	statement.Define(peg.NewSequence(
		ident,
		S0,
		peg.NewLiteral("<-"),
		S0,
		expression,
	))

	// expression <- sequence ("/" S0 sequence)*
	expression.Define(peg.NewSequence(
		sequence,
		peg.NewZeroOrMore(peg.NewSequence(
			peg.NewLiteral("/"),
			S0,
			sequence,
		)),
	))

	// sequence <- (term S0)+
	sequence.Define(peg.NewOneOrMore(
		peg.NewSequence(
			term,
			S0,
		),
	))

	// term <-
	//   andpred /
	//   notpred /
	//   factor
	term.Define(peg.NewChoice(
		andpred,
		notpred,
		factor,
	))

	// andpred <- "&" factor
	andpred.Define(peg.NewSequence(
		peg.NewLiteral("&"),
		factor,
	))

	// notpred <- "!" factor
	notpred.Define(peg.NewSequence(
		peg.NewLiteral("!"),
		factor,
	))

	// factor <- primary ("?" / "*" / "+" / repeat)?
	factor.Define(peg.NewSequence(
		primary,
		peg.NewOptional(peg.NewChoice(
			peg.NewLiteral("?"),
			peg.NewLiteral("*"),
			peg.NewLiteral("+"),
			repeat,
		)),
	))

	// repeat <- "{" S0 (
	//   digits S0 "," S0 digits /
	//   digits S0 "," /
	//   "," S0 digits /
	//   digits
	//   ) S0 "}"
	repeat.Define(peg.NewSequence(
		peg.NewLiteral("{"),
		S0,
		peg.NewChoice(
			peg.NewSequence(
				digits,
				S0,
				peg.NewLiteral(","),
				S0,
				digits,
			),
			peg.NewSequence(
				digits,
				S0,
				peg.NewLiteral(","),
			),
			peg.NewSequence(
				peg.NewLiteral(","),
				S0,
				digits,
			),
			digits,
		),
		S0,
		peg.NewLiteral("}"),
	))

	// primary <-
	//   "(" S0 expression ")" /
	//   "EOT" /
	//   charclass /
	//   refident /
	//   literal /
	//   "."
	primary.Define(peg.NewChoice(
		peg.NewSequence(
			peg.NewLiteral("("),
			S0,
			expression,
			peg.NewLiteral(")"),
		),
		peg.NewLiteral("EOT"),
		charclass,
		refident,
		literal,
		peg.NewLiteral("."),
	))

	// charclass <- "[" "^"? (!"]" Range)+ "]"
	charclass.Define(peg.NewSequence(
		peg.NewLiteral("["),
		peg.NewOptional(peg.NewLiteral("^")),
		peg.NewOneOrMore(peg.NewSequence(
			peg.NewNot(peg.NewLiteral("]")),
			Range,
		)),
		peg.NewLiteral("]"),
	))

	// refident <- ident !(S0 "<-")
	refident.Define(peg.NewSequence(
		ident,
		peg.NewNot(peg.NewSequence(
			S0,
			peg.NewLiteral("<-"),
		)),
	))

	// ident <- [a-za-Z_] [0-9a-zA-Z_]*
	ident.Define(peg.NewSequence(
		peg.NewCharclass(
			peg.RuneUnion{
				peg.RuneRange{'a', 'z'},
				peg.RuneRange{'A', 'Z'},
				peg.RuneValue('_'),
			},
		),
		peg.NewZeroOrMore(
			peg.NewCharclass(
				peg.RuneUnion{
					peg.RuneRange{'0', '9'},
					peg.RuneRange{'a', 'z'},
					peg.RuneRange{'A', 'Z'},
					peg.RuneValue('_'),
				},
			),
		),
	))

	// literal <-
	//   '"' (!'"' Char)* '"' /
	//   "'' (!"'" Char)* "'"
	literal.Define(peg.NewChoice(
		peg.NewSequence(
			peg.NewLiteral("\""),
			peg.NewZeroOrMore(peg.NewSequence(
				peg.NewNot(peg.NewLiteral("\"")),
				Char,
			)),
			peg.NewLiteral("\""),
		),
		peg.NewSequence(
			peg.NewLiteral("'"),
			peg.NewZeroOrMore(peg.NewSequence(
				peg.NewNot(peg.NewLiteral("'")),
				Char,
			)),
			peg.NewLiteral("'"),
		),
	))

	// Range <- Char "-" Char / Char
	Range.Define(peg.NewChoice(
		peg.NewSequence(
			Char,
			peg.NewLiteral("-"),
			Char,
		),
		Char,
	))

	// Char <-
	//   "\\" [abefnrtv'"\[\]\\] /
	//   "\\" [0-3] [0-7] [0-7] /
	//   "\\" [0-7] [0-7]? /
	//   "\\" "-" /
	//   !"\\" .
	Char.Define(peg.NewChoice(
		peg.NewSequence(
			peg.NewLiteral("\\"),
			peg.NewCharclass(peg.RuneUnion{
				peg.RuneValue('a'),
				peg.RuneValue('b'),
				peg.RuneValue('e'),
				peg.RuneValue('f'),
				peg.RuneValue('n'),
				peg.RuneValue('r'),
				peg.RuneValue('t'),
				peg.RuneValue('v'),
				peg.RuneValue('\''),
				peg.RuneValue('"'),
				peg.RuneValue('['),
				peg.RuneValue(']'),
				peg.RuneValue('\\'),
			}),
		),
		peg.NewSequence(
			peg.NewLiteral("\\"),
			peg.NewCharclass(peg.RuneRange{'0', '3'}),
			peg.NewCharclass(peg.RuneRange{'0', '7'}),
			peg.NewCharclass(peg.RuneRange{'0', '7'}),
		),
		peg.NewSequence(
			peg.NewLiteral("\\"),
			peg.NewCharclass(peg.RuneRange{'0', '7'}),
			peg.NewOptional(
				peg.NewCharclass(peg.RuneRange{'0', '7'}),
			),
		),
		peg.NewSequence(
			peg.NewLiteral("\\"),
			peg.NewLiteral("-"),
		),
		peg.NewSequence(
			peg.NewNot(peg.NewLiteral("\\")),
			peg.Any,
		),
	))

	// digits <- [1-9] [0-9]* / "0"
	digits.Define(peg.NewChoice(
		peg.NewSequence(
			peg.NewCharclass(peg.RuneRange{'1', '9'}),
			peg.NewZeroOrMore(
				peg.NewCharclass(peg.RuneRange{'0', '9'}),
			),
		),
		peg.NewLiteral("0"),
	))

	// S0 <- space*
	S0.Define(peg.NewZeroOrMore(space))

	// space <- [ \t\n\r]
	space.Define(peg.NewCharclass(
		peg.RuneUnion{
			peg.RuneValue(' '),
			peg.RuneValue('\t'),
			peg.RuneValue('\n'),
			peg.RuneValue('\r'),
		},
	))

	return program
}
