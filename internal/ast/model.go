package ast

import "custom-database/internal/lex"

type Ast struct {
	Statements []*Statement
}

type AstKind uint

const (
	SelectKind AstKind = iota
	CreateTableKind
	InsertKind
)

type Statement struct {
	SelectStatement      *SelectStatement
	CreateTableStatement *CreateTableStatement
	InsertStatement      *InsertStatement
	Kind                 AstKind
}

type expressionKind uint

const (
	literalKind expressionKind = iota
)

type expression struct {
	literal *lex.Token
	kind    expressionKind
}
