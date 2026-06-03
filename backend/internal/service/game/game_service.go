package game

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pythonistD/Guess-The-Flag/internal/service/game/algo"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pythonistD/Guess-The-Flag/internal/db/models"
	"github.com/pythonistD/Guess-The-Flag/internal/db/repo"
	"github.com/pythonistD/Guess-The-Flag/internal/schema"
	"github.com/pythonistD/Guess-The-Flag/internal/service/game/storage"
)

const (
	questionNumToPregenerate int = 5
)

type StartGameResult struct {
	GameId          uuid.UUID
	Variant         GameVariant
	MaxFlagsPerGame int
}

type Service struct {
	gamesRepo     *repo.GamesRepo
	questionsRepo *repo.QuestionsRepo
	answersRepo   *repo.AnswersRepo
	countriesRepo *repo.CountriesRepo

	gameStorage     storage.GameStorage
	countryStore    storage.CountryStorage
	maxFlagsPerGame int
}

func NewService(
	db *sqlx.DB,
	gameStore storage.GameStorage,
	countryStore storage.CountryStorage,
	maxFlagsPerGame int,
) *Service {
	if maxFlagsPerGame <= 0 {
		maxFlagsPerGame = 30
	}
	gamesRepo := repo.NewGamesRepo(db)
	questionsRepo := repo.NewQuestionsRepo(db)
	answersRepo := repo.NewAnswersRepo(db)
	countriesRepo := repo.NewCountriesRepo(db)
	return &Service{
		gamesRepo:       gamesRepo,
		questionsRepo:   questionsRepo,
		answersRepo:     answersRepo,
		countriesRepo:   countriesRepo,
		gameStorage:     gameStore,
		countryStore:    countryStore,
		maxFlagsPerGame: maxFlagsPerGame,
	}
}

func (s *Service) StartGame(ctx context.Context, langCode string) (*StartGameResult, error) {
	userId := ctx.Value("userId").(uuid.UUID)
	gameId := uuid.New()
	variant := assignGameVariant()

	s.gameStorage.InitStorageState(gameId, langCode, string(variant))

	gameModel := models.Game{
		GameId:       gameId,
		UserId:       userId,
		LanguageCode: langCode,
		GameVariant:  string(variant),
		StartedAt:    time.Now(),
	}

	err := s.gamesRepo.Start(ctx, &gameModel)
	if err != nil {
		return nil, fmt.Errorf("failed to start game: %w", err)
	}
	questions, err := s.generateQuestions(gameId)
	if err != nil {
		return nil, fmt.Errorf("failed to start game: %w", err)
	}
	err = s.gameStorage.SetQuestions(gameId, questions)
	if err != nil {
		return nil, fmt.Errorf("failed to start game: %w", err)
	}
	return &StartGameResult{
		GameId:          gameId,
		Variant:         variant,
		MaxFlagsPerGame: s.maxFlagsPerGame,
	}, nil
}

// EndGame deletes questions from gameState, returns game`s answers
func (s *Service) EndGame(ctx context.Context, gameId uuid.UUID) ([]models.QuestionWithAnswers, error) {
	langCode, err := s.gameStorage.GetGameLangCode(gameId)
	if err != nil {
		return nil, fmt.Errorf("failed to end game: %w", err)
	}
	err = s.gamesRepo.End(ctx, gameId)
	if err != nil {
		return nil, fmt.Errorf("failed to end game: %w", err)
	}
	questionsWithAnswer, err := s.questionsRepo.GetQuestionsWithAnswers(ctx, gameId, langCode)
	if err != nil {
		return nil, fmt.Errorf("failed to end game: %w", err)
	}
	// delete current game state
	err = s.gameStorage.DeleteGameSession(gameId)
	if err != nil {
		return nil, fmt.Errorf("failed to end game: %w", err)
	}
	return questionsWithAnswer, nil
}

func (s *Service) GetQuestion(ctx context.Context, gameId uuid.UUID) (any, error) {
	served, err := s.gameStorage.GetQuestionsServed(gameId)
	if err != nil {
		return nil, fmt.Errorf("failed to get question: %w", err)
	}
	if served >= s.maxFlagsPerGame {
		return nil, fmt.Errorf("game flag limit reached: %w", storage.GameLimitReached)
	}

	question, err := s.gameStorage.GetQuestion(gameId)
	if err != nil {
		if errors.Is(err, storage.GameIdError) {
			return nil, fmt.Errorf("game is already ended, please start new game: %w", err)
		}
		if errors.Is(err, storage.EmptyBuffer) || errors.Is(err, storage.GettingQuestionsError) {
			questionsGenerated, genErr := s.generateQuestions(gameId)
			if genErr != nil {
				return nil, fmt.Errorf("failed to get question from gameStateStorage: %w", genErr)
			}
			if len(questionsGenerated) == 0 {
				return nil, fmt.Errorf("game flag limit reached: %w", storage.GameLimitReached)
			}
			if setErr := s.gameStorage.SetQuestions(gameId, questionsGenerated); setErr != nil {
				return nil, fmt.Errorf("failed to save questions: %w", setErr)
			}
			question, err = s.gameStorage.GetQuestion(gameId)
			if err != nil {
				return nil, fmt.Errorf("failed to get question after refill: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to get question: %w", err)
		}
	}

	if question == nil {
		return nil, fmt.Errorf("failed to get question: question is nil")
	}

	step, err := s.gameStorage.IncrementQuestionsServed(gameId)
	if err != nil {
		return nil, fmt.Errorf("failed to update question progress: %w", err)
	}

	questionDB, err := s.saveQuestionToDB(ctx, question)
	if err != nil {
		return nil, fmt.Errorf("failed to save question to DB: %w", err)
	}

	variant, err := s.gameStorage.GetGameVariant(gameId)
	if err != nil {
		return nil, fmt.Errorf("failed to get game variant: %w", err)
	}

	base := schema.QuestionResp{
		QuestionText:    question.QuestionText,
		FlagSVG:         question.FlagSVG,
		QuestionID:      questionDB.QuestionId,
		Variant:         variant,
		Step:            step,
		MaxFlags:        s.maxFlagsPerGame,
		StepsCompleted:  step - 1,
	}

	if GameVariant(variant) == GameVariantMultipleChoice {
		langCode, err := s.gameStorage.GetGameLangCode(gameId)
		if err != nil {
			return nil, fmt.Errorf("failed to get game lang code: %w", err)
		}
		usedCountries, err := s.gameStorage.GetUsedCountries(gameId)
		if err != nil {
			return nil, fmt.Errorf("failed to get used countries: %w", err)
		}
		options, err := GenerateOptions(s.countryStore, question.CountryId, langCode, usedCountries)
		if err != nil {
			return nil, fmt.Errorf("failed to generate options: %w", err)
		}
		return &schema.MultipleChoiceQuestionResp{
			QuestionResp: base,
			Options:      options,
		}, nil
	}

	return &base, nil
}

func (s *Service) AnswerTheQuestion(
	ctx context.Context,
	gameId uuid.UUID,
	questionId uuid.UUID,
	req schema.AnswerReq,
) (*schema.AnswerResp, error) {
	countryId, err := s.getQuestionAnswer(ctx, questionId)
	if err != nil {
		return nil, fmt.Errorf("failed to get question answer: %w", err)
	}

	langCode, err := s.gameStorage.GetGameLangCode(gameId)
	if err != nil {
		return nil, fmt.Errorf("failed to get game lang code: %w", err)
	}

	variantStr, err := s.gameStorage.GetGameVariant(gameId)
	if err != nil {
		return nil, fmt.Errorf("failed to get game variant: %w", err)
	}
	variant := GameVariant(variantStr)

	isCorrect, answerText, selectedCountryId, err := s.gradeAnswer(ctx, variant, req, countryId, langCode)
	if err != nil {
		return nil, fmt.Errorf("failed to answer the question: %w", err)
	}

	err = s.SaveAnswer(ctx, questionId, answerText, isCorrect, selectedCountryId)
	if err != nil {
		return nil, fmt.Errorf("failed to answer the question: %w", err)
	}

	served, err := s.gameStorage.GetQuestionsServed(gameId)
	if err != nil {
		return nil, fmt.Errorf("failed to get game progress: %w", err)
	}

	return &schema.AnswerResp{
		IsCorrect:      isCorrect,
		Step:           served,
		MaxFlags:       s.maxFlagsPerGame,
		StepsCompleted: served,
		GameCompleted:  served >= s.maxFlagsPerGame,
	}, nil
}

func (s *Service) getQuestionAnswer(ctx context.Context, questionId uuid.UUID) (int, error) {
	q, err := s.questionsRepo.GetQuestion(ctx, questionId)
	if err != nil {
		return 0, fmt.Errorf("failed to get question from DB: %w", err)
	}
	return q.CountryId, nil

}

func (s *Service) saveQuestionToDB(ctx context.Context, questionStorage *storage.QuestionInStorage) (*models.Question, error) {
	questionDB := models.Question{
		GameId:     questionStorage.GameId,
		QuestionId: questionStorage.QuestionId,
		CountryId:  questionStorage.CountryId,
		CreatedAt:  questionStorage.CreatedAt,
	}
	err := s.questionsRepo.Create(ctx, &questionDB)
	if err != nil {
		return nil, fmt.Errorf("failed to save question to db: %w", err)
	}
	return &questionDB, nil
}

func (s *Service) isAnswerCorrect(ctx context.Context, countryId int, answer, langCode string) (bool, error) {
	var normalizedAnswer string
	normalizedAnswer = strings.ToLower(answer)
	normalizedAnswer = strings.TrimSpace(normalizedAnswer)
	country, err := s.countryStore.GetByID(countryId)
	if err != nil {
		return false, fmt.Errorf("failed to get country from DB: %w", err)
	}
	commonName := country.CommonCountryNames[langCode]
	dist := algo.WordDistance(commonName.NormalizedName, normalizedAnswer)
	if dist <= commonName.Threshold {
		return true, nil
	}
	allNames := country.AllCountryNames[langCode]
	for _, name := range allNames {
		dist = algo.WordDistance(name.NormalizedName, normalizedAnswer)
		if dist <= name.Threshold {
			return true, nil
		}
	}
	return false, nil
}

func (s *Service) SaveAnswer(ctx context.Context, questionId uuid.UUID, answerText string, isCorrect bool, selectedCountryId *int) error {
	answer := models.Answer{
		AnswerId:          uuid.New(),
		QuestionId:        questionId,
		AnswerText:        answerText,
		IsCorrect:         isCorrect,
		AnsweredAt:        time.Now(),
		SelectedCountryId: selectedCountryId,
	}
	err := s.answersRepo.Create(ctx, &answer)
	if err != nil {
		return fmt.Errorf("failed to save answer: %w", err)
	}
	return nil
}

func (s *Service) generateQuestions(gameId uuid.UUID) ([]storage.QuestionInStorage, error) {
	usedCountries, err := s.gameStorage.GetUsedCountries(gameId)
	if err != nil {
		return nil, fmt.Errorf("failed to generate questions: %w", err)
	}
	remaining := s.maxFlagsPerGame - len(usedCountries)
	if remaining <= 0 {
		return nil, nil
	}
	questionsGenerateNum := storage.QuestionsBuffSize
	if remaining < questionsGenerateNum {
		questionsGenerateNum = remaining
	}

	questions := make([]storage.QuestionInStorage, questionsGenerateNum)
	for i := 0; i < questionsGenerateNum; i++ {
		country, err := s.extractNewCountry(gameId)
		if err != nil {
			return nil, fmt.Errorf("failed to generate questions: %w", err)
		}
		langCode, err := s.gameStorage.GetGameLangCode(gameId)
		if err != nil {
			return nil, fmt.Errorf("failed to generate questions: %w", err)
		}
		question := storage.QuestionInStorage{
			QuestionId: uuid.New(),
			GameId:     gameId,
			CountryId:  country.Id,
			// todo: в зависимости от langcode должен выводить вопрос на выбранном языке
			QuestionText: "Guess The Flag",
			FlagSVG:      country.FlagSVG,
			CreatedAt:    time.Now(),
			Answer:       country.CommonCountryNames[langCode].Name,
		}
		questions[i] = question
	}
	return questions, nil
}

// FlagInfo используется для отладочного эндпоинта со всеми флагами.
type FlagInfo struct {
	CountryId int    `json:"country_id"`
	FlagSVG   string `json:"flag_svg"`
}

// GetAllFlags возвращает все флаги из in-memory хранилища стран.
// Предназначено для отладки на фронтенде, чтобы быстро проверить рендер всех SVG.
func (s *Service) GetAllFlags() []FlagInfo {
	countries := s.countryStore.GetAll()
	flags := make([]FlagInfo, 0, len(countries))
	for _, c := range countries {
		flags = append(flags, FlagInfo{
			CountryId: c.Id,
			FlagSVG:   c.FlagSVG,
		})
	}
	return flags
}

func (s *Service) extractNewCountry(gameId uuid.UUID) (storage.Country, error) {
	// todo: тут узкое место, если прошли много стран, то будет много итераций
	// так как использовано уже много, а осталось мало
	for i := 0; i < storage.NumberOfCountries; i++ {
		country, err := s.countryStore.GetRandom()
		if err != nil {
			return storage.Country{}, fmt.Errorf("failed to extract new country: %w", err)
		}
		if ok := s.gameStorage.IsCountryUsed(gameId, country.Id); !ok {
			err = s.gameStorage.SetCountry(gameId, country.Id)
			if err != nil {
				return storage.Country{}, fmt.Errorf("failed to extract new country: %w", err)
			}
			return country, nil
		}
	}
	return storage.Country{}, fmt.Errorf("failed to extract new country from game with id %s", gameId)
}
