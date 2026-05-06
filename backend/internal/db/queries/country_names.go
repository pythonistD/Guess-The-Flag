package queries

var CountryNamesQueries = struct {
	Create                    string
	GetCommonNamesByCountryId string
	GetAllNamesByCountryId    string
}{
	Create: `
		INSERT INTO country_names (language_code, country_id, name, normalized_name, threshold, is_display_name)
		VALUES (:language_code, :country_id, :name, :normalized_name, :threshold, :is_display_name)
	`,
	GetCommonNamesByCountryId: `
		select * from country_names where country_id = $1 and is_display_name = true
	`,
	GetAllNamesByCountryId: `
		select * from country_names where country_id = $1
	`,
}
