package filldb

import (
	"os"
	"regexp"
	"testing"
)

func TestNormalizeSVG(t *testing.T) {
	filePath := "tst.svg"
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file: %s", filePath)
	}
	svg := string(data)
	normalized, err := normalizeSVG(svg)
	if err != nil {
		t.Fatalf("failed to normalize svg: %s", err)
	}

	hw := regexp.MustCompile(`\s*(width|height)="600"`)
	t.Log(normalized)
	if hw.MatchString(normalized) {
		t.Error("height and weight props haven`t deleted")
	}
}
