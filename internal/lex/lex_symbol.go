package lex

func lexSymbol(source string, ic Cursor) (*Token, Cursor, bool) {
	c := source[ic.Pointer]
	cur := ic
	// Will get overwritten later if not an ignored syntax
	cur.Pointer++
	cur.Loc.Col++

	switch c {
	// Syntax that should be thrown away
	case '\n':
		cur.Loc.Line++
		cur.Loc.Col = 0
		fallthrough
	case '\t':
		fallthrough
	case ' ':
		return nil, cur, true
	}

	// Syntax that should be kept
	symbols := []symbol{
		commaSymbol,
		leftparenSymbol,
		rightparenSymbol,
		semicolonSymbol,
		asteriskSymbol,
	}

	var options []string
	for _, s := range symbols {
		options = append(options, string(s))
	}

	// Use `ic`, not `cur`
	match := longestMatch(source, ic, options)
	// Unknown character
	if match == "" {
		return nil, ic, false
	}

	cur.Pointer = ic.Pointer + uint(len(match))
	cur.Loc.Col = ic.Loc.Col + uint(len(match))

	return &Token{
		Value: match,
		Loc:   ic.Loc,
		Kind:  SymbolToken,
	}, cur, true
}
