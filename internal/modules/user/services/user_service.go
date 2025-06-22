package user

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ziyadrw/faslah/internal/base"
	userDTOs "github.com/ziyadrw/faslah/internal/modules/user/dtos"
	userEnums "github.com/ziyadrw/faslah/internal/modules/user/enums"
	userModels "github.com/ziyadrw/faslah/internal/modules/user/models"
	userRepo "github.com/ziyadrw/faslah/internal/modules/user/repositories"
)

type UserService struct {
	UserRepository *userRepo.UserRepository
}

func NewUserService(userRepository *userRepo.UserRepository) *UserService {
	return &UserService{UserRepository: userRepository}
}

// Signup registers a new user
func (us *UserService) Signup(dto userDTOs.SignupRequest) base.Response {
	existingUser, err := us.UserRepository.GetUserByEmail(dto.Email)
	if err != nil {
		return base.SetErrorMessage("خطأ في الخادم", err.Error())
	}
	if existingUser != nil {
		return base.SetErrorMessage("البريد الإلكتروني مستخدم بالفعل", "يرجى استخدام بريد إلكتروني آخر أو تسجيل الدخول")
	}

	role := dto.Role
	if role == "" {
		role = userEnums.TypeViewer
	}

	user := &userModels.User{
		Email: dto.Email,
		Name:  dto.Name,
		Role:  role,
	}

	if err := us.UserRepository.CreateUser(user, dto.Password); err != nil {
		return base.SetErrorMessage("فشل في إنشاء المستخدم", err.Error())
	}

	token, err := us.generateJWT(user.ID.String())
	if err != nil {
		return base.SetErrorMessage("فشل في إنشاء رمز المصادقة", err.Error())
	}

	response := userDTOs.AuthResponse{
		Token: token,
		User: userDTOs.UserResponse{
			ID:        user.ID.String(),
			Email:     user.Email,
			Name:      user.Name,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
	}

	return base.SetData(response, "تم التسجيل بنجاح")
}

// Login authenticates a user
func (us *UserService) Login(dto userDTOs.LoginRequest) base.Response {
	user, err := us.UserRepository.GetUserByEmail(dto.Email)
	if err != nil {
		return base.SetErrorMessage("خطأ في الخادم", err.Error())
	}
	if user == nil {
		return base.SetErrorMessage("بيانات الاعتماد غير صالحة", "البريد الإلكتروني أو كلمة المرور غير صحيحة")
	}

	if !us.UserRepository.VerifyPassword(user, dto.Password) {
		return base.SetErrorMessage("بيانات الاعتماد غير صالحة", "البريد الإلكتروني أو كلمة المرور غير صحيحة")
	}

	token, err := us.generateJWT(user.ID.String())
	if err != nil {
		return base.SetErrorMessage("فشل في إنشاء رمز المصادقة", err.Error())
	}

	response := userDTOs.AuthResponse{
		Token: token,
		User: userDTOs.UserResponse{
			ID:        user.ID.String(),
			Email:     user.Email,
			Name:      user.Name,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
	}

	return base.SetData(response, "تم تسجيل الدخول بنجاح")
}

// GetProfile retrieves the user profile
func (us *UserService) GetProfile(userID string) base.Response {
	id, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("معرف المستخدم غير صالح", err.Error())
	}

	user, err := us.UserRepository.GetUserByID(id)
	if err != nil {
		return base.SetErrorMessage("خطأ في الخادم", err.Error())
	}
	if user == nil {
		return base.SetErrorMessage("المستخدم غير موجود", "لم يتم العثور على المستخدم")
	}

	response := userDTOs.UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}

	return base.SetData(response)
}

func (us *UserService) generateJWT(userID string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("JWT_SECRET environment variable not set")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
