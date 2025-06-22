package user

import (
	"errors"
	"github.com/google/uuid"
	"github.com/ziyadrw/faslah/internal/base/utils"
	userModels "github.com/ziyadrw/faslah/internal/modules/user/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) *UserRepository {
	return &UserRepository{DB: DB}
}

// CreateUser creates a new user in the database
func (ur *UserRepository) CreateUser(user *userModels.User, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword)
	user.Email = utils.FormatEmail(user.Email)

	return ur.DB.Create(user).Error
}

// GetUserByEmail retrieves a user by email
func (ur *UserRepository) GetUserByEmail(email string) (*userModels.User, error) {
	var user userModels.User
	email = utils.FormatEmail(email)

	result := ur.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID
func (ur *UserRepository) GetUserByID(id uuid.UUID) (*userModels.User, error) {
	var user userModels.User

	result := ur.DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &user, nil
}

// VerifyPassword checks if the provided password matches the stored hash
func (ur *UserRepository) VerifyPassword(user *userModels.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}
