package queries

var CountryQueries = struct {
	GetAll           string
	GetByID          string
	GetByCode        string
	Create           string
	GetAllWithImage  string
	GetByIdWithImage string
}{
	GetAll: `
		SELECT country_id, code, flag_image_id
		FROM countries
	`,
	GetByID: `
		SELECT country_id, code, flag_image_id
		FROM countries
		WHERE country_id = $1
	`,
	GetByCode: `
		SELECT country_id, code, flag_image_id
		FROM countries
		WHERE code = $1
	`,
	Create: `
		INSERT INTO countries (code, flag_image_id)
		VALUES (:code, :flag_image_id)
		RETURNING country_id
	`,
	GetAllWithImage: `
		select country_id, code, flag_image_id, svg_data, image_hash
		from countries c
		join images on images.image_id = c.flag_image_id
	`,
	GetByIdWithImage: `
		select country_id, code, flag_image_id, svg_data, image_hash
		from countries c
		join images on images.image_id = c.flag_image_id
		where country_id = $1
	`,
}
