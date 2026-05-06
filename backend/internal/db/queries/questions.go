package queries

var QuestionQueries = struct {
	Create                  string
	GetGameQuestions        string
	GetQuestionsWithAnswers string
	GetQuestion             string
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
		SELECT a.answer, a.is_correct, cn.name, c.code, i.svg_data AS flag_svg
		FROM questions q
		JOIN public.countries c ON c.country_id = q.country_id
		JOIN public.images i ON i.image_id = c.flag_image_id
		JOIN public.answers a ON a.question_id = q.question_id
		JOIN public.country_names cn ON cn.country_id = c.country_id
			AND cn.language_code = $2
			AND cn.is_display_name = true
		WHERE q.game_id = $1
	`,
	GetQuestion: `
	SELECT * from questions
	WHERE question_id = $1
	`,
}
