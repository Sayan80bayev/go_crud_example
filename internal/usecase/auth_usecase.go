package usecase

import (
	"errors"
	"go_crud_example/internal/models"
	"go_crud_example/internal/repository"

	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepo *repository.UserRepository
	jwtKey   []byte
}

func NewAuthUsecase(userRepo *repository.UserRepository, jwtKey string) *AuthUsecase {
	return &AuthUsecase{userRepo, []byte(jwtKey)}
}

func (uc *AuthUsecase) Register(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Username: username,
		Password: string(hashedPassword),
	}

	return uc.userRepo.CreateUser(user)
}

func (uc *AuthUsecase) Login(username, password string) (string, error) {
	user, err := uc.userRepo.GetUserByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(uc.jwtKey)
}
