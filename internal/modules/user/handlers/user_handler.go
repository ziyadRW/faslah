package user

import (
	"github.com/labstack/echo/v4"
	"github.com/ziyadrw/faslah/internal/base"
	userDTOs "github.com/ziyadrw/faslah/internal/modules/user/dtos"
	userServices "github.com/ziyadrw/faslah/internal/modules/user/services"
)

type UserHandler struct {
	UserService *userServices.UserService
}

func NewUserHandler(userService *userServices.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

// Signup godoc
// @Summary تسجيل مستخدم جديد
// @Description تسجيل مستخدم جديد وإنشاء حساب
// @Tags المستخدمين والمصادقة
// @Accept json
// @Produce json
// @Param request body userDTOs.SignupRequest true "بيانات التسجيل"
// @Success 200 {object} base.Response{data=userDTOs.AuthResponse} "تم التسجيل بنجاح"
// @Failure 400 {object} base.Response "خطأ في البيانات المدخلة"
// @Failure 500 {object} base.Response "خطأ في الخادم"
// @Router /users/signup [post]
func (uh *UserHandler) Signup(c echo.Context) error {
	var dto userDTOs.SignupRequest
	if res, ok := base.BindAndValidate(c, &dto); !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	response := uh.UserService.Signup(dto)
	return c.JSON(response.HTTPStatus, response)
}

// Login godoc
// @Summary تسجيل الدخول
// @Description مصادقة المستخدم وإنشاء رمز JWT
// @Tags المستخدمين والمصادقة
// @Accept json
// @Produce json
// @Param request body userDTOs.LoginRequest true "بيانات تسجيل الدخول"
// @Success 200 {object} base.Response{data=userDTOs.AuthResponse} "تم تسجيل الدخول بنجاح"
// @Failure 400 {object} base.Response "خطأ في البيانات المدخلة"
// @Failure 401 {object} base.Response "بيانات الاعتماد غير صالحة"
// @Failure 500 {object} base.Response "خطأ في الخادم"
// @Router /users/login [post]
func (uh *UserHandler) Login(c echo.Context) error {
	var dto userDTOs.LoginRequest
	if res, ok := base.BindAndValidate(c, &dto); !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	response := uh.UserService.Login(dto)
	return c.JSON(response.HTTPStatus, response)
}

// GetProfile godoc
// @Summary الحصول على الملف الشخصي
// @Description استرجاع معلومات الملف الشخصي للمستخدم الحالي
// @Tags المستخدمين والمصادقة
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} base.Response{data=userDTOs.UserResponse} "تم استرجاع الملف الشخصي بنجاح"
// @Failure 401 {object} base.Response "غير مصرح"
// @Failure 500 {object} base.Response "خطأ في الخادم"
// @Router /users/me [get]
func (uh *UserHandler) GetProfile(c echo.Context) error {
	userID := c.Get("user_id").(string)

	response := uh.UserService.GetProfile(userID)
	return c.JSON(response.HTTPStatus, response)
}
