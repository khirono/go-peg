package main

//go:generate go run ../../cmd/gen ipv6addr.peg

import (
	"fmt"

	"github.com/khirono/go-peg"
)

func main() {
	g := peg.NewSequence(NewGrammar(), peg.EOT)
	Do(g, "2001:0db8:85a3:0000:0000:8a2e:0370:7334")
	Do(g, "2001:0db8:85a3:0000:0000:8a2e:0370:")
	Do(g, "2001:db8:85a3:0:0:8a2e:370:7334")
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
