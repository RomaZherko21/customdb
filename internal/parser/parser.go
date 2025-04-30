package parser

import (
	"custom-database/internal/parser/ast"
)

type ParserService interface {
	Parse(query string) (*ast.Ast, error)
}

type parser struct {
}

func NewParser() ParserService {
	return &parser{}
}

func (p *parser) Parse(query string) (*ast.Ast, error) {
	astService := ast.NewAst()
	return astService.Parse(query)
}
