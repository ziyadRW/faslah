package cms

import (
	"github.com/labstack/echo/v4"
	"github.com/ziyadrw/faslah/internal/base"
	podcastDTOs "github.com/ziyadrw/faslah/internal/modules/cms/dtos"
	podcastServices "github.com/ziyadrw/faslah/internal/modules/cms/services"
)

type PodcastHandler struct {
	PodcastService *podcastServices.PodcastService
}

func NewPodcastHandler(podcastService *podcastServices.PodcastService) *PodcastHandler {
	return &PodcastHandler{
		PodcastService: podcastService,
	}
}

// CreateContent godoc
// @Summary إنشاء بودكاست جديد
// @Description إنشاء بودكاست جديد من خلال رفع ملف MP4 مباشرة أو من خلال رابط يوتيوب. يمكنك إما تقديم رابط فيديو يوتيوب وسنقوم بتنزيله واستخراج البيانات الوصفية تلقائيًا، أو يمكنك رفع ملف MP4 مباشرة وتقديم البيانات الوصفية يدويًا. في كلتا الحالتين، سيتم تخزين الفيديو في خدمة Cloudflare R2 الخاصة بنا تحت نطاق media.faslah.net للوصول السريع والآمن.
// @Tags إنشاء المحتوى
// @Accept multipart/form-data
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param source_url formData string false "رابط يوتيوب (مطلوب إذا لم يتم تقديم ملف). سنقوم بتنزيل الفيديو واستخراج العنوان والوصف والوسوم تلقائيًا."
// @Param file formData file false "ملف MP4 (مطلوب إذا لم يتم تقديم رابط يوتيوب). يجب أن يكون بتنسيق MP4."
// @Param title formData string false "العنوان (مطلوب لرفع الملف، اختياري لرابط يوتيوب). إذا قدمت رابط يوتيوب، يمكنك تجاوز العنوان المستخرج تلقائيًا."
// @Param description formData string false "الوصف (مطلوب لرفع الملف، اختياري لرابط يوتيوب). إذا قدمت رابط يوتيوب، يمكنك تجاوز الوصف المستخرج تلقائيًا."
// @Param tags formData []string false "الوسوم (مطلوبة لرفع الملف، اختيارية لرابط يوتيوب). قائمة من الكلمات المفتاحية لتسهيل البحث."
// @Param published_at formData string false "تاريخ النشر (اختياري). بتنسيق ISO 8601 (YYYY-MM-DDTHH:MM:SSZ)."
// @Success 200 {object} base.Response{data=podcastDTOs.PodcastResponse} "تم إنشاء البودكاست بنجاح"
// @Failure 400 {object} base.Response "خطأ في البيانات المدخلة"
// @Failure 401 {object} base.Response "غير مصرح"
// @Failure 500 {object} base.Response "خطأ في الخادم"
// @Router /cms/create-content [post]
func (ph *PodcastHandler) CreateContent(c echo.Context) error {
	userID := c.Get("user_id").(string)
	var dto podcastDTOs.CreatePodcastRequest

	file, err := c.FormFile("file")
	if err == nil {
		dto.File = file
	}

	if res, ok := base.BindAndValidate(c, &dto); !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	response := ph.PodcastService.CreateContent(dto, userID)
	return c.JSON(response.HTTPStatus, response)
}

// GetContent godoc
// @Summary الحصول على محتوى
// @Description استرجاع بودكاست بواسطة المعرف
// @Tags إدارة المحتوى
// @Accept json
// @Produce json
// @Param id path string true "معرف البودكاست"
// @Success 200 {object} base.Response{data=podcastDTOs.PodcastResponse} "تم استرجاع البودكاست بنجاح"
// @Failure 400 {object} base.Response "معرف البودكاست غير صالح"
// @Failure 404 {object} base.Response "البودكاست غير موجود"
// @Failure 500 {object} base.Response "خطأ في الخادم"
// @Router /cms/retreive-content/{id} [get]
func (ph *PodcastHandler) GetContent(c echo.Context) error {
	id := c.Param("id")
	response := ph.PodcastService.GetContent(id)
	return c.JSON(response.HTTPStatus, response)
}

// UpdateContent godoc
// @Summary تحديث محتوى
// @Description تحديث بودكاست بواسطة المعرف
// @Tags إدارة المحتوى
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "معرف البودكاست"
// @Param request body podcastDTOs.UpdateContentRequest true "بيانات التحديث"
// @Success 200 {object} base.Response{data=podcastDTOs.PodcastResponse} "تم تحديث البودكاست بنجاح"
// @Failure 400 {object} base.Response "خطأ في البيانات المدخلة"
// @Failure 401 {object} base.Response "غير مصرح"
// @Failure 403 {object} base.Response "ليس لديك الصلاحيات الكافية"
// @Failure 404 {object} base.Response "البودكاست غير موجود"
// @Failure 500 {object} base.Response "خطأ في الخادم"
// @Router /cms/update-content/{id} [put]
func (ph *PodcastHandler) UpdateContent(c echo.Context) error {
	id := c.Param("id")

	var dto podcastDTOs.UpdateContentRequest
	if res, ok := base.BindAndValidate(c, &dto); !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	response := ph.PodcastService.UpdateContent(id, dto)
	return c.JSON(response.HTTPStatus, response)
}

// DeleteContent godoc
// @Summary حذف محتوى
// @Description حذف بودكاست بواسطة المعرف
// @Tags إدارة المحتوى
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "معرف البودكاست"
// @Success 200 {object} base.Response "تم حذف البودكاست بنجاح"
// @Failure 400 {object} base.Response "معرف البودكاست غير صالح"
// @Failure 401 {object} base.Response "غير مصرح"
// @Failure 403 {object} base.Response "ليس لديك الصلاحيات الكافية"
// @Failure 404 {object} base.Response "البودكاست غير موجود"
// @Failure 500 {object} base.Response "خطأ في الخادم"
// @Router /cms/delete-content/{id} [delete]
func (ph *PodcastHandler) DeleteContent(c echo.Context) error {
	id := c.Param("id")
	response := ph.PodcastService.DeleteContent(id)
	return c.JSON(response.HTTPStatus, response)
}

// MyContent godoc
// @Summary الحصول على المحتوى الخاص بي
// @Description استرجاع جميع البودكاست التي أنشأها المستخدم الحالي
// @Tags إدارة المحتوى
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} base.Response{data=[]podcastDTOs.PodcastResponse} "تم استرجاع المحتوى الخاص بك بنجاح"
// @Failure 401 {object} base.Response "غير مصرح"
// @Failure 500 {object} base.Response "خطأ في الخادم"
// @Router /cms/my-content [post]
func (ph *PodcastHandler) MyContent(c echo.Context) error {
	userID := c.Get("user_id").(string)
	response := ph.PodcastService.GetMyContent(userID)
	return c.JSON(response.HTTPStatus, response)
}
