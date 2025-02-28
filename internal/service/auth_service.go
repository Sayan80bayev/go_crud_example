package service

import (
	"errors"
	"go_crud_example/internal/models"
	"go_crud_example/internal/repository"

	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
	jwtKey   []byte
}

func NewAuthService(userRepo *repository.UserRepository, jwtKey string) *AuthService {
	return &AuthService{userRepo, []byte(jwtKey)}
}

func (uc *AuthService) Register(username, password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &models.User{
		Username: username,
		Password: string(hashedPassword),
	}
	err = uc.userRepo.CreateUser(user)
	if err != nil {
		return "", err
	}

	token := generateJwtToken(user)

	return token.SignedString(uc.jwtKey)
}

func (uc *AuthService) Login(username, password string) (string, error) {
	user, err := uc.userRepo.GetUserByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token := generateJwtToken(user)

	return token.SignedString(uc.jwtKey)
}

func generateJwtToken(user *models.User) *jwt.Token {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"username":  user.Username,
		"user_role": user.Role,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	})
	return token
}
