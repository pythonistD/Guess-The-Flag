package queries

var ImageQueries = struct {
	Create      string
	CreateOrGet string
	GetByHash   string
	GetById     string
	CreateAll   string
}{
	Create: `
		INSERT INTO images (svg_data, image_hash, file_size, created_at, updated_at)
		VALUES (:svg_data, :image_hash, :file_size, :created_at, :updated_at)
		RETURNING image_id
	`,
	CreateOrGet: `
		INSERT INTO images (svg_data, image_hash, file_size, created_at, updated_at)
		VALUES (:svg_data, :image_hash, :file_size, :created_at, :updated_at)
		ON CONFLICT (image_hash) DO UPDATE SET updated_at = EXCLUDED.updated_at
		RETURNING image_id
	`,
	GetByHash: `
		SELECT image_id, svg_data, image_hash, file_size, created_at, updated_at
		FROM images
		WHERE image_hash = $1
	`,
	GetById: `
		SELECT image_id, svg_data, image_hash, file_size, created_at, updated_at
		FROM images
		WHERE image_id = $1
	`,
	CreateAll: `
		INSERT INTO images (svg_data, image_hash, file_size, created_at, updated_at)
		VALUES (:svg_data, :image_hash, :file_size, :created_at, :updated_at)
		ON CONFLICT (image_hash) DO NOTHING
	`,
}
