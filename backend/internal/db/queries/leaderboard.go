package queries

var LeaderboardQueries = struct {
	Create  string
	GetByID string
	Update  string
	Delete  string
	Upsert  string
}{
	Create: `
		INSERT INTO leaderboard (leaderboard_id, user_id, score, games_played, last_game)
		VALUES (:leaderboard_id, :user_id, :score, :games_played, :last_game)
	`,
	GetByID: `
		SELECT leaderboard_id, user_id, score, games_played, last_game
		FROM leaderboard
		WHERE leaderboard_id = $1
	`,
	Update: `
		UPDATE leaderboard
		SET score = :score, games_played = :games_played, last_game = :last_game
		WHERE leaderboard_id = :leaderboard_id
	`,
	Delete: `
		DELETE FROM leaderboard
		WHERE leaderboard_id = $1
	`,
	Upsert: `
		INSERT INTO leaderboard (user_id, score, games_played, last_game)
		VALUES (:user_id, :score, :games_played, :last_game)
		ON CONFLICT (user_id) DO UPDATE
		SET score = leaderboard.score + EXCLUDED.score,
			games_played = leaderboard.games_played + 1,
			last_game = EXCLUDED.last_game
	`,
}
