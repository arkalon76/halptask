package ui

import (
	"strings"
	"testing"
)

func TestNoKeybindingPrefixCollisions(t *testing.T) {
	bindings := GetAllKeyBindings()

	for i, b1 := range bindings {
		for j, b2 := range bindings {
			if i == j {
				continue
			}

			// Check if b1.Keys is a proper prefix of b2.Keys
			if len(b1.Keys) < len(b2.Keys) {
				isPrefix := true
				for k := 0; k < len(b1.Keys); k++ {
					if b1.Keys[k] != b2.Keys[k] {
						isPrefix = false
						break
					}
				}

				if isPrefix {
					t.Errorf("Keybinding collision detected: %q (%v) is a prefix of %q (%v). The shorter binding will hijack keystrokes!",
						b1.KeyString, strings.Join(b1.Keys, " "),
						b2.KeyString, strings.Join(b2.Keys, " "),
					)
				}
			}
		}
	}
}
