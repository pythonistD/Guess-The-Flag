package storage

import (
	"testing"

	"github.com/google/uuid"
)

func TestInMemoryGameStorage_GameVariant(t *testing.T) {
	storage := NewInMemoryGameStorage()
	gameID := uuid.New()
	variant := "multiple_choice"

	storage.InitStorageState(gameID, "rus", variant)

	got, err := storage.GetGameVariant(gameID)
	if err != nil {
		t.Fatalf("GetGameVariant returned error: %v", err)
	}
	if got != variant {
		t.Fatalf("expected variant %s, got %s", variant, got)
	}

	langCode, err := storage.GetGameLangCode(gameID)
	if err != nil {
		t.Fatalf("GetGameLangCode returned error: %v", err)
	}
	if langCode != "rus" {
		t.Fatalf("expected lang rus, got %s", langCode)
	}
}

func TestInMemoryGameStorage_GetUsedCountries(t *testing.T) {
	storage := NewInMemoryGameStorage()
	gameID := uuid.New()
	storage.InitStorageState(gameID, "rus", "text_input")

	if err := storage.SetCountry(gameID, 10); err != nil {
		t.Fatalf("SetCountry returned error: %v", err)
	}
	if err := storage.SetCountry(gameID, 20); err != nil {
		t.Fatalf("SetCountry returned error: %v", err)
	}

	used, err := storage.GetUsedCountries(gameID)
	if err != nil {
		t.Fatalf("GetUsedCountries returned error: %v", err)
	}
	if len(used) != 2 {
		t.Fatalf("expected 2 used countries, got %d", len(used))
	}
	if _, ok := used[10]; !ok {
		t.Fatal("expected country 10 to be used")
	}
	if _, ok := used[20]; !ok {
		t.Fatal("expected country 20 to be used")
	}
}
