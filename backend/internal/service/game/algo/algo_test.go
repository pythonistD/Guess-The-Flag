package algo

import (
	"fmt"
	"testing"
)

func TestWordDistance(t *testing.T) {
	tests := []struct {
		c1   string
		c2   string
		dist int
	}{
		{
			c1:   "Турция",
			c2:   "Тунис",
			dist: 3,
		},
		{
			c1:   "",
			c2:   "",
			dist: 0,
		},
		{
			c1:   "a",
			c2:   "",
			dist: 1,
		},
		{
			c1:   "",
			c2:   "a",
			dist: 1,
		},
		{
			c1:   "кот",
			c2:   "кот",
			dist: 0,
		},
		{
			c1:   "кот",
			c2:   "кит",
			dist: 1,
		},
		{
			c1:   "кот",
			c2:   "коты",
			dist: 1,
		},
		{
			c1:   "кот",
			c2:   "крот",
			dist: 1,
		},
		{
			c1:   "стол",
			c2:   "стул",
			dist: 1,
		},
		{
			c1:   "мама",
			c2:   "папа",
			dist: 2,
		},
		{
			c1:   "kitten",
			c2:   "sitting",
			dist: 3,
		},
		{
			c1:   "flaw",
			c2:   "lawn",
			dist: 2,
		},
		{
			c1:   "Москва",
			c2:   "Масква",
			dist: 1,
		},
		{
			c1:   "ёж",
			c2:   "еж",
			dist: 1,
		},
		{
			c1:   "привет",
			c2:   "превед",
			dist: 2,
		},
	}

	for _, tt := range tests {
		tName := fmt.Sprintf("w1:%s,w2:%s", tt.c1, tt.c2)
		t.Run(tName, func(t *testing.T) {
			distGot := WordDistance(tt.c1, tt.c2)
			if distGot != tt.dist {
				t.Errorf("dist expected: %d, got: %d", tt.dist, distGot)
			}
		})

	}

}
