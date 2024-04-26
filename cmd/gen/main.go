package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/khirono/go-peg"
)

func main() {
	var outfile string
	var pkgname string
	var funcname string
	flag.StringVar(&outfile, "outfile", "grammar.go", "output filename")
	flag.StringVar(&pkgname, "pkgname", "main", "package name")
	flag.StringVar(&funcname, "funcname", "NewGrammar", "function name")
	flag.Parse()
	if flag.NArg() < 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	infile := flag.Arg(0)
	prog, err := LoadFile(infile)
	if err != nil {
		fmt.Printf("Load Error: %v\n", err)
		os.Exit(1)
	}
	code, err := GenerateCode(pkgname, funcname, prog)
	if err != nil {
		fmt.Printf("Generate Error: %v\n", err)
		os.Exit(1)
	}
	err = os.WriteFile(outfile, code, 0644)
	if err != nil {
		fmt.Printf("Write Error: %v\n", err)
		os.Exit(1)
	}
}

func LoadFile(filename string) (*Program, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	scan := peg.NewScanner(string(data))
	g := NewPEGGrammar()
	t, ok := g.Parse(scan)
	if !ok {
		fmt.Printf("longest: %q\n", scan.Longest())
		return nil, fmt.Errorf("not accepted")
	}
	b := NewASTBuilder(scan.Text)
	return b.Build(t)
}
