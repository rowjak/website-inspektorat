package routes

import (
	"github.com/goravel/framework/facades"

	"rowjak/website-inspektorat/app/http/controllers"
)

func Api() {
	userController := controllers.NewUserController()
	facades.Route().Get("/users/{id}", userController.Show)
}
