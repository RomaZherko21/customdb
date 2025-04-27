package lex

func lexKeyword(source string, ic Cursor) (*Token, Cursor, bool) {
	cur := ic
	keywords := []keyword{
		SelectKeyword,
		InsertKeyword,
		ValuesKeyword,
		TableKeyword,
		CreateKeyword,
		FromKeyword,
		IntoKeyword,
		IntKeyword,
		TextKeyword,
	}

	var options []string
	for _, k := range keywords {
		options = append(options, string(k))
	}

	match := longestMatch(source, ic, options)
	if match == "" {
		return nil, ic, false
	}

	cur.Pointer = ic.Pointer + uint(len(match))
	cur.Loc.Col = ic.Loc.Col + uint(len(match))

	return &Token{
		Value: match,
		Kind:  KeywordToken,
		Loc:   ic.Loc,
	}, cur, true
}
