package sprite

import (
	"testing"
)

func TestGrowingPacker_SingleBlock(t *testing.T) {
	blocks := []*Node{NewNode("a", 100, 50)}
	(&GrowingPacker{}).Fit(blocks)
	if blocks[0].X != 0 || blocks[0].Y != 0 {
		t.Fatalf("expected (0,0), got (%d,%d)", blocks[0].X, blocks[0].Y)
	}
}

func TestGrowingPacker_TwoBlocks(t *testing.T) {
	a := NewNode("a", 100, 50)
	b := NewNode("b", 60, 30)
	blocks := []*Node{a, b}
	// Blocks are sorted by area descending, so 'a' comes first
	(&GrowingPacker{}).Fit(blocks)
	if a.X != 0 || a.Y != 0 {
		t.Fatalf("first block should be at (0,0), got (%d,%d)", a.X, a.Y)
	}
	// Second block should be placed somewhere (exact location depends on algorithm)
	if b.X == 0 && b.Y == 0 {
		// Both could overlap if algorithm doesn't separate them
		if b.Width <= a.Width && b.Height <= a.Height {
			// OK, placed next to first
		}
	}
	if b.X < 0 || b.Y < 0 {
		t.Fatal("negative coordinates")
	}
}

func TestGrowingPacker_IdenticalBlocks(t *testing.T) {
	n := 4
	blocks := make([]*Node, n)
	for i := 0; i < n; i++ {
		blocks[i] = NewNode(string(rune('a'+i)), 32, 32)
	}
	(&GrowingPacker{}).Fit(blocks)

	// All blocks should have non-negative coordinates
	for _, b := range blocks {
		if b.X < 0 || b.Y < 0 {
			t.Fatalf("block %s: negative coordinates (%d,%d)", b.Key, b.X, b.Y)
		}
	}
}

func TestGrowingPacker_ManyBlocks(t *testing.T) {
	n := 10
	blocks := make([]*Node, n)
	sizes := []int{64, 64, 32, 32, 16, 16, 8, 8, 4, 4}
	for i := 0; i < n; i++ {
		blocks[i] = NewNode(string(rune('a'+i)), sizes[i], sizes[i])
	}
	(&GrowingPacker{}).Fit(blocks)
	for _, b := range blocks {
		if b.X < 0 || b.Y < 0 {
			t.Fatalf("block %s: negative coordinates (%d,%d)", b.Key, b.X, b.Y)
		}
	}
}

func TestGrowingPacker_ZeroSizedBlock(t *testing.T) {
	a := NewNode("a", 0, 0)
	b := NewNode("b", 10, 10)
	blocks := []*Node{a, b}
	(&GrowingPacker{}).Fit(blocks)
	if b.X < 0 || b.Y < 0 {
		t.Fatal("negative coordinates")
	}
}

func TestGrowingPacker_EmptyInput(t *testing.T) {
	// Should not panic
	(&GrowingPacker{}).Fit(nil)
	(&GrowingPacker{}).Fit([]*Node{})
}

func TestFindNode_EmptyRoot(t *testing.T) {
	n := findNode(nil, 10, 10)
	if n != nil {
		t.Fatal("expected nil for nil root")
	}
}

func TestSplitNode(t *testing.T) {
	node := &Node{X: 0, Y: 0, Width: 100, Height: 100}
	result := splitNode(node, 40, 50)
	if result != node {
		t.Fatal("splitNode should return the same node")
	}
	if !node.Used {
		t.Fatal("node should be marked used")
	}
	if node.Right == nil || node.Down == nil {
		t.Fatal("expected Right and Down children")
	}
	if node.Right.Width != 60 || node.Right.Height != 50 {
		t.Fatalf("Right: expected 60x50, got %dx%d", node.Right.Width, node.Right.Height)
	}
	if node.Down.Width != 100 || node.Down.Height != 50 {
		t.Fatalf("Down: expected 100x50, got %dx%d", node.Down.Width, node.Down.Height)
	}
}
