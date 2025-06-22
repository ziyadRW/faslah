package user

import (
	"time"

	userEnums "github.com/ziyadrw/faslah/internal/modules/user/enums"
)

// SignupRequest represents the request body for user registration
type SignupRequest struct {
	Email    string         `json:"email" validate:"required,email" message:"البريد الإلكتروني مطلوب وغير صالح"`
	Password string         `json:"password" validate:"required,passwordvalidator" message:"كلمة المرور مطلوبة ويجب أن تحتوي على 6 أحرف على الأقل وتتضمن حرفًا واحدًا ورقمًا واحدًا على الأقل"`
	Name     string         `json:"name" validate:"required" message:"الاسم مطلوب"`
	Role     userEnums.Type `json:"role" validate:"omitempty,oneof=viewer creator" message:"الدور يجب أن يكون إما مشاهد أو صانع محتوى"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" message:"البريد الإلكتروني مطلوب وغير صالح"`
	Password string `json:"password" validate:"required" message:"كلمة المرور مطلوبة"`
}

// UserResponse represents the response body for user profile
type UserResponse struct {
	ID        string         `json:"id"`
	Email     string         `json:"email"`
	Name      string         `json:"name"`
	Role      userEnums.Type `json:"role"`
	CreatedAt time.Time      `json:"created_at"`
}

// AuthResponse represents the response body for authentication
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// WatchHistoryResponse represents the response body for watch history
type WatchHistoryResponse struct {
	PodcastID      string     `json:"podcast_id"`
	PlaybackSecond int        `json:"playback_second"`
	LastPlayedAt   *time.Time `json:"last_played_at"`
}

// TrackPlayRequest represents the request body for tracking podcast play
type TrackPlayRequest struct {
	PlaybackSecond int `json:"playback_second" validate:"required,min=0" message:"ثانية التشغيل مطلوبة ويجب أن تكون قيمة موجبة"`
}
