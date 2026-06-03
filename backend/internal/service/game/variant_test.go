package game

import "testing"

func TestAssignGameVariant(t *testing.T) {
	seen := make(map[GameVariant]bool)
	for i := 0; i < 100; i++ {
		seen[assignGameVariant()] = true
	}
	if !seen[GameVariantTextInput] || !seen[GameVariantMultipleChoice] {
		t.Fatal("assignGameVariant should produce both variants over many calls")
	}
}
