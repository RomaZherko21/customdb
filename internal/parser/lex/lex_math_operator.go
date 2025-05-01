package lex

func lexMathOperator(source string, ic Cursor) (*Token, Cursor, bool) {
	cur := ic

	var options []string
	for _, s := range mathOperators {
		options = append(options, string(s))
	}

	match := longestMatch(source, ic, options)
	if match == "" {
		return nil, ic, false
	}

	cur.Pointer = ic.Pointer + uint(len(match))
	cur.Loc.Col = ic.Loc.Col + uint(len(match))

	return &Token{
		Value: match,
		Loc:   ic.Loc,
		Kind:  MathOperatorToken,
	}, cur, true
}
