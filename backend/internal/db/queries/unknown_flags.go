package queries

var UnknownFlagQueries = struct {
	Create           string
	GetById          string
	GetUserCountries string
	Delete           string
}{
	Create: `
		INSERT INTO unknown_flags (user_id, country_id)
		VALUES (:user_id, :country_id)
		ON CONFLICT DO NOTHING
	`,
	GetById: `
		SELECT unknown_flags_id, user_id, country_id
		FROM unknown_flags
		WHERE unknown_flags_id = $1
	`,
	Delete: `
		DELETE FROM unknown_flags
		WHERE unknown_flags_id = $1
	`,
	GetUserCountries: `
		SELECT country_id FROM unknown_flags
		WHERE user_id = $1
	`,
}
