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

const CREATE_TABLE SqlActionType = "CREATE_TABLE"
const DROP_TABLE SqlActionType = "DROP_TABLE"
const INSERT_INTO SqlActionType = "INSERT_INTO"
const SELECT_ACTION SqlActionType = "SELECT"

// keywords
type SqlKeywordType string

const TABLE SqlKeywordType = "TABLE"
const INTO SqlKeywordType = "INTO"
const SELECT SqlKeywordType = "SELECT"
const DROP SqlKeywordType = "DROP"
const INSERT SqlKeywordType = "INSERT"
const CREATE SqlKeywordType = "CREATE"

// TokenType = SYMBOL
const LPAREN = '('
const RPAREN = ')'
const COMMA = ','
const SEMICOLON = ';'
