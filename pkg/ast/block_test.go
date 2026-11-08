package ast

import (
	"testing"

	"codeberg.org/textwire/textwire/v5/pkg/token"
)

func TestToDefaultPassDir(t *testing.T) {
	text := func(lit string) Chunk {
		return NewText(token.Token{Lit: lit})
	}

	pass := func() Chunk {
		return NewPassDir(token.Token{}, &StrExpr{Val: "name"})
	}

	cases := []struct {
		name   string
		chunks []Chunk
		// expectChunks is nil when ToDefaultPassDir should return nil
		expectChunks []string
	}{
		{
			name:   "only @pass returns nil",
			chunks: []Chunk{pass()},
		},
		{
			name:   "only blank text returns nil",
			chunks: []Chunk{text("\n  ")},
		},
		{
			name:   "blank text around @pass returns nil",
			chunks: []Chunk{text("\n    "), pass(), text("\n")},
		},
		{
			name:   "empty last text after @pass returns nil",
			chunks: []Chunk{pass(), text("")},
		},
		{
			name:         "content before empty last text is kept",
			chunks:       []Chunk{text("x"), text("")},
			expectChunks: []string{"x", ""},
		},
		{
			name:         "empty first text before content is kept",
			chunks:       []Chunk{text(""), text("x")},
			expectChunks: []string{"", "x"},
		},
		{
			name:         "blank text next to content is kept",
			chunks:       []Chunk{text("\n"), text("x"), text("\n")},
			expectChunks: []string{"\n", "x", "\n"},
		},
		{
			name:         "@pass is removed from content",
			chunks:       []Chunk{pass(), text("x")},
			expectChunks: []string{"x"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			block := NewBlock(token.Token{})
			block.Chunks = tc.chunks

			passDir := block.ToDefaultPassDir("comp")

			if tc.expectChunks == nil {
				if passDir != nil {
					t.Fatalf("expect nil, got pass with %d chunks", len(passDir.Block.Chunks))
				}
				return
			}

			if passDir == nil {
				t.Fatal("expect pass, got nil")
			}

			if len(passDir.Block.Chunks) != len(tc.expectChunks) {
				t.Fatalf(
					"expect %d chunks, got %d",
					len(tc.expectChunks),
					len(passDir.Block.Chunks),
				)
			}

			for i, want := range tc.expectChunks {
				text, ok := passDir.Block.Chunks[i].(*Text)
				if !ok {
					t.Fatalf("chunk %d is not *Text, got %T", i, passDir.Block.Chunks[i])
				}

				if text.Token.Lit != want {
					t.Fatalf("chunk %d: expect %q, got %q", i, want, text.Token.Lit)
				}
			}
		})
	}
}
