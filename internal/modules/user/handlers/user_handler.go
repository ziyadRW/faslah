package user

import (
	"net/http"

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

// GetWatchHistory godoc
// @Summary الحصول على سجل المشاهدة
// @Description استرجاع سجل مشاهدة البودكاست للمستخدم الحالي
// @Tags المستخدمين والمصادقة
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} base.Response{data=[]userDTOs.WatchHistoryResponse} "تم استرجاع سجل المشاهدة بنجاح"
// @Failure 401 {object} base.Response "غير مصرح"
// @Failure 500 {object} base.Response "خطأ في الخادم"
// @Router /users/me/history [get]
func (uh *UserHandler) GetWatchHistory(c echo.Context) error {
	userID := c.Get("user_id").(string)

	response := uh.UserService.GetWatchHistory(userID)
	return c.JSON(response.HTTPStatus, response)
}

// TrackPlay godoc
// @Summary تتبع تشغيل البودكاست
// @Description تسجيل موضع التشغيل الحالي للبودكاست
// @Tags تشغيل البودكاست
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "معرف البودكاست"
// @Param request body userDTOs.TrackPlayRequest true "بيانات التشغيل"
// @Success 204 "تم تتبع التشغيل بنجاح"
// @Failure 400 {object} base.Response "خطأ في البيانات المدخلة"
// @Failure 401 {object} base.Response "غير مصرح"
// @Failure 500 {object} base.Response "خطأ في الخادم"
// @Router /podcasts/{id}/track-play [post]
func (uh *UserHandler) TrackPlay(c echo.Context) error {
	userID := c.Get("user_id").(string)
	podcastID := c.Param("id")

	var dto userDTOs.TrackPlayRequest
	if res, ok := base.BindAndValidate(c, &dto); !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	response := uh.UserService.TrackPlay(userID, podcastID, dto)
	if response.HTTPStatus == http.StatusOK {
		return c.NoContent(http.StatusNoContent)
	}
	return c.JSON(response.HTTPStatus, response)
}
