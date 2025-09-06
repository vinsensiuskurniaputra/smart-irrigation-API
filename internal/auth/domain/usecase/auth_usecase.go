package usecase

import (
	"errors"

	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/auth/data/models"
	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/auth/data/repositories"
	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/core/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepo repositories.UserRepository
	jwtKey   []byte
}

func NewAuthUsecase(userRepo repositories.UserRepository) *AuthUsecase {
	return &AuthUsecase{userRepo: userRepo}
}

func (uc *AuthUsecase) Login(username, password string) (string, *models.User, error) {
	user, err := uc.userRepo.FindByUsername(username)
	if err != nil {
		return "", nil, err
	}

	// cek password
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", nil, errors.New("invalid password")
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (uc *AuthUsecase) Register(username, password, name, email string) error {
	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Username: username,
		Password: string(hash),
		Name:     name,
		Email:    email,
	}

	return uc.userRepo.Create(user)
}

func (uc *AuthUsecase) GetMe(userID uint) (*models.User, error) {
	return uc.userRepo.FindByID(userID)
}

func (uc *AuthUsecase) UpdateMyPassword(userID uint, oldPassword, newPassword string) error {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	// verify old password
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)) != nil {
		return errors.New("invalid old password")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hash)
	return uc.userRepo.Update(user)
}
