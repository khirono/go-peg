package main

import (
	"fmt"
	"os"

	"github.com/khirono/go-peg"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: calc <expr>")
		os.Exit(1)
	}
	scan := peg.NewScanner(os.Args[1])
	g := NewCalcGrammar()
	t, ok := g.Parse(scan)
	if !ok {
		fmt.Println("Error: not accepted.")
		fmt.Printf("longest match: %q\n", scan.Longest())
		os.Exit(1)
	}
	calc := NewCalc(scan.Text)
	val, err := calc.Eval(t)
	if err != nil {
		fmt.Printf("Evaluation Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%v\n", val)
}
