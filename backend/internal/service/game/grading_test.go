package game

import (
	"context"
	"testing"

	"github.com/pythonistD/Guess-The-Flag/internal/schema"
)

func TestGradeAnswer_Skipped(t *testing.T) {
	s := &Service{countryStore: newMockCountryStorage()}
	isCorrect, answerText, selectedId, err := s.gradeAnswer(
		context.Background(),
		GameVariantTextInput,
		schema.AnswerReq{Skipped: true},
		1,
		"rus",
	)
	if err != nil {
		t.Fatalf("gradeAnswer returned error: %v", err)
	}
	if isCorrect || answerText != "" || selectedId != nil {
		t.Fatalf("expected skipped result false/empty/nil, got %v %q %v", isCorrect, answerText, selectedId)
	}
}

func TestGradeAnswer_MultipleChoice_Correct(t *testing.T) {
	s := &Service{countryStore: newMockCountryStorage()}
	isCorrect, answerText, selectedId, err := s.gradeAnswer(
		context.Background(),
		GameVariantMultipleChoice,
		schema.AnswerReq{SelectedCountry: 1},
		1,
		"rus",
	)
	if err != nil {
		t.Fatalf("gradeAnswer returned error: %v", err)
	}
	if !isCorrect {
		t.Fatal("expected correct answer")
	}
	if answerText != "Франция" {
		t.Fatalf("expected answer text Франция, got %s", answerText)
	}
	if selectedId == nil || *selectedId != 1 {
		t.Fatalf("expected selected country id 1, got %v", selectedId)
	}
}

func TestGradeAnswer_MultipleChoice_Incorrect(t *testing.T) {
	s := &Service{countryStore: newMockCountryStorage()}
	isCorrect, _, _, err := s.gradeAnswer(
		context.Background(),
		GameVariantMultipleChoice,
		schema.AnswerReq{SelectedCountry: 2},
		1,
		"rus",
	)
	if err != nil {
		t.Fatalf("gradeAnswer returned error: %v", err)
	}
	if isCorrect {
		t.Fatal("expected incorrect answer")
	}
}

func TestGradeAnswer_TextInput(t *testing.T) {
	s := &Service{countryStore: newMockCountryStorage()}
	isCorrect, answerText, selectedId, err := s.gradeAnswer(
		context.Background(),
		GameVariantTextInput,
		schema.AnswerReq{Answer: "Франция"},
		1,
		"rus",
	)
	if err != nil {
		t.Fatalf("gradeAnswer returned error: %v", err)
	}
	if !isCorrect {
		t.Fatal("expected correct text answer")
	}
	if answerText != "Франция" {
		t.Fatalf("expected answer text Франция, got %s", answerText)
	}
	if selectedId != nil {
		t.Fatalf("expected nil selected country id, got %v", selectedId)
	}
}
