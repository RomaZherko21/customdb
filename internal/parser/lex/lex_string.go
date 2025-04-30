package lex

func lexString(source string, ic Cursor) (*Token, Cursor, bool) {
	return lexCharacterDelimited(source, ic, '\'')
}
