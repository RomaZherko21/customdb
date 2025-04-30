package parser

import (
	"custom-database/internal/parser/ast"
)

type ParserService interface {
	Parse(query string) (*ast.Ast, error)
}

type Parser struct {
}

func NewParser() ParserService {
	return &Parser{}
}

func (p *Parser) Parse(query string) (*ast.Ast, error) {
	astService := ast.NewAst()
	return astService.Parse(query)
}
