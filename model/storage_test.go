package model

import (
	"testing"
)

func TestTreeOperations(t *testing.T) {
	tree := NewTree()
	r1 := tree.InsertBelow("", "Root 1")
	r2 := tree.InsertBelow(r1.ID, "Root 2")
	c1 := tree.AddChild(r1.ID, "Child 1.1")

	if len(tree.Roots) != 2 {
		t.Fatalf("expected 2 roots, got %d", len(tree.Roots))
	}
	if len(r1.Children) != 1 {
		t.Fatalf("expected 1 child for r1, got %d", len(r1.Children))
	}
	if c1.Text != "Child 1.1" {
		t.Fatalf("unexpected child text %s", c1.Text)
	}

	// Indent r2 to become child of r1
	ok := tree.Indent(r2.ID)
	if !ok {
		t.Fatalf("failed to indent r2")
	}
	if len(tree.Roots) != 1 {
		t.Fatalf("expected 1 root after indenting r2, got %d", len(tree.Roots))
	}
	if len(r1.Children) != 2 {
		t.Fatalf("expected 2 children for r1, got %d", len(r1.Children))
	}

	// Cycle task status
	tree.ToggleTask(c1.ID)
	if !c1.IsTask || c1.Status != StatusTodo {
		t.Fatalf("expected task status todo, got isTask=%v status=%s", c1.IsTask, c1.Status)
	}
	tree.CycleStatus(c1.ID)
	if c1.Status != StatusInProgress {
		t.Fatalf("expected status in_progress, got %s", c1.Status)
	}
	tree.CycleStatus(c1.ID)
	if c1.Status != StatusDone {
		t.Fatalf("expected status done, got %s", c1.Status)
	}
}

func TestMarkdownSerializeAndParse(t *testing.T) {
	tree := NewTree()
	r1 := tree.InsertBelow("", "Bullet 1")
	t1 := tree.InsertBelow(r1.ID, "Task Todo")
	tree.ToggleTask(t1.ID)

	t2 := tree.AddChild(r1.ID, "Subtask Done")
	tree.SetStatus(t2.ID, StatusDone)

	t3 := tree.AddChild(r1.ID, "Subtask In Progress")
	tree.SetStatus(t3.ID, StatusInProgress)

	r1.Folded = true

	md := SerializeMarkdown(tree)
	parsedTree := ParseMarkdown(md)

	if len(parsedTree.Roots) != 2 {
		t.Fatalf("expected 2 roots in parsed tree, got %d", len(parsedTree.Roots))
	}

	parsedR1 := parsedTree.Roots[0]
	if !parsedR1.Folded {
		t.Fatalf("expected parsed root to be folded")
	}
	if len(parsedR1.Children) != 2 {
		t.Fatalf("expected 2 children for parsed root, got %d", len(parsedR1.Children))
	}
	if parsedR1.Children[0].Status != StatusDone {
		t.Fatalf("expected status done for first child, got %s", parsedR1.Children[0].Status)
	}
	if parsedR1.Children[1].Status != StatusInProgress {
		t.Fatalf("expected status in_progress for second child, got %s", parsedR1.Children[1].Status)
	}
}

func TestEncryptionDecryption(t *testing.T) {
	plainText := "- [ ] First task\n  - [x] Subtask done\n"
	passphrase := "secret123"

	encrypted, err := encryptContent(plainText, passphrase)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	decrypted, err := decryptContent(encrypted, passphrase)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if decrypted != plainText {
		t.Fatalf("decrypted content mismatch: expected %q, got %q", plainText, decrypted)
	}

	_, err = decryptContent(encrypted, "wrongpass")
	if err == nil {
		t.Fatalf("expected error when decrypting with wrong password")
	}
}
