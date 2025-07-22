package routes

import (
	"rowjak/website-inspektorat/app/http/controllers"
	clientControllers "rowjak/website-inspektorat/app/http/controllers/Client"
	"rowjak/website-inspektorat/app/http/middleware"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support"
)

func Web() {
	facades.Route().Static("assets", "public")

	authController := controllers.NewAuthController()
	userController := controllers.NewUserController()
	postController := controllers.NewPostController()
	fileEncryptController := controllers.NewFileEncryptController()
	homeController := clientControllers.NewHomeController()

	facades.Route().Get("/home", homeController.Index)

	facades.Route().Get("/login", authController.ShowLoginForm)
	facades.Route().Post("/login", authController.Login)

	facades.Route().Get("/", func(ctx http.Context) http.Response {
		return ctx.Response().View().Make("layout/client/header.tmpl", map[string]any{
			"version": support.Version,
		})
	})

	facades.Route().Get("/carousel", func(ctx http.Context) http.Response {
		return ctx.Response().View().Make("layout/client/carousels.tmpl", map[string]any{
			"version": support.Version,
		})
	})

	facades.Route().Middleware(middleware.Auth()).Group(func(router route.Router) {

		router.Get("/logout", authController.Logout)

		router.Get("/dashboard", authController.Dashboard)

		router.Resource("/user", userController)
		router.Resource("/post", postController)
		router.Get("/stream/file-encrypt/{id}", fileEncryptController.UnduhEncryptedFileStream)
		router.Resource("/file-encrypt", fileEncryptController)

	})
}
