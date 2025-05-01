package lex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLexLogicalOperator(t *testing.T) {
	t.Run("invalid operator", func(t *testing.T) {
		input := "not a logical operator"
		cursor := Cursor{}

		_, _, isValid := lexLogicalOperator(input, cursor)

		require.False(t, isValid)
	})

	t.Run("multiple valid operators", func(t *testing.T) {
		tests := []struct {
			input       string
			want        string
			wantPointer uint
		}{
			{"AND", string(AndOperator), 3},
			{"OR", string(OrOperator), 2},
		}

		for _, tt := range tests {
			t.Run(tt.input, func(t *testing.T) {
				cursor := Cursor{}
				got, newCursor, isValid := lexLogicalOperator(tt.input, cursor)

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
			{"ANDAND", string(AndOperator), 3},
			{"OROR", string(OrOperator), 2},
		}

		for _, tt := range tests {
			t.Run(tt.input, func(t *testing.T) {
				cursor := Cursor{}
				got, newCursor, isValid := lexLogicalOperator(tt.input, cursor)

				require.True(t, isValid)
				require.Equal(t, tt.want, got.Value)
				require.Equal(t, tt.wantPointer, newCursor.Pointer)
			})
		}
	})
}
