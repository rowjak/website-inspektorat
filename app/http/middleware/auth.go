package middleware

import (
	"github.com/goravel/framework/contracts/http"
)

func Auth() http.Middleware {
	return func(ctx http.Context) {
		userID := ctx.Request().Session().Get("user_id")

		if userID == nil {
			ctx.Response().Redirect(302, "/login").Abort()
			return
		}

		ctx.Request().Next()

	}
}
