package game

import "math/rand"

type GameVariant string

const (
	GameVariantTextInput      GameVariant = "text_input"
	GameVariantMultipleChoice GameVariant = "multiple_choice"
)

func assignGameVariant() GameVariant {
	if rand.Intn(2) == 0 {
		return GameVariantTextInput
	}
	return GameVariantMultipleChoice
}
