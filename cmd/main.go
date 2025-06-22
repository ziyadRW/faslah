package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/ziyadrw/faslah/config"
	"github.com/ziyadrw/faslah/docs"
	"github.com/ziyadrw/faslah/internal/base"
	"github.com/ziyadrw/faslah/internal/middlewares"
	"github.com/ziyadrw/faslah/internal/migrations"
	"github.com/ziyadrw/faslah/internal/routes"
)

func init() {
	config.LoadEnv()
	baseURL := config.GetEnv("BASE_URL", "localhost:8080")
	appDomain := config.GetEnv("APP_DOMAIN", "faslah.net")
	docs.SwaggerInfo.Description = fmt.Sprintf(
		`<div dir="rtl" style="text-align: right;">
	إنشاء بودكاست جديد من خلال رفع ملف MP4 مباشرة أو من خلال رابط <strong>يوتيوب</strong>.
	يمكنك إما تقديم رابط فيديو <strong>يوتيوب</strong> وسنقوم بتنزيله واستخراج البيانات الوصفية تلقائيًا،
	أو يمكنك رفع ملف MP4 مباشرة وتقديم البيانات الوصفية يدويًا.
	في كلتا الحالتين، سيتم تخزين الفيديو في خدمة Cloudflare R2 الخاصة بنا تحت نطاق media.%s للوصول السريع والآمن.
	</div>`,
		appDomain,
	)
	docs.SwaggerInfo.Title = "🔊 فاصلة | واجهات برمجة خلفية"
	docs.SwaggerInfo.Host = baseURL
}

// @title فاصلة API
// @version 1.0
// @description واجهة برمجة التطبيقات لمنصة فاصلة للبودكاست
// @termsOfService http://swagger.io/terms/
// @contact.name فريق دعم فاصلة
// @contact.url https://github.com/ziyadrw/faslah
// @contact.email zeadAlrouasheed@gmail.com

// @BasePath /
// @schemes http https

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description أدخل رمز JWT مع البادئة Bearer، مثال: "Bearer abcdef123456"
func main() {
	e := echo.New()

	base.RegisterValidator(e)
	middlewares.RegisterAllGlobalMiddlewares(e)

	config.Connect()
	db := config.GetDB()
	migrations.Migrate()

	routes.RegisterAllRoutes(e, db)

	startServer(e)
}
