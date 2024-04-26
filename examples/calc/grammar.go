package main

import (
	"github.com/khirono/go-peg"
)

func NewCalcGrammar() peg.Expr {
	program := peg.NewRule("program")
	expr := peg.NewRule("expr")
	term := peg.NewRule("term")
	factor := peg.NewRule("factor")
	number := peg.NewRule("number")
	S0 := peg.NewRule("S0")
	space := peg.NewRule("space")

	// program <- expr EOT
	program.Define(peg.NewSequence(expr, peg.EOT))

	// expr <- term (("+" / "-") S0 term)*
	expr.Define(peg.NewSequence(
		term,
		peg.NewZeroOrMore(peg.NewSequence(
			peg.NewChoice(
				peg.NewLiteral("+"),
				peg.NewLiteral("-"),
			),
			S0,
			term,
		)),
	))

	// term <- factor (("*" / "/") S0 factor)*
	term.Define(peg.NewSequence(
		factor,
		peg.NewZeroOrMore(peg.NewSequence(
			peg.NewChoice(
				peg.NewLiteral("*"),
				peg.NewLiteral("/"),
			),
			S0,
			factor,
		)),
	))

	// factor <-
	//   "(" S0 expr ")" S0 /
	//   number S0
	factor.Define(peg.NewChoice(
		peg.NewSequence(
			peg.NewLiteral("("),
			S0,
			expr,
			peg.NewLiteral(")"),
			S0,
		),
		peg.NewSequence(
			number,
			S0,
		),
	))

	// number <- [1-9] [0-9]* / "0"
	number.Define(peg.NewChoice(
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
	space.Define(peg.NewCharclass(peg.RuneUnion{
		peg.RuneValue(' '),
		peg.RuneValue('\t'),
		peg.RuneValue('\n'),
		peg.RuneValue('\r'),
	}))

	return program
}
