package middlewares

import (
	"github.com/ziyadrw/faslah/config"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/ziyadrw/faslah/internal/base"
	"github.com/ziyadrw/faslah/internal/base/utils"
	userEnums "github.com/ziyadrw/faslah/internal/modules/user/enums"
	"gorm.io/gorm"
)

func RoleMiddleware(db *gorm.DB, allowedRoles ...userEnums.Type) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.GetEnv("APP_ENV", "") != "production" {
				var testingUserID string
				if err := db.Table("users").Select("id").Order("RANDOM()").Limit(1).Scan(&testingUserID).Error; err != nil {
					return next(c)
				}
				c.Set("user_id", testingUserID)
				return next(c)
			}

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, base.SetErrorMessage("غير مصرح", "يرجى تسجيل الدخول للوصول إلى هذا المورد"))
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				return c.JSON(http.StatusUnauthorized, base.SetErrorMessage("تنسيق غير صالح", "يجب أن يكون التنسيق: Bearer <token>"))
			}

			userID, err := utils.ExtractUserIDFromToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, base.SetErrorMessage("رمز غير صالح", err.Error()))
			}

			var count int64
			if err := db.Table("users").Where("id = ? AND deleted_at IS NULL", userID).Count(&count).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, base.SetErrorMessage("خطأ في الخادم", "فشل التحقق من المستخدم"))
			}

			if count == 0 {
				return c.JSON(http.StatusUnauthorized, base.SetErrorMessage("غير مصرح", "المستخدم غير موجود أو تم حذفه"))
			}

			c.Set("user_id", userID.String())

			if len(allowedRoles) == 0 {
				return next(c)
			}

			var userRole string
			if err := db.Table("users").Select("role").Where("id = ?", userID).Scan(&userRole).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, base.SetErrorMessage("خطأ في الخادم", "فشل التحقق من دور المستخدم"))
			}

			hasPermission := false
			for _, role := range allowedRoles {
				if string(role) == userRole {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				return c.JSON(http.StatusForbidden, base.SetErrorMessage("غير مصرح", "ليس لديك الصلاحيات الكافية للوصول إلى هذا المورد"))
			}

			c.Set("user_role", userRole)
			return next(c)
		}
	}
}
