package lex

type keyword string

const (
	SelectKeyword keyword = "select"
	FromKeyword   keyword = "from"
	AsKeyword     keyword = "as"
	TableKeyword  keyword = "table"
	CreateKeyword keyword = "create"
	InsertKeyword keyword = "insert"
	IntoKeyword   keyword = "into"
	ValuesKeyword keyword = "values"
	IntKeyword    keyword = "int"
	TextKeyword   keyword = "text"
	DropKeyword   keyword = "drop"
)

var keywords = []keyword{
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

type symbol string

const (
	SemicolonSymbol  symbol = ";"
	AsteriskSymbol   symbol = "*"
	CommaSymbol      symbol = ","
	LeftparenSymbol  symbol = "("
	RightparenSymbol symbol = ")"
)

var symbols = []symbol{
	CommaSymbol,
	LeftparenSymbol,
	RightparenSymbol,
	SemicolonSymbol,
	AsteriskSymbol,
}
