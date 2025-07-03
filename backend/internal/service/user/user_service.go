package user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pythonistD/Guess-The-Flag/internal/db/models"
	"github.com/pythonistD/Guess-The-Flag/internal/db/repo"
	"github.com/pythonistD/Guess-The-Flag/internal/schema"
)

type Service interface {
	Register(ctx context.Context, register schema.Register) (schema.Token, error)
	Login(ctx context.Context, login schema.Login) (schema.Token, error)
}

type ServiceImpl struct {
	usersRepo   *repo.UsersRepo
	jwtManager  TokenManager
	passManager PasswordManager
}

func NewService(db *sqlx.DB, jwtManager TokenManager, passManager PasswordManager) *ServiceImpl {
	usersRepo := repo.NewUsersRepo(db)
	return &ServiceImpl{
		usersRepo:   usersRepo,
		jwtManager:  jwtManager,
		passManager: passManager,
	}
}

func (s *ServiceImpl) Login(ctx context.Context, login schema.Login) (schema.Token, error) {
	user, err := s.usersRepo.GetByUsername(ctx, login.Username)
	if err != nil {
		return schema.Token{}, fmt.Errorf("authentication error: %w", err)
	}
	err = s.passManager.Validate(login.Password, user.Password)
	if err != nil {
		return schema.Token{}, fmt.Errorf("authentication error: %w", err)
	}
	token, err := s.jwtManager.GenerateToken(user.Uuid)
	if err != nil {
		return schema.Token{}, fmt.Errorf("authentication error: %w", err)
	}
	return schema.Token{Token: token}, nil
}

func (s *ServiceImpl) Register(ctx context.Context, register schema.Register) (schema.Token, error) {
	userId := uuid.New()
	token, err := s.jwtManager.GenerateToken(userId)
	if err != nil {
		return schema.Token{}, fmt.Errorf("could not generate token: %w", err)
	}
	passHash, err := s.passManager.Hash(register.Password)
	if err != nil {
		return schema.Token{}, fmt.Errorf("could not hash password: %w", err)
	}
	user := models.User{
		Uuid:      userId,
		Name:      register.Username,
		Email:     register.Email,
		Password:  passHash,
		CreatedAt: time.Time{},
	}
	err = s.usersRepo.Create(ctx, &user)
	if err != nil {
		return schema.Token{}, fmt.Errorf("registration error: %w", err)
	}
	return schema.Token{Token: token}, nil
}
