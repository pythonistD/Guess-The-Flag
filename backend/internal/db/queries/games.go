package queries

var GameQueries = struct {
	Start           string
	End             string
	GetByID         string
	GetLastUserGame string
}{
	Start: `
		INSERT INTO games (game_id, user_id, started_at)
		VALUES (:game_id, :user_id, :started_at)
	`,
	End: `
		UPDATE games
		SET ended_at = :ended_at
		WHERE game_id = :game_id
	`,
	GetByID: `
		SELECT game_id, user_id, started_at, ended_at
		FROM games
		WHERE game_id = $1
	`,
	GetLastUserGame: `
		SELECT game_id, user_id, started_at, ended_at
				FROM games
				WHERE user_id = $1
				ORDER BY started_at DESC
				LIMIT 1
	`,
}
