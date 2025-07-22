// app/http/middleware/active_menu.go
package middleware

import (
	"rowjak/website-inspektorat/app/helpers" // sesuaikan dengan module name

	"github.com/goravel/framework/contracts/http"
)

func ActiveMenu() http.Middleware {
	return func(ctx http.Context) {
		helpers.CurrentPath = ctx.Request().Path()
		ctx.Request().Next()
	}
}
