package lex

// lexCharacterDelimited looks through a source string starting at the
// given cursor to find a start- and end- delimiter. The delimiter can
// be escaped be preceeding the delimiter with itself.
func lexCharacterDelimited(source string, ic Cursor, delimiter byte) (*Token, Cursor, bool) {
	cur := ic

	if len(source[cur.Pointer:]) == 0 {
		return nil, ic, false
	}

	if source[cur.Pointer] != delimiter {
		return nil, ic, false
	}

	cur.Loc.Col++
	cur.Pointer++

	var value []byte
	for ; cur.Pointer < uint(len(source)); cur.Pointer++ {
		c := source[cur.Pointer]

		if c == delimiter {
			// SQL escapes are via double characters, not backslash.
			if cur.Pointer+1 >= uint(len(source)) || source[cur.Pointer+1] != delimiter {
				cur.Pointer++
				cur.Loc.Col++
				return &Token{
					Value: string(value),
					Loc:   ic.Loc,
					Kind:  StringToken,
				}, cur, true
			}
			value = append(value, delimiter)
			cur.Pointer++
			cur.Loc.Col++
		}

		value = append(value, c)
		cur.Loc.Col++
	}

	return nil, ic, false
}
