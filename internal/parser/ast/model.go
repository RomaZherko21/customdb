package ast

import "custom-database/internal/parser/lex"

type Ast struct {
	Statements []*Statement
}

type AstKind uint

const (
	SelectKind AstKind = iota
	CreateTableKind
	InsertKind
	DropTableKind
)

type Statement struct {
	SelectStatement      *SelectStatement
	CreateTableStatement *CreateTableStatement
	InsertStatement      *InsertStatement
	DropTableStatement   *DropTableStatement
	Kind                 AstKind
}

type ExpressionKind uint

const (
	LiteralKind ExpressionKind = iota
)

type Expression struct {
	Literal *lex.Token
	Kind    ExpressionKind
}
