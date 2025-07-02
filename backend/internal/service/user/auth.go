package user

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type TokenManager interface {
	GenerateToken(userId uuid.UUID) (string, error)
	ValidateToken(token string) (uuid.UUID, error)
}

type PasswordManager interface {
	Hash(password string) (string, error)
	Validate(passwordProvided string, passwordHash string) error
}

type JWTManager struct {
	secretKey     string
	tokenLifetime time.Duration
}

func NewJWTManager(secretKey string, tokenLifetime time.Duration) *JWTManager {
	return &JWTManager{secretKey, tokenLifetime}
}

func (j *JWTManager) GenerateToken(userId uuid.UUID) (string, error) {
	payload := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(j.tokenLifetime).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTManager) ValidateToken(token string) (uuid.UUID, error) {
	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
	if err != nil || !tokenParsed.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}
	payload, ok := tokenParsed.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid payload")
	}
	userIdStr, ok := payload["user_id"].(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("user_id missing in token")
	}
	return uuid.Parse(userIdStr)
}

func (j *JWTManager) GetSecretKey() string {
	return j.secretKey
}

type PasswordManagerImpl struct {
}

func NewPasswordManager() *PasswordManagerImpl {
	return &PasswordManagerImpl{}
}

func (p *PasswordManagerImpl) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}
	return string(hash), nil
}

func (p *PasswordManagerImpl) Validate(passwordProvided string, passwordHash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(passwordProvided))
	if err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}
	return nil
}
