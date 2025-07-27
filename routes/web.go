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
	facades.Route().Static("storage", "storage/app/public").Name("storage")
	facades.Route().Static("assets", "public").Name("assets")

	authController := controllers.NewAuthController()
	userController := controllers.NewUserController()
	postController := controllers.NewPostController()
	tagsController := controllers.NewTagsController()
	kategoriController := controllers.NewKategoriController()
	carouselController := controllers.NewCarouselsController()
	dokumenController := controllers.NewDokumenController()

	fileEncryptController := controllers.NewFileEncryptController()
	homeController := clientControllers.NewHomeController()

	facades.Route().Get("/home", homeController.Index).Name("home.index")

	facades.Route().Get("/login", authController.ShowLoginForm).Name("auth.login")
	facades.Route().Post("/login", authController.Login).Name("auth.login.store")

	facades.Route().Get("/", func(ctx http.Context) http.Response {
		return ctx.Response().View().Make("layout/client/header.tmpl", map[string]any{
			"version": support.Version,
		})
	})

	facades.Route().Get("/carousel", func(ctx http.Context) http.Response {
		// route := facades.Route().Info("auth.logout")
		route := facades.Route().Info("dokumen.unduh")

		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{
			"status":  true,
			"message": route,
		})
	})

	facades.Route().Get("/dokumen/unduh/{id}", dokumenController.Download).Name("dokumen.unduh")

	facades.Route().Prefix("admin").Middleware(middleware.Auth()).Group(func(router route.Router) {
		router.Get("/logout", authController.Logout).Name("auth.logout")

		router.Get("/dashboard", authController.Dashboard).Name("dashboard.index")
		router.Get("/berita/create", postController.Create).Name("post-create")

		router.Resource("/user", userController).Name("user")
		router.Resource("/tags", tagsController).Name("tags")
		router.Resource("/kategori", kategoriController).Name("kategori")

		router.Resource("/berita", postController).Name("berita")
		router.Resource("/carousel", carouselController).Name("carousel")
		router.Resource("/dokumen", dokumenController).Name("dokumen")
		router.Get("/stream/file-encrypt/{id}", fileEncryptController.UnduhEncryptedFileStream).Name("fileEncrypt.stream")
		router.Resource("/file-encrypt", fileEncryptController)

	})
}
