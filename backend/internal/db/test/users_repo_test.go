package test

import (
	"context"
	"github.com/google/uuid"
	"github.com/pythonistD/Guess-The-Flag/internal/db/models"
	"github.com/pythonistD/Guess-The-Flag/internal/db/repo"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	db := newTestDB(t)
	defer clearTables(t, db, []string{"users"})
	userRepo := repo.NewUsersRepo(db)
	user := models.User{
		Uuid:      uuid.New(),
		Name:      "John Doe",
		Email:     "test@mail.com",
		CreatedAt: time.Now().UTC(),
	}
	ctx := context.Background()
	err := userRepo.Create(ctx, &user)
	if err != nil {
		t.Fatal(err)
	}
	userFetched, err := userRepo.GetById(ctx, user.Uuid)
	assert.Equal(t, user.Uuid, userFetched.Uuid)
	assert.Equal(t, user.Name, userFetched.Name)
	assert.Equal(t, user.Email, userFetched.Email)
	t.Logf("userFetched CreatedAt: %+v", userFetched.CreatedAt)
	t.Logf("user CreatedAt: %+v", user.CreatedAt)
	assert.WithinDuration(t, user.CreatedAt, userFetched.CreatedAt, time.Second)
}
