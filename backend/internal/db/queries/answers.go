package queries

var AnswerQueries = struct {
	Create            string
	GetQuestionAnswer string
}{
	Create: `
		INSERT INTO answers (answer_id, question_id, answer, answered_at, is_correct)
		VALUES (:answer_id, :question_id, :answer, :answered_at, :is_correct)
	`,
	GetQuestionAnswer: `
		SELECT answer_id, question_id, answer, answered_at, is_correct
		FROM answers
		WHERE question_id = $1
	`,
}
