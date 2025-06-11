package queries

var QuestionQueries = struct {
	Create           string
	GetGameQuestions string
}{
	Create: `
		INSERT INTO questions (question_id, game_id, country_id, created_at)
		VALUES (:question_id, :game_id, :country_id, :created_at)
	`,
	GetGameQuestions: `
		SELECT question_id, game_id, country_id, created_at
		FROM questions
		WHERE game_id = $1
	`,
}
