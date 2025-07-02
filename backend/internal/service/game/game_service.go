package game

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pythonistD/Guess-The-Flag/internal/db/models"
	"github.com/pythonistD/Guess-The-Flag/internal/db/repo"
	"github.com/pythonistD/Guess-The-Flag/internal/schema"
	"github.com/pythonistD/Guess-The-Flag/internal/service/game/storage"
	"strings"
	"time"
)

const (
	questionNumToPregenerate int = 5
)

type Service struct {
	gamesRepo     *repo.GamesRepo
	questionsRepo *repo.QuestionsRepo
	answersRepo   *repo.AnswersRepo
	countriesRepo *repo.CountriesRepo

	gameStorage  storage.GameStorage
	countryStore storage.CountryStorage
}

func NewService(db *sqlx.DB, gameStore storage.GameStorage, countryStore storage.CountryStorage) *Service {
	gamesRepo := repo.NewGamesRepo(db)
	questionsRepo := repo.NewQuestionsRepo(db)
	answersRepo := repo.NewAnswersRepo(db)
	countriesRepo := repo.NewCountriesRepo(db)
	return &Service{
		gamesRepo:     gamesRepo,
		questionsRepo: questionsRepo,
		answersRepo:   answersRepo,
		countriesRepo: countriesRepo,
		gameStorage:   gameStore,
		countryStore:  countryStore,
	}
}

func (s *Service) StartGame(ctx context.Context) (uuid.UUID, error) {
	userId := ctx.Value("userId").(uuid.UUID)
	gameId := uuid.New()

	gameModel := models.Game{
		GameId:    gameId,
		UserId:    userId,
		StartedAt: time.Now(),
	}

	err := s.gamesRepo.Start(ctx, &gameModel)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to start game: %w", err)
	}
	questions, err := s.generateQuestions(0, questionNumToPregenerate, gameId)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to start game: %w", err)
	}
	err = s.gameStorage.SetQuestions(gameId, questions)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to start game: %w", err)
	}
	return gameId, nil
}

// EndGame deletes questions from gameState, returns game`s answers
func (s *Service) EndGame(ctx context.Context, gameId uuid.UUID) ([]models.QuestionWithAnswers, error) {
	err := s.gamesRepo.End(ctx, gameId)
	if err != nil {
		return nil, fmt.Errorf("failed to end game: %w", err)
	}
	questionsWithAnswer, err := s.questionsRepo.GetQuestionsWithAnswers(ctx, gameId)
	if err != nil {
		return nil, fmt.Errorf("failed to end game: %w", err)
	}
	err = s.gameStorage.DeleteQuestions(gameId)
	if err != nil {
		return nil, fmt.Errorf("failed to end game: %w", err)
	}
	return questionsWithAnswer, nil
}

func (s *Service) GetAndDeleteStorageQuestion(ctx context.Context, gameId uuid.UUID, questionNum int) (*storage.QuestionInStorage, error) {
	questionFromStorage, err := s.gameStorage.PopQuestion(gameId, questionNum)
	if err != nil {
		//return nil, fmt.Errorf("failed to pop question from gameStateStorage: %w", err)
		offset := questionNumToPregenerate + questionNum
		var questionsNew []storage.QuestionInStorage
		questionsNew, err = s.generateQuestions(questionNum, offset, gameId)
		if err != nil {
			return nil, fmt.Errorf("failed to get question from gameStateStorage: %w", err)
		}
		err = s.gameStorage.SetQuestions(gameId, questionsNew)
		if err != nil {
			return nil, fmt.Errorf("failed to save questions: %w", err)
		}
	}
	return questionFromStorage, nil
}

func (s *Service) GetQuestion(ctx context.Context, gameId uuid.UUID, questionNum int) (*schema.QuestionResp, error) {
	question, err := s.gameStorage.GetQuestion(gameId, questionNum)
	if err != nil && errors.Is(err, storage.GameIdError) {
		return nil, fmt.Errorf("game is already ended, please start new game: %w", err)
	}
	if question == nil {
		offset := questionNumToPregenerate + questionNum
		questionsNew, err := s.generateQuestions(questionNum, offset, gameId)
		if err != nil {
			return nil, fmt.Errorf("failed to get question from gameStateStorage: %w", err)
		}
		err = s.gameStorage.SetQuestions(gameId, questionsNew)
		if err != nil {
			return nil, fmt.Errorf("failed to save questions: %w", err)
		}
		question, _ = s.gameStorage.GetQuestion(gameId, questionNum)
		return &schema.QuestionResp{QuestionText: question.QuestionText, FlagUrl: question.FlagUrl}, nil
	}
	return &schema.QuestionResp{QuestionText: question.QuestionText, FlagUrl: question.FlagUrl}, nil
}

func (s *Service) AnswerTheQuestion(ctx context.Context, gameId uuid.UUID, questionNum int, answerText string) (*schema.AnswerResp, error) {
	// save the question
	questionStorage, err := s.gameStorage.PopQuestion(gameId, questionNum)
	if err != nil && errors.Is(err, storage.GameIdError) {
		return nil, fmt.Errorf("game is already ended, please start new game: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get question from gameStateStorage: %w", err)
	}
	questionDb, err := s.saveQuestionToDB(ctx, questionStorage)
	if err != nil {
		return nil, fmt.Errorf("failed to answer the question: %w", err)
	}
	isCorrect, err := s.isAnswerCorrect(ctx, questionStorage, answerText)
	if err != nil {
		return nil, fmt.Errorf("failed to answer the question: %w", err)
	}
	err = s.SaveAnswer(ctx, questionDb.QuestionId, answerText, isCorrect)
	if err != nil {
		return nil, fmt.Errorf("failed to answer the question: %w", err)
	}
	return &schema.AnswerResp{IsCorrect: isCorrect}, nil
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

func (s *Service) isAnswerCorrect(ctx context.Context, questionFromStorage *storage.QuestionInStorage, answer string) (bool, error) {
	lowered := strings.ToLower(answer)
	if lowered != strings.ToLower(questionFromStorage.Answer) {
		return false, nil
	}
	return true, nil
}

func (s *Service) SaveAnswer(ctx context.Context, questionId uuid.UUID, answerText string, isCorrect bool) error {
	answer := models.Answer{
		AnswerId:   uuid.New(),
		QuestionId: questionId,
		AnswerText: answerText,
		IsCorrect:  isCorrect,
		AnsweredAt: time.Now(),
	}
	err := s.answersRepo.Create(ctx, &answer)
	if err != nil {
		return fmt.Errorf("failed to save answer: %w", err)
	}
	return nil
}

func (s *Service) generateQuestions(questionNumStart int, questionNumEnd int, gameId uuid.UUID) ([]storage.QuestionInStorage, error) {
	questions := make([]storage.QuestionInStorage, questionNumEnd)
	for i := questionNumStart; i < questionNumEnd; i++ {
		country, err := s.extractNewCountry(gameId)
		if err != nil {
			return nil, fmt.Errorf("failed to generate questions: %w", err)
		}
		question := storage.QuestionInStorage{
			QuestionId:   uuid.New(),
			GameId:       gameId,
			CountryId:    country.CountryId,
			QuestionText: "Guess The Flag",
			FlagUrl:      country.FlagUrl,
			CreatedAt:    time.Now(),
			Answer:       country.Name,
		}
		questions[i] = question
	}
	return questions, nil
}

func (s *Service) extractNewCountry(gameId uuid.UUID) (*models.Country, error) {
	for i := 0; i < 5; i++ {
		country, err := s.countryStore.GetRandom()
		if err != nil {
			return nil, fmt.Errorf("failed to extract new country: %w", err)
		}
		if ok := s.gameStorage.IsCountryUsed(gameId, country.CountryId); !ok {
			err = s.gameStorage.SetCountry(gameId, country.CountryId)
			if err != nil {
				return nil, fmt.Errorf("failed to extract new country: %w", err)
			}
			return country, nil
		}
	}
	return nil, fmt.Errorf("failed to extract new country from game with id %s", gameId)
}
