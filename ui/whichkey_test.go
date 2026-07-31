package ui

import (
	"testing"
)

func TestWhichKeyOptions(t *testing.T) {
	wk := NewWhichKeyModel()
	allBindings := GetAllKeyBindings()

	// Initial leader key
	wk.Active = true
	wk.PrefixKeys = []string{" "}

	title, options := wk.GetOptions(allBindings)
	if title != " Space " {
		t.Fatalf("expected title ' Space ', got %q", title)
	}
	if len(options) == 0 {
		t.Fatalf("expected non-empty options for space leader key")
	}

	// Check group options presence (+bullets, +tasks, +folds)
	foundBulletsGroup := false
	for _, opt := range options {
		if opt.Key == "b" && opt.IsGroup && opt.Label == "+bullets" {
			foundBulletsGroup = true
		}
	}
	if !foundBulletsGroup {
		t.Fatalf("expected +bullets group under leader key")
	}

	// Sub-prefix Space > Bullets
	wk.PrefixKeys = []string{" ", "b"}
	title, subOptions := wk.GetOptions(allBindings)
	if title != " Space > b " {
		t.Fatalf("expected title ' Space > b ', got %q", title)
	}

	foundNewBelow := false
	for _, opt := range subOptions {
		if opt.Key == "n" && !opt.IsGroup && opt.Label == "new bullet below" {
			foundNewBelow = true
		}
	}
	if !foundNewBelow {
		t.Fatalf("expected 'new bullet below' under Space > b")
	}
}
