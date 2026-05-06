package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pythonistD/Guess-The-Flag/internal/db/models"
	"github.com/pythonistD/Guess-The-Flag/internal/db/queries"
)

type ImagesRepo struct {
	db *sqlx.DB
}

func NewImagesRepo(db *sqlx.DB) *ImagesRepo {
	return &ImagesRepo{db: db}
}

// Create создаёт новое изображение и возвращает его ID
func (i *ImagesRepo) Create(ctx context.Context, image *models.FlagImage) (*models.FlagImage, error) {
	// Устанавливаем время создания, если не установлено
	if image.CreatedAt.IsZero() {
		image.CreatedAt = time.Now()
	}
	if image.UpdatedAt.IsZero() {
		image.UpdatedAt = time.Now()
	}

	rows, err := i.db.NamedQueryContext(ctx, queries.ImageQueries.Create, image)
	if err != nil {
		return nil, fmt.Errorf("failed to create image: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&image.ImageId)
		if err != nil {
			return nil, fmt.Errorf("failed to scan image id: %w", err)
		}
	}

	return image, nil
}

// CreateOrGet создаёт изображение или возвращает существующее по хэшу (дедупликация)
func (i *ImagesRepo) CreateOrGet(ctx context.Context, image *models.FlagImage) (int, error) {
	// Устанавливаем время создания, если не установлено
	if image.CreatedAt.IsZero() {
		image.CreatedAt = time.Now()
	}
	if image.UpdatedAt.IsZero() {
		image.UpdatedAt = time.Now()
	}

	var imageId int
	rows, err := i.db.NamedQueryContext(ctx, queries.ImageQueries.CreateOrGet, image)
	if err != nil {
		return 0, fmt.Errorf("failed to create or get image: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&imageId)
		if err != nil {
			return 0, fmt.Errorf("failed to scan image id: %w", err)
		}
	}

	return imageId, nil
}

// CreateAll вставляет массив изображений с игнорированием дубликатов
func (i *ImagesRepo) CreateAll(ctx context.Context, images []models.FlagImage) error {
	if len(images) == 0 {
		return nil
	}

	// Устанавливаем время для всех изображений
	now := time.Now()
	for idx := range images {
		if images[idx].CreatedAt.IsZero() {
			images[idx].CreatedAt = now
		}
		if images[idx].UpdatedAt.IsZero() {
			images[idx].UpdatedAt = now
		}
	}

	// Используем транзакцию для батчевой вставки
	tx, err := i.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, image := range images {
		_, err := tx.NamedExecContext(ctx, queries.ImageQueries.CreateAll, image)
		if err != nil {
			return fmt.Errorf("failed to batch insert image: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetByHash получает изображение по хэшу
func (i *ImagesRepo) GetByHash(ctx context.Context, hash string) (*models.FlagImage, error) {
	var image models.FlagImage
	err := i.db.GetContext(ctx, &image, queries.ImageQueries.GetByHash, hash)
	if err != nil {
		return nil, fmt.Errorf("error getting image by hash: %w", err)
	}
	return &image, nil
}

// GetById получает изображение по ID
func (i *ImagesRepo) GetById(ctx context.Context, id int) (*models.FlagImage, error) {
	var image models.FlagImage
	err := i.db.GetContext(ctx, &image, queries.ImageQueries.GetById, id)
	if err != nil {
		return nil, fmt.Errorf("error getting image by id: %w", err)
	}
	return &image, nil
}
