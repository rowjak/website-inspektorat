package controllers

import (
	"rowjak/website-inspektorat/app/models"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type HomeController struct {
	// Dependent services
}

func NewHomeController() *HomeController {
	return &HomeController{
		// Inject services
	}
}

func (r *HomeController) Index(ctx http.Context) http.Response {
	var carousels []models.Carousel
	err := facades.Orm().Query().
		Select("keterangan", "image_sm", "image_lg", "link").
		Where("status", "Ditampilkan").
		OrderBy("created_at", "desc").
		Get(&carousels)

	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "Terjadi kesalahan pada server: " + err.Error(),
		})
	}

	// route := facades.Route().Info("storage")

	// return ctx.Response().Json(http.StatusInternalServerError, map[string]any{
	// 	"status":  true,
	// 	"message": route,
	// })

	return ctx.Response().View().Make("layout/client/header.tmpl", map[string]any{
		"Carousel": carousels,
	})
}
