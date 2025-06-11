package queries

var CountryQueries = struct {
	GetAll  string
	GetByID string
	Create  string
}{
	GetAll: `
		SELECT country_id, name, code, flag_url
		FROM countries
		ORDER BY name
	`,
	GetByID: `
		SELECT country_id, name, code, flag_url
		FROM countries
		WHERE country_id = $1
	`,
	Create: `
		INSERT INTO countries (name, code, flag_url)
		VALUES (:name, :code, :flag_url)
	`,
}
