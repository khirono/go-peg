package main

import (
	"fmt"

	"github.com/khirono/go-peg"
)

func main() {
	// g = IPv6address EOT
	g := peg.NewSequence(NewIPv6AddressGrammar(), peg.EOT)
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

func NewIPv6AddressGrammar() peg.Expr {
	// RFC 5954
	// IPv6address =                            6( h16 ":" ) ls32
	//             /                       "::" 5( h16 ":" ) ls32
	//             / [               h16 ] "::" 4( h16 ":" ) ls32
	//             / [ *1( h16 ":" ) h16 ] "::" 3( h16 ":" ) ls32
	//             / [ *2( h16 ":" ) h16 ] "::" 2( h16 ":" ) ls32
	//             / [ *3( h16 ":" ) h16 ] "::"    h16 ":"   ls32
	//             / [ *4( h16 ":" ) h16 ] "::"              ls32
	//             / [ *5( h16 ":" ) h16 ] "::"              h16
	//             / [ *6( h16 ":" ) h16 ] "::"
	//
	// h16         = 1*4HEXDIG
	// ls32        = ( h16 ":" h16 ) / IPv4address
	// IPv4address = dec-octet "." dec-octet "." dec-octet "." dec-octet
	// dec-octet   = DIGIT                 ; 0-9
	//             / %x31-39 DIGIT         ; 10-99
	//             / "1" 2DIGIT            ; 100-199
	//             / "2" %x30-34 DIGIT     ; 200-249
	//             / "25" %x30-35          ; 250-255
	ipv6addr := peg.NewRule("ipv6-addr")
	h16 := peg.NewRule("h16")
	ls32 := peg.NewRule("ls32")
	ipv4addr := peg.NewRule("ipv4-addr")
	decoctet := peg.NewRule("dec-octet")
	digit := peg.NewRule("digit")
	hexdig := peg.NewRule("hexdig")

	// IPv6address =                            6( h16 ":" ) ls32
	//             /                       "::" 5( h16 ":" ) ls32
	//             / [               h16 ] "::" 4( h16 ":" ) ls32
	//             / [ *1( h16 ":" ) h16 ] "::" 3( h16 ":" ) ls32
	//             / [ *2( h16 ":" ) h16 ] "::" 2( h16 ":" ) ls32
	//             / [ *3( h16 ":" ) h16 ] "::"    h16 ":"   ls32
	//             / [ *4( h16 ":" ) h16 ] "::"              ls32
	//             / [ *5( h16 ":" ) h16 ] "::"              h16
	//             / [ *6( h16 ":" ) h16 ] "::"
	//
	ipv6addr.Define(peg.NewChoice(
		peg.NewSequence(
			peg.NewRepeat(peg.NewSequence(
				h16,
				peg.NewLiteral(":"),
			), peg.NewLimit(6, 6)),
			ls32,
		),
		peg.NewSequence(
			peg.NewLiteral("::"),
			peg.NewRepeat(peg.NewSequence(
				h16,
				peg.NewLiteral(":"),
			), peg.NewLimit(5, 5)),
			ls32,
		),
		peg.NewSequence(
			peg.NewOptional(h16),
			peg.NewLiteral("::"),
			peg.NewRepeat(peg.NewSequence(
				h16,
				peg.NewLiteral(":"),
			), peg.NewLimit(4, 4)),
			ls32,
		),
		peg.NewSequence(
			peg.NewOptional(peg.NewSequence(
				peg.NewRepeat(peg.NewSequence(
					h16,
					peg.NewLiteral(":"),
				), peg.NewLimitUpper(1)),
				h16,
			)),
			peg.NewLiteral("::"),
			peg.NewRepeat(peg.NewSequence(
				h16,
				peg.NewLiteral(":"),
			), peg.NewLimit(3, 3)),
			ls32,
		),
		peg.NewSequence(
			peg.NewOptional(peg.NewSequence(
				peg.NewRepeat(peg.NewSequence(
					h16,
					peg.NewLiteral(":"),
				), peg.NewLimitUpper(2)),
				h16,
			)),
			peg.NewLiteral("::"),
			peg.NewRepeat(peg.NewSequence(
				h16,
				peg.NewLiteral(":"),
			), peg.NewLimit(2, 2)),
			ls32,
		),
		peg.NewSequence(
			peg.NewOptional(peg.NewSequence(
				peg.NewRepeat(peg.NewSequence(
					h16,
					peg.NewLiteral(":"),
				), peg.NewLimitUpper(3)),
				h16,
			)),
			peg.NewLiteral("::"),
			h16,
			peg.NewLiteral(":"),
			ls32,
		),
		peg.NewSequence(
			peg.NewOptional(peg.NewSequence(
				peg.NewRepeat(peg.NewSequence(
					h16,
					peg.NewLiteral(":"),
				), peg.NewLimitUpper(4)),
				h16,
			)),
			peg.NewLiteral("::"),
			ls32,
		),
		peg.NewSequence(
			peg.NewOptional(peg.NewSequence(
				peg.NewRepeat(peg.NewSequence(
					h16,
					peg.NewLiteral(":"),
				), peg.NewLimitUpper(5)),
				h16,
			)),
			peg.NewLiteral("::"),
			h16,
		),
		peg.NewSequence(
			peg.NewOptional(peg.NewSequence(
				peg.NewRepeat(peg.NewSequence(
					h16,
					peg.NewLiteral(":"),
				), peg.NewLimitUpper(6)),
				h16,
			)),
			peg.NewLiteral("::"),
		),
	))

	// h16 = 1*4HEXDIG
	// XXX: h16 = 1*4HEXDIG !"."
	h16.Define(peg.NewSequence(
		peg.NewRepeat(hexdig, peg.NewLimit(1, 4)),
		peg.NewNot(peg.NewLiteral(".")),
	))

	// ls32 = ( h16 ":" h16 ) / IPv4address
	ls32.Define(peg.NewChoice(
		peg.NewSequence(
			h16,
			peg.NewLiteral(":"),
			h16,
		),
		ipv4addr,
	))

	// IPv4address = dec-octet "." dec-octet "." dec-octet "." dec-octet
	ipv4addr.Define(peg.NewSequence(
		decoctet,
		peg.NewLiteral("."),
		decoctet,
		peg.NewLiteral("."),
		decoctet,
		peg.NewLiteral("."),
		decoctet,
	))

	// dec-octet   = DIGIT                 ; 0-9
	//             / %x31-39 DIGIT         ; 10-99
	//             / "1" 2DIGIT            ; 100-199
	//             / "2" %x30-34 DIGIT     ; 200-249
	//             / "25" %x30-35          ; 250-255
	decoctet.Define(peg.NewChoice(
		peg.NewSequence(
			peg.NewLiteral("25"),
			peg.NewCharclass(peg.RuneRange{'0', '5'}),
		),
		peg.NewSequence(
			peg.NewLiteral("2"),
			peg.NewCharclass(peg.RuneRange{'0', '4'}),
			digit,
		),
		peg.NewSequence(
			peg.NewLiteral("1"),
			peg.NewRepeat(digit, peg.NewLimit(2, 2)),
		),
		peg.NewSequence(
			peg.NewCharclass(peg.RuneRange{'1', '9'}),
			digit,
		),
		digit,
	))

	digit.Define(peg.NewCharclass(peg.RuneRange{'0', '9'}))

	hexdig.Define(peg.NewCharclass(peg.RuneUnion{
		peg.RuneRange{'0', '9'},
		peg.RuneRange{'a', 'f'},
		peg.RuneRange{'A', 'F'},
	}))

	return ipv6addr
}
