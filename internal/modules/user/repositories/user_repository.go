package user

import (
	"errors"
	"time"

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

// GetWatchHistory retrieves the watch history for a user
func (ur *UserRepository) GetWatchHistory(userID uuid.UUID) ([]userModels.WatchHistory, error) {
	var history []userModels.WatchHistory

	result := ur.DB.Where("user_id = ?", userID).Order("last_played_at DESC").Find(&history)
	if result.Error != nil {
		return nil, result.Error
	}

	return history, nil
}

// UpsertWatchHistory creates or updates a watch history record
func (ur *UserRepository) UpsertWatchHistory(userID, podcastID uuid.UUID, playbackSecond int) error {
	now := time.Now()

	var count int64
	if err := ur.DB.Model(&userModels.WatchHistory{}).
		Where("user_id = ? AND podcast_id = ?", userID, podcastID).
		Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return ur.DB.Model(&userModels.WatchHistory{}).
			Where("user_id = ? AND podcast_id = ?", userID, podcastID).
			Updates(map[string]interface{}{
				"playback_second": playbackSecond,
				"last_played_at":  now,
				"updated_at":      now,
			}).Error
	} else {
		history := userModels.WatchHistory{
			UserID:         userID.String(),
			PodcastID:      podcastID.String(),
			PlaybackSecond: playbackSecond,
			LastPlayedAt:   &now,
		}
		return ur.DB.Create(&history).Error
	}
}
