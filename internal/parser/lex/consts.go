package lex

type Keyword string

const (
	SelectKeyword Keyword = "select"
	FromKeyword   Keyword = "from"
	AsKeyword     Keyword = "as"
	TableKeyword  Keyword = "table"
	CreateKeyword Keyword = "create"
	InsertKeyword Keyword = "insert"
	IntoKeyword   Keyword = "into"
	ValuesKeyword Keyword = "values"
	IntKeyword    Keyword = "int"
	TextKeyword   Keyword = "text"
	DropKeyword   Keyword = "drop"
)

var keywords = []Keyword{
	SelectKeyword,
	InsertKeyword,
	ValuesKeyword,
	TableKeyword,
	CreateKeyword,
	FromKeyword,
	IntoKeyword,
	IntKeyword,
	TextKeyword,
	DropKeyword,
}

type Symbol string

const (
	SemicolonSymbol  Symbol = ";"
	AsteriskSymbol   Symbol = "*"
	CommaSymbol      Symbol = ","
	LeftparenSymbol  Symbol = "("
	RightparenSymbol Symbol = ")"
)

var symbols = []Symbol{
	CommaSymbol,
	LeftparenSymbol,
	RightparenSymbol,
	SemicolonSymbol,
	AsteriskSymbol,
}
