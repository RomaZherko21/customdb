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

type CreateTableStatement struct {
	Name lex.Token
	Cols *[]*columnDefinition
}

type DropTableStatement struct {
	Table lex.Token
}

type InsertStatement struct {
	Table  lex.Token
	Values *[]*Expression
}

type SelectStatement struct {
	SelectedColumns []*Expression
	From            lex.Token
	Where           []*WhereClause
}

type WhereClause struct {
	Left  *Expression
	Right *Expression
	Op    lex.Token
}
