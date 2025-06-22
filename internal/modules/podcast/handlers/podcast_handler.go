package podcast

import (
	"github.com/labstack/echo/v4"
	"github.com/ziyadrw/faslah/internal/base"
	podcastDTOs "github.com/ziyadrw/faslah/internal/modules/podcast/dtos"
	podcastServices "github.com/ziyadrw/faslah/internal/modules/podcast/services"
)

type PodcastHandler struct {
	PodcastService *podcastServices.PodcastService
}

func NewPodcastHandler(podcastService *podcastServices.PodcastService) *PodcastHandler {
	return &PodcastHandler{
		PodcastService: podcastService,
	}
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
