package game

import "github.com/pythonistD/Guess-The-Flag/internal/db/repo"

type GameService struct {
	gamesRepo     *repo.GamesRepo
	questionsRepo *repo.QuestionsRepo
	answersRepo   *repo.AnswersRepo
	countriesRepo *repo.CountriesRepo
}
