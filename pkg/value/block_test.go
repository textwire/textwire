package value

import "testing"

func TestBlockAppend(t *testing.T) {
	t.Run("appends chunks in order", func(t *testing.T) {
		first := &Text{Val: "a"}
		second := &Text{Val: "b"}

		block := NewBlock(2)
		block.Append(first)
		block.Append(second)

		if len(block.chunks) != 2 {
			t.Fatalf("expected 2 chunks, got %d", len(block.chunks))
		}

		if block.chunks[0] != first || block.chunks[1] != second {
			t.Fatal("expected chunks in append order")
		}
	})

	t.Run("break chunk sets only break flag", func(t *testing.T) {
		block := NewBlock(1)
		block.Append(&Break{})

		if !block.hasBreak {
			t.Fatal("expected hasBreak")
		}

		if block.hasContinue {
			t.Fatal("unexpected hasContinue")
		}
	})

	t.Run("continue chunk sets only continue flag", func(t *testing.T) {
		block := NewBlock(1)
		block.Append(&Continue{})

		if !block.hasContinue {
			t.Fatal("expected hasContinue")
		}

		if block.hasBreak {
			t.Fatal("unexpected hasBreak")
		}
	})

	t.Run("nested block with break sets break flag", func(t *testing.T) {
		inner := NewBlock(1)
		inner.Append(&Break{})

		block := NewBlock(1)
		block.Append(inner)

		if !block.hasBreak {
			t.Fatal("expected hasBreak from nested block")
		}
	})

	t.Run("later chunk keeps earlier flag", func(t *testing.T) {
		block := NewBlock(2)
		block.Append(&Break{})
		block.Append(&Text{Val: "a"})

		if !block.hasBreak {
			t.Fatal("expected hasBreak to stay")
		}
	})
}

func TestBlockHasControl(t *testing.T) {
	t.Run("finds break in nested block", func(t *testing.T) {
		inner := NewBlock(1)
		inner.Append(&Break{})

		outer := NewBlock(1)
		outer.Append(&Text{Val: "a"})
		outer.Append(inner)

		if !outer.HasControl(BREAK_VAL) {
			t.Fatal("expected break in nested block")
		}

		if outer.HasControl(CONTINUE_VAL) {
			t.Fatal("unexpected continue in block")
		}
	})

	t.Run("keeps break after later chunks", func(t *testing.T) {
		block := NewBlock(2)
		block.Append(&Break{})
		block.Append(&Text{Val: "a"})

		if !block.HasControl(BREAK_VAL) {
			t.Fatal("expected break to stay after a later chunk")
		}
	})

	t.Run("keeps continue after later chunks", func(t *testing.T) {
		block := NewBlock(2)
		block.Append(&Continue{})
		block.Append(&Text{Val: "a"})

		if !block.HasControl(CONTINUE_VAL) {
			t.Fatal("expected continue to stay after a later chunk")
		}
	})

	t.Run("false without control values", func(t *testing.T) {
		block := NewBlock(1)
		block.Append(&Text{Val: "a"})
		block.Append(&Text{Val: "b"})

		if block.HasControl(BREAK_VAL) || block.HasControl(CONTINUE_VAL) {
			t.Fatal("expected no control values")
		}
	})
}
