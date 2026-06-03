package game

import (
	"context"
	"fmt"

	"github.com/pythonistD/Guess-The-Flag/internal/schema"
)

func (s *Service) gradeAnswer(
	ctx context.Context,
	variant GameVariant,
	req schema.AnswerReq,
	correctCountryId int,
	langCode string,
) (isCorrect bool, answerText string, selectedCountryId *int, err error) {
	if req.Skipped {
		return false, "", nil, nil
	}

	switch variant {
	case GameVariantTextInput:
		if req.Answer == "" {
			return false, "", nil, fmt.Errorf("answer is required for text_input variant")
		}
		isCorrect, err = s.isAnswerCorrect(ctx, correctCountryId, req.Answer, langCode)
		if err != nil {
			return false, "", nil, err
		}
		return isCorrect, req.Answer, nil, nil

	case GameVariantMultipleChoice:
		if req.SelectedCountry <= 0 {
			return false, "", nil, fmt.Errorf("selected_country is required for multiple_choice variant")
		}
		isCorrect = req.SelectedCountry == correctCountryId
		country, err := s.countryStore.GetByID(req.SelectedCountry)
		if err != nil {
			return false, "", nil, fmt.Errorf("failed to get selected country: %w", err)
		}
		commonName, ok := country.CommonCountryNames[langCode]
		if !ok {
			return false, "", nil, fmt.Errorf("no common name for country %d in lang %s", req.SelectedCountry, langCode)
		}
		selectedId := req.SelectedCountry
		return isCorrect, commonName.Name, &selectedId, nil

	default:
		return false, "", nil, fmt.Errorf("unknown game variant: %s", variant)
	}
}
