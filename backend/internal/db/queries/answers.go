package queries

var AnswerQueries = struct {
	Create               string
	GetQuestionAnswer    string
	GetQuestionsByGameId string
}{
	Create: `
		INSERT INTO answers (answer_id, question_id, answer, answered_at, is_correct, selected_country_id)
		VALUES (:answer_id, :question_id, :answer, :answered_at, :is_correct, :selected_country_id)
	`,
	GetQuestionAnswer: `
		SELECT answer_id, question_id, answer, answered_at, is_correct, selected_country_id
		FROM answers
		WHERE question_id = $1
	`,
	GetQuestionsByGameId: `
		SELECT a.answer_id,
			   a.question_id,
			   a.answer,
			   a.answered_at,
			   a.is_correct,
			   a.selected_country_id
		FROM answers a
		JOIN questions q ON a.question_id = q.question_id
		WHERE q.game_id = $1
	`,
}
