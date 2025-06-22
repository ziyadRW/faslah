package discovery

import (
	"github.com/labstack/echo/v4"
	"github.com/ziyadrw/faslah/internal/base"
	discoveryDTOs "github.com/ziyadrw/faslah/internal/modules/discovery/dtos"
	discoveryServices "github.com/ziyadrw/faslah/internal/modules/discovery/services"
)

type DiscoveryHandler struct {
	DiscoveryService *discoveryServices.DiscoveryService
}

func NewDiscoveryHandler(discoveryService *discoveryServices.DiscoveryService) *DiscoveryHandler {
	return &DiscoveryHandler{
		DiscoveryService: discoveryService,
	}
}

// ListPodcasts godoc
// @Summary قائمة البودكاست
// @Description استرجاع قائمة البودكاست مع إمكانية التصفية والترتيب
// @Tags اكتشاف البودكاست
// @Accept json
// @Produce json
// @Param page query int false "رقم الصفحة" default(1)
// @Param per_page query int false "عدد العناصر في الصفحة" default(10)
// @Param sort query string false "الترتيب (newest, oldest, popular)" default(newest)
// @Param published_from query string false "تاريخ النشر من (YYYY-MM-DD)"
// @Param published_to query string false "تاريخ النشر إلى (YYYY-MM-DD)"
// @Param tag query string false "الوسم"
// @Success 200 {object} base.Response "تم استرجاع قائمة البودكاست بنجاح"
// @Failure 400 {object} base.Response "خطأ في البيانات المدخلة"
// @Failure 500 {object} base.Response "خطأ في الخادم"
// @Router /discovery [get]
func (dh *DiscoveryHandler) ListPodcasts(c echo.Context) error {
	var dto discoveryDTOs.PodcastListRequest
	dto.BindPaginationParams(c)
	if res, ok := base.BindAndValidate(c, &dto); !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	response := dh.DiscoveryService.ListPodcasts(dto)
	return c.JSON(response.HTTPStatus, response)
}

// SearchPodcasts godoc
// @Summary البحث في البودكاست
// @Description البحث في البودكاست بواسطة العنوان أو الوصف أو الوسوم
// @Tags اكتشاف البودكاست
// @Accept json
// @Produce json
// @Param page query int false "رقم الصفحة" default(1)
// @Param per_page query int false "عدد العناصر في الصفحة" default(10)
// @Param q query string true "مصطلح البحث"
// @Success 200 {object} base.Response "تم استرجاع نتائج البحث بنجاح"
// @Failure 400 {object} base.Response "خطأ في البيانات المدخلة"
// @Failure 500 {object} base.Response "خطأ في الخادم"
// @Router /discovery/search [get]
func (dh *DiscoveryHandler) SearchPodcasts(c echo.Context) error {
	var dto discoveryDTOs.PodcastSearchRequest
	dto.BindPaginationParams(c)
	if res, ok := base.BindAndValidate(c, &dto); !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	response := dh.DiscoveryService.SearchPodcasts(dto)
	return c.JSON(response.HTTPStatus, response)
}
