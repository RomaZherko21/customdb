package lex

type Keyword string

const (
	// Statements
	CreateKeyword Keyword = "create"
	DropKeyword   Keyword = "drop"
	SelectKeyword Keyword = "select"
	InsertKeyword Keyword = "insert"
	// Keywords
	FromKeyword   Keyword = "from"
	TableKeyword  Keyword = "table"
	AsKeyword     Keyword = "as"
	IntoKeyword   Keyword = "into"
	ValuesKeyword Keyword = "values"
	WhereKeyword  Keyword = "where"
	LimitKeyword  Keyword = "limit"
	OffsetKeyword Keyword = "offset"
	// Datatypes
	IntKeyword         Keyword = "int"
	TextKeyword        Keyword = "text"
	BooleanTypeKeyword Keyword = "boolean"
	NullKeyword        Keyword = "null"
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
	WhereKeyword,
	BooleanTypeKeyword,
	NullKeyword,
	LimitKeyword,
	OffsetKeyword,
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

type BooleanKeyword string

const (
	TrueKeyword  BooleanKeyword = "true"
	FalseKeyword BooleanKeyword = "false"
)

var booleanKeywords = []BooleanKeyword{
	TrueKeyword,
	FalseKeyword,
}

type MathOperator string

const (
	EqualOperator       MathOperator = "="
	NotEqualOperator    MathOperator = "!="
	GreaterThanOperator MathOperator = ">"
	LessThanOperator    MathOperator = "<"
)

var mathOperators = []MathOperator{
	EqualOperator,
	NotEqualOperator,
	GreaterThanOperator,
	LessThanOperator,
}

type LogicalOperator string

const (
	AndOperator LogicalOperator = "and"
	OrOperator  LogicalOperator = "or"
)

var logicalOperators = []LogicalOperator{
	AndOperator,
	OrOperator,
}
