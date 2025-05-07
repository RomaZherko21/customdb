package lex

type TokenKind uint

const (
	KeywordToken TokenKind = iota
	SymbolToken
	IdentifierToken
	StringToken
	NumericToken
	BooleanToken
	NullToken
	MathOperatorToken
	LogicalOperatorToken
)

type Location struct {
	Line uint
	Col  uint
}

type Token struct {
	Value string
	Kind  TokenKind
	Loc   Location
}

func (t *Token) Equals(other *Token) bool {
	return t.Value == other.Value && t.Kind == other.Kind
}

type Cursor struct {
	Pointer uint
	Loc     Location
}
