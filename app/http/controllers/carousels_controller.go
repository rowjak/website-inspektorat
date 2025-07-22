package controllers

import (
	"github.com/goravel/framework/contracts/http"
)

type CarouselsController struct {
	// Dependent services
}

func NewCarouselsController() *CarouselsController {
	return &CarouselsController{
		// Inject services
	}
}

func (r *CarouselsController) Index(ctx http.Context) http.Response {
	return nil
}	

func (r *CarouselsController) Show(ctx http.Context) http.Response {
	return nil
}

func (r *CarouselsController) Store(ctx http.Context) http.Response {
	return nil
}

func (r *CarouselsController) Update(ctx http.Context) http.Response {
	return nil
}

func (r *CarouselsController) Destroy(ctx http.Context) http.Response {
	return nil
}
