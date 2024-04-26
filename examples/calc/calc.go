package main

import (
	"fmt"
	"strconv"

	"github.com/khirono/go-peg"
)

type Calc struct {
	text string
}

func NewCalc(text string) *Calc {
	c := new(Calc)
	c.text = text
	return c
}

func (c *Calc) Text(t *peg.Tree) string {
	return c.text[t.Start:t.End]
}

func (c *Calc) Eval(t *peg.Tree) (int, error) {
	return c.Program(t)
}

func (c *Calc) Program(t *peg.Tree) (int, error) {
	// program <- expr EOT
	return c.Expr(t.Child[0])
}

func (c *Calc) Expr(t *peg.Tree) (int, error) {
	// expr <- term (("+" / "-") S0 term)*
	v, err := c.Term(t.Child[0])
	if err != nil {
		return 0, err
	}
	for _, child := range t.Child[1].Child {
		op := child.Child[0]
		x, err := c.Term(child.Child[2])
		if err != nil {
			return v, err
		}
		switch op.Index {
		case 0:
			v += x
		case 1:
			v -= x
		default:
			return v, fmt.Errorf("invalid index %v", op.Index)
		}
	}
	return v, nil
}

func (c *Calc) Term(t *peg.Tree) (int, error) {
	// term <- factor (("*" / "/") S0 factor)*
	v, err := c.Factor(t.Child[0])
	if err != nil {
		return 0, err
	}
	for _, child := range t.Child[1].Child {
		op := child.Child[0]
		x, err := c.Factor(child.Child[2])
		if err != nil {
			return v, err
		}
		switch op.Index {
		case 0:
			v *= x
		case 1:
			v /= x
		default:
			return v, fmt.Errorf("invalid index %v", op.Index)
		}
	}
	return v, nil
}

func (c *Calc) Factor(t *peg.Tree) (int, error) {
	// factor <-
	//   "(" S0 expr ")" S0 /
	//   number S0
	switch t.Index {
	case 0:
		return c.Expr(t.Child[2])
	case 1:
		return c.Number(t.Child[0])
	default:
		return 0, fmt.Errorf("invalid index %v", t.Index)
	}
}

func (c *Calc) Number(t *peg.Tree) (int, error) {
	// number <- [1-9] [0-9]* / "0"
	v, err := strconv.Atoi(c.Text(t))
	if err != nil {
		return 0, err
	}
	return int(v), nil
}
