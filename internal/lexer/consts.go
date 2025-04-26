package lexer

// Типы токенов
type TokenType string

const (
	KEYWORD    TokenType = "KEYWORD"
	IDENTIFIER TokenType = "IDENTIFIER"
	TYPE       TokenType = "TYPE"
	SYMBOL     TokenType = "SYMBOL"
)

type KeywordType string

const CREATE_TABLE KeywordType = "CREATE_TABLE"
const INSERT_INTO KeywordType = "INSERT_INTO"
const SELECT KeywordType = "SELECT"

// TokenType = SYMBOL
const LPAREN = '('
const RPAREN = ')'
const COMMA = ','
const SEMICOLON = ';'

// ColumnType
type ColumnType string

const INT ColumnType = "INT"
const TEXT ColumnType = "TEXT"
