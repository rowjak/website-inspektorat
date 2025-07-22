package controllers

import (
	"github.com/goravel/framework/contracts/http"
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
	return ctx.Response().Success().Json(http.Json{
		"message": "Ini adalah HomeController di dalam folder Client",
	})
}
