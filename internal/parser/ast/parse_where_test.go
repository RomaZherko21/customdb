package ast

import (
	"custom-database/internal/parser/lex"
	"testing"
)

func TestMapWhereTokens(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []*lex.Token
		expected []*tokenWithPrior
	}{
		{
			name: "простой случай с идентификатором",
			tokens: []*lex.Token{
				{Kind: lex.IdentifierToken, Value: "name"},
			},
			expected: []*tokenWithPrior{
				{Token: &lex.Token{Kind: lex.IdentifierToken, Value: "name"}, Priority: 41},
			},
		},
		{
			name: "выражение с операторами",
			tokens: []*lex.Token{
				{Kind: lex.IdentifierToken, Value: "age"},
				{Kind: lex.SymbolToken, Value: ">"},
				{Kind: lex.NumericToken, Value: "18"},
			},
			expected: []*tokenWithPrior{
				{Token: &lex.Token{Kind: lex.IdentifierToken, Value: "age"}, Priority: 41},
				{Token: &lex.Token{Kind: lex.SymbolToken, Value: ">"}, Priority: 40},
				{Token: &lex.Token{Kind: lex.NumericToken, Value: "18"}, Priority: 42},
			},
		},
		{
			name: "выражение со скобками",
			tokens: []*lex.Token{
				{Kind: lex.SymbolToken, Value: "("},
				{Kind: lex.IdentifierToken, Value: "name"},
				{Kind: lex.SymbolToken, Value: "="},
				{Kind: lex.StringToken, Value: "'John'"},
				{Kind: lex.SymbolToken, Value: ")"},
			},
			expected: []*tokenWithPrior{
				{Token: &lex.Token{Kind: lex.IdentifierToken, Value: "name"}, Priority: 141},
				{Token: &lex.Token{Kind: lex.SymbolToken, Value: "="}, Priority: 140},
				{Token: &lex.Token{Kind: lex.StringToken, Value: "'John'"}, Priority: 142},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := addTokensPriority(tt.tokens)
			if len(result) != len(tt.expected) {
				t.Errorf("ожидалось %d токенов, получено %d", len(tt.expected), len(result))
				return
			}

			for i, expected := range tt.expected {
				if result[i].Token.Kind != expected.Token.Kind {
					t.Errorf("токен %d: ожидался тип %v, получен %v", i, expected.Token.Kind, result[i].Token.Kind)
				}
				if result[i].Token.Value != expected.Token.Value {
					t.Errorf("токен %d: ожидалось значение %v, получено %v", i, expected.Token.Value, result[i].Token.Value)
				}
				if result[i].Priority != expected.Priority {
					t.Errorf("токен %d: ожидался приоритет %d, получен %d", i, expected.Priority, result[i].Priority)
				}
			}
		})
	}
}

func TestGetPriority(t *testing.T) {
	tests := []struct {
		name        string
		token       *lex.Token
		subPriority uint
		identPrior  uint
		expected    uint
	}{
		{
			name:        "идентификатор",
			token:       &lex.Token{Kind: lex.IdentifierToken, Value: "name"},
			subPriority: 0,
			identPrior:  0,
			expected:    40,
		},
		{
			name:        "оператор OR",
			token:       &lex.Token{Kind: lex.SymbolToken, Value: string(lex.OrOperator)},
			subPriority: 100,
			identPrior:  0,
			expected:    120,
		},
		{
			name:        "оператор AND",
			token:       &lex.Token{Kind: lex.SymbolToken, Value: string(lex.AndOperator)},
			subPriority: 100,
			identPrior:  0,
			expected:    130,
		},
		{
			name:        "оператор сравнения",
			token:       &lex.Token{Kind: lex.SymbolToken, Value: string(lex.EqualOperator)},
			subPriority: 100,
			identPrior:  0,
			expected:    140,
		},
		{
			name:        "неизвестный оператор",
			token:       &lex.Token{Kind: lex.SymbolToken, Value: "unknown"},
			subPriority: 100,
			identPrior:  0,
			expected:    1100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getPriority(tt.token, tt.subPriority, tt.identPrior)
			if result != tt.expected {
				t.Errorf("ожидался приоритет %d, получен %d", tt.expected, result)
			}
		})
	}
}

func TestParseTree(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []*tokenWithPrior
		expected *WhereClause
	}{
		{
			name:     "пустой список токенов",
			tokens:   []*tokenWithPrior{},
			expected: nil,
		},
		{
			name: "один токен",
			tokens: []*tokenWithPrior{
				{Token: &lex.Token{Kind: lex.IdentifierToken, Value: "name"}, Priority: 40},
			},
			expected: &WhereClause{
				Token: &lex.Token{Kind: lex.IdentifierToken, Value: "name"},
			},
		},
		{
			name: "простое сравнение",
			tokens: []*tokenWithPrior{
				{Token: &lex.Token{Kind: lex.IdentifierToken, Value: "age"}, Priority: 41},
				{Token: &lex.Token{Kind: lex.SymbolToken, Value: ">"}, Priority: 40},
				{Token: &lex.Token{Kind: lex.NumericToken, Value: "18"}, Priority: 42},
			},
			expected: &WhereClause{
				Token: &lex.Token{Kind: lex.SymbolToken, Value: ">"},
				Left: &WhereClause{
					Token: &lex.Token{Kind: lex.IdentifierToken, Value: "age"},
				},
				Right: &WhereClause{
					Token: &lex.Token{Kind: lex.NumericToken, Value: "18"},
				},
			},
		},
		{
			name: "сложное выражение с AND",
			tokens: []*tokenWithPrior{
				{Token: &lex.Token{Kind: lex.IdentifierToken, Value: "age"}, Priority: 41},
				{Token: &lex.Token{Kind: lex.SymbolToken, Value: ">"}, Priority: 40},
				{Token: &lex.Token{Kind: lex.NumericToken, Value: "18"}, Priority: 42},
				{Token: &lex.Token{Kind: lex.SymbolToken, Value: "AND"}, Priority: 30},
				{Token: &lex.Token{Kind: lex.IdentifierToken, Value: "name"}, Priority: 41},
				{Token: &lex.Token{Kind: lex.SymbolToken, Value: "="}, Priority: 40},
				{Token: &lex.Token{Kind: lex.StringToken, Value: "'John'"}, Priority: 42},
			},
			expected: &WhereClause{
				Token: &lex.Token{Kind: lex.SymbolToken, Value: "AND"},
				Left: &WhereClause{
					Token: &lex.Token{Kind: lex.SymbolToken, Value: ">"},
					Left: &WhereClause{
						Token: &lex.Token{Kind: lex.IdentifierToken, Value: "age"},
					},
					Right: &WhereClause{
						Token: &lex.Token{Kind: lex.NumericToken, Value: "18"},
					},
				},
				Right: &WhereClause{
					Token: &lex.Token{Kind: lex.SymbolToken, Value: "="},
					Left: &WhereClause{
						Token: &lex.Token{Kind: lex.IdentifierToken, Value: "name"},
					},
					Right: &WhereClause{
						Token: &lex.Token{Kind: lex.StringToken, Value: "'John'"},
					},
				},
			},
		},
		{
			name: "выражение со скобками",
			tokens: []*tokenWithPrior{
				{Token: &lex.Token{Kind: lex.IdentifierToken, Value: "age"}, Priority: 41},
				{Token: &lex.Token{Kind: lex.SymbolToken, Value: ">"}, Priority: 40},
				{Token: &lex.Token{Kind: lex.NumericToken, Value: "18"}, Priority: 42},
				{Token: &lex.Token{Kind: lex.SymbolToken, Value: "AND"}, Priority: 30},
				{Token: &lex.Token{Kind: lex.IdentifierToken, Value: "name"}, Priority: 141},
				{Token: &lex.Token{Kind: lex.SymbolToken, Value: "="}, Priority: 140},
				{Token: &lex.Token{Kind: lex.StringToken, Value: "'John'"}, Priority: 142},
				{Token: &lex.Token{Kind: lex.SymbolToken, Value: "OR"}, Priority: 120},
				{Token: &lex.Token{Kind: lex.IdentifierToken, Value: "name"}, Priority: 141},
				{Token: &lex.Token{Kind: lex.SymbolToken, Value: "="}, Priority: 140},
				{Token: &lex.Token{Kind: lex.StringToken, Value: "'Jane'"}, Priority: 142},
			},
			expected: &WhereClause{
				Token: &lex.Token{Kind: lex.SymbolToken, Value: "AND"},
				Left: &WhereClause{
					Token: &lex.Token{Kind: lex.SymbolToken, Value: ">"},
					Left: &WhereClause{
						Token: &lex.Token{Kind: lex.IdentifierToken, Value: "age"},
					},
					Right: &WhereClause{
						Token: &lex.Token{Kind: lex.NumericToken, Value: "18"},
					},
				},
				Right: &WhereClause{
					Token: &lex.Token{Kind: lex.SymbolToken, Value: "OR"},
					Left: &WhereClause{
						Token: &lex.Token{Kind: lex.SymbolToken, Value: "="},
						Left: &WhereClause{
							Token: &lex.Token{Kind: lex.IdentifierToken, Value: "name"},
						},
						Right: &WhereClause{
							Token: &lex.Token{Kind: lex.StringToken, Value: "'John'"},
						},
					},
					Right: &WhereClause{
						Token: &lex.Token{Kind: lex.SymbolToken, Value: "="},
						Left: &WhereClause{
							Token: &lex.Token{Kind: lex.IdentifierToken, Value: "name"},
						},
						Right: &WhereClause{
							Token: &lex.Token{Kind: lex.StringToken, Value: "'Jane'"},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseTree(tt.tokens)
			if !compareTrees(result, tt.expected) {
				t.Errorf("деревья не совпадают")
			}
		})
	}
}

func compareTrees(a, b *WhereClause) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	if a.Token.Kind != b.Token.Kind || a.Token.Value != b.Token.Value {
		return false
	}

	return compareTrees(a.Left, b.Left) && compareTrees(a.Right, b.Right)
}
