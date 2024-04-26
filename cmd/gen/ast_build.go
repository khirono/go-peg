package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/khirono/go-peg"
)

type ASTBuilder struct {
	text string
}

func NewASTBuilder(text string) *ASTBuilder {
	b := new(ASTBuilder)
	b.text = text
	return b
}

func (b *ASTBuilder) Text(t *peg.Tree) string {
	return b.text[t.Start:t.End]
}

func (b *ASTBuilder) Build(t *peg.Tree) (*Program, error) {
	return b.Program(t)
}

func (b *ASTBuilder) Program(t *peg.Tree) (*Program, error) {
	prog := &Program{}
	// program <- S0 (statement S0)* EOT
	for _, child := range t.Child[1].Child {
		stmt, err := b.Statement(child.Child[0])
		if err != nil {
			return prog, err
		}
		prog.Stmts = append(prog.Stmts, *stmt)
	}
	return prog, nil
}

func (b *ASTBuilder) Statement(t *peg.Tree) (*DefineStmt, error) {
	// statement <- ident S0 "<-" S0 expression
	stmt := &DefineStmt{}
	ident, err := b.Ident(t.Child[0])
	if err != nil {
		return stmt, err
	}
	expr, err := b.Expression(t.Child[4])
	if err != nil {
		return stmt, err
	}
	stmt.Ident = *ident
	stmt.Expr = expr
	return stmt, nil
}

func (b *ASTBuilder) Expression(t *peg.Tree) (Expr, error) {
	// expression <- sequence ("/" S0 sequence)*
	expr, err := b.Sequence(t.Child[0])
	if err != nil {
		return nil, err
	}
	if len(t.Child[1].Child) == 0 {
		return expr, nil
	}
	choice := &ChoiceExpr{}
	choice.Exprs = append(choice.Exprs, expr)
	for _, child := range t.Child[1].Child {
		expr, err := b.Sequence(child.Child[2])
		if err != nil {
			return choice, err
		}
		choice.Exprs = append(choice.Exprs, expr)
	}
	return choice, nil
}

func (b *ASTBuilder) Sequence(t *peg.Tree) (Expr, error) {
	// sequence <- (term S0)+
	if len(t.Child) == 1 {
		return b.Term(t.Child[0].Child[0])
	}
	seq := &SequenceExpr{}
	for _, child := range t.Child {
		expr, err := b.Term(child.Child[0])
		if err != nil {
			return seq, err
		}
		seq.Exprs = append(seq.Exprs, expr)
	}
	return seq, nil
}

func (b *ASTBuilder) Term(t *peg.Tree) (Expr, error) {
	// term <-
	//   andpred /
	//   notpred /
	//   factor
	switch t.Index {
	case 0:
		// andpred <- "&" factor
		expr, err := b.Factor(t.Child[1])
		if err != nil {
			return nil, err
		}
		return &AndExpr{expr}, nil
	case 1:
		// notpred <- "!" factor
		expr, err := b.Factor(t.Child[1])
		if err != nil {
			return nil, err
		}
		return &NotExpr{expr}, nil
	case 2:
		return b.Factor(t)
	default:
		return nil, fmt.Errorf("invalid index %v", t.Index)
	}
}

func (b *ASTBuilder) Factor(t *peg.Tree) (Expr, error) {
	// factor <- primary ("?" / "*" / "+" / repeat)?
	expr, err := b.Primary(t.Child[0])
	if err != nil {
		return nil, err
	}
	if t.Child[1] == nil {
		return expr, nil
	}
	switch t.Child[1].Index {
	case 0:
		return &OptionalExpr{expr}, nil
	case 1:
		return &ZeroOrMoreExpr{expr}, nil
	case 2:
		return &OneOrMoreExpr{expr}, nil
	case 3:
		limit, err := b.Repeat(t.Child[1])
		if err != nil {
			return expr, err
		}
		return &RepeatExpr{expr, *limit}, nil
	default:
		return nil, fmt.Errorf("invalid index %v", t.Index)
	}
}

func (b *ASTBuilder) Repeat(t *peg.Tree) (*Limit, error) {
	// repeat <- "{" S0 (
	//   digits S0 "," S0 digits /
	//   digits S0 "," /
	//   "," S0 digits /
	//   digits
	//   ) S0 "}"
	l := &Limit{}
	t = t.Child[2]
	switch t.Index {
	case 0:
		lower, err := b.Digits(t.Child[0])
		if err != nil {
			return nil, err
		}
		upper, err := b.Digits(t.Child[4])
		if err != nil {
			return nil, err
		}
		l.Lower = lower
		l.Upper = upper
		l.LowerValid = true
		l.UpperValid = true
		return l, nil
	case 1:
		lower, err := b.Digits(t.Child[0])
		if err != nil {
			return nil, err
		}
		l.Lower = lower
		l.LowerValid = true
		return l, nil
	case 2:
		upper, err := b.Digits(t.Child[2])
		if err != nil {
			return nil, err
		}
		l.Upper = upper
		l.UpperValid = true
		return l, nil
	case 3:
		v, err := b.Digits(t)
		if err != nil {
			return nil, err
		}
		l.Lower = v
		l.Upper = v
		l.LowerValid = true
		l.UpperValid = true
		return l, nil
	default:
		return nil, fmt.Errorf("invalid index %v", t.Index)
	}
}

func (b *ASTBuilder) Primary(t *peg.Tree) (Expr, error) {
	// primary <-
	//   "(" S0 expression ")" /
	//   "EOT" /
	//   charclass /
	//   refident /
	//   literal /
	//   "."
	switch t.Index {
	case 0:
		return b.Expression(t.Child[2])
	case 1:
		return &EOT{}, nil
	case 2:
		return b.Charclass(t)
	case 3:
		return b.RefIdent(t)
	case 4:
		return b.Literal(t)
	case 5:
		return &Any{}, nil
	default:
		return nil, fmt.Errorf("invalid index %v", t.Index)
	}
}

func (b *ASTBuilder) Charclass(t *peg.Tree) (*Charclass, error) {
	// charclass <- "[" "^"? (!"]" Range)+ "]"
	charclass := &Charclass{}
	if t.Child[1] != nil {
		charclass.Invert = true
	}
	for _, child := range t.Child[2].Child {
		r, err := b.Range(child.Child[1])
		if err != nil {
			return charclass, err
		}
		charclass.Set = append(charclass.Set, *r)
	}
	return charclass, nil
}

func (b *ASTBuilder) RefIdent(t *peg.Tree) (*Ident, error) {
	// refident <- ident !(S0 "<-")
	ident := &Ident{}
	ident.Name = b.Text(t.Child[0])
	return ident, nil
}

func (b *ASTBuilder) Ident(t *peg.Tree) (*Ident, error) {
	// ident <- [a-za-Z_] [0-9a-zA-Z_]*
	ident := &Ident{}
	ident.Name = b.Text(t)
	return ident, nil
}

func (b *ASTBuilder) Literal(t *peg.Tree) (*Literal, error) {
	// literal <-
	//   '"' (!'"' Char)* '"' /
	//   "'' (!"'" Char)* "'"
	l := &Literal{}
	var sb strings.Builder
	for _, child := range t.Child[1].Child {
		ch, err := b.Char(child.Child[1])
		if err != nil {
			return l, err
		}
		sb.WriteRune(ch)
	}
	l.Text = sb.String()
	return l, nil
}

func (b *ASTBuilder) Range(t *peg.Tree) (*CharRange, error) {
	// Range <- Char "-" Char / Char
	r := &CharRange{}
	switch t.Index {
	case 0:
		lower, err := b.Char(t.Child[0])
		if err != nil {
			return nil, err
		}
		r.Lower = lower
		upper, err := b.Char(t.Child[2])
		if err != nil {
			return nil, err
		}
		r.Upper = upper
		return r, nil
	case 1:
		ch, err := b.Char(t)
		if err != nil {
			return nil, err
		}
		r.Lower = ch
		r.Upper = ch
		return r, nil
	default:
		return nil, fmt.Errorf("invalid index %v", t.Index)
	}
}

func (b *ASTBuilder) Char(t *peg.Tree) (rune, error) {
	// Char <-
	//   "\\" [abefnrtv'"\[\]\\] /
	//   "\\" [0-3] [0-7] [0-7] /
	//   "\\" [0-7] [0-7]? /
	//   "\\" "-" /
	//   !"\\" .
	switch t.Index {
	case 0:
		ch, _ := utf8.DecodeRuneInString(b.Text(t.Child[1]))
		switch ch {
		case 'a':
			return '\a', nil
		case 'b':
			return '\b', nil
		case 'e':
			return '\x1b', nil
		case 'f':
			return '\f', nil
		case 'n':
			return '\n', nil
		case 'r':
			return '\r', nil
		case 't':
			return '\t', nil
		case 'v':
			return '\v', nil
		default:
			return ch, nil
		}
	case 1:
		text := b.text[t.Child[1].Start:t.Child[3].End]
		v, err := strconv.ParseUint(text, 8, 32)
		if err != nil {
			return 0, err
		}
		return rune(v), nil
	case 2:
		text := b.Text(t.Child[1])
		if t.Child[2] != nil {
			text += b.Text(t.Child[2])
		}
		v, err := strconv.ParseUint(text, 8, 32)
		if err != nil {
			return 0, err
		}
		return rune(v), nil
	case 3:
		return '-', nil
	case 4:
		ch, _ := utf8.DecodeRuneInString(b.Text(t.Child[1]))
		return ch, nil
	default:
		return 0, fmt.Errorf("invalid index %v", t.Index)
	}
}

func (b *ASTBuilder) Digits(t *peg.Tree) (int, error) {
	// digits <- [1-9] [0-9]* / "0"
	v, err := strconv.ParseInt(b.Text(t), 10, 64)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}
