package queries

var QuestionQueries = struct {
	Create                  string
	GetGameQuestions        string
	GetQuestionsWithAnswers string
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
	GetQuestionsWithAnswers: `
		SELECT answer, is_correct, name, code, flag_url
		FROM questions
		JOIN public.countries c on c.country_id = questions.country_id
		JOIN public.answers a on questions.question_id = a.question_id
		WHERE game_id = $1
	`,
}
