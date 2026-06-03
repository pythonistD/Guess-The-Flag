package queries

var CommonQueries = struct {
	ClearAllTables string
}{
	ClearAllTables: `
	truncate table answers, countries, country_names, games, images, leaderboard,
	questions, unknown_flags, users restart identity cascade;
	`,
}
