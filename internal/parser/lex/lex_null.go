package lex

func lexNull(source string, ic Cursor) (*Token, Cursor, bool) {
	cur := ic

	match := longestMatch(source, ic, []string{string(NullValueKeyword)})
	if match == "" {
		return nil, ic, false
	}

	cur.Pointer = ic.Pointer + uint(len(match))
	cur.Loc.Col = ic.Loc.Col + uint(len(match))

	return &Token{
		Value: match,
		Loc:   ic.Loc,
		Kind:  NullToken,
	}, cur, true
}
