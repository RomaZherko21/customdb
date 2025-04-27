package lexer

// Типы токенов
type TokenType string

const (
	KEYWORD    TokenType = "KEYWORD"
	IDENTIFIER TokenType = "IDENTIFIER"
	TYPE       TokenType = "TYPE"
	SYMBOL     TokenType = "SYMBOL"
)

// actions
type SqlActionType string

const (
	CREATE_TABLE  SqlActionType = "CREATE_TABLE"
	DROP_TABLE    SqlActionType = "DROP_TABLE"
	INSERT_INTO   SqlActionType = "INSERT_INTO"
	SELECT_ACTION SqlActionType = "SELECT"
)

// keywords
type SqlKeywordType string

const (
	TABLE  SqlKeywordType = "TABLE"
	INTO   SqlKeywordType = "INTO"
	SELECT SqlKeywordType = "SELECT"
	DROP   SqlKeywordType = "DROP"
	INSERT SqlKeywordType = "INSERT"
	CREATE SqlKeywordType = "CREATE"
)

// TokenType = SYMBOL
const (
	LPAREN    = '('
	RPAREN    = ')'
	COMMA     = ','
	SEMICOLON = ';'
)
