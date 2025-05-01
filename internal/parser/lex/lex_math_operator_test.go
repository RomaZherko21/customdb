package lex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLexOperator(t *testing.T) {
	t.Run("invalid operator", func(t *testing.T) {
		input := "not an operator"
		cursor := Cursor{}

		_, _, isValid := lexMathOperator(input, cursor)

		require.False(t, isValid)
	})

	t.Run("multiple valid operators", func(t *testing.T) {
		tests := []struct {
			input       string
			want        string
			wantPointer uint
		}{
			{"=", string(EqualOperator), 1},
			{"<", string(LessThanOperator), 1},
			{">", string(GreaterThanOperator), 1},
			{"!=", string(NotEqualOperator), 2},
		}

		for _, tt := range tests {
			t.Run(tt.input, func(t *testing.T) {
				cursor := Cursor{}
				got, newCursor, isValid := lexMathOperator(tt.input, cursor)

				require.True(t, isValid)
				require.Equal(t, tt.want, got.Value)
				require.Equal(t, tt.wantPointer, newCursor.Pointer)
			})
		}
	})

	t.Run("multiple strange strings", func(t *testing.T) {
		tests := []struct {
			input       string
			want        string
			wantPointer uint
		}{
			{"====", string(EqualOperator), 1},
			{"<=-", string(LessThanOperator), 1},
			{">--=", string(GreaterThanOperator), 1},
			{"!===", string(NotEqualOperator), 2},
		}

		for _, tt := range tests {
			t.Run(tt.input, func(t *testing.T) {
				cursor := Cursor{}
				got, newCursor, isValid := lexMathOperator(tt.input, cursor)

				require.True(t, isValid)
				require.Equal(t, tt.want, got.Value)
				require.Equal(t, tt.wantPointer, newCursor.Pointer)
			})
		}
	})
}
