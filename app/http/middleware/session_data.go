package middleware

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

func SessionData() http.Middleware {
	return func(ctx http.Context) {
		// Cek apakah ada session yang aktif
		if ctx.Request().Session().Has("user_id") {
			// Ambil data session
			userData := map[string]any{
				"id":    ctx.Request().Session().Get("user_id"),
				"name":  ctx.Request().Session().Get("user_name"),
				"email": ctx.Request().Session().Get("user_email"),
				"role":  ctx.Request().Session().Get("user_role"),
			}

			// Share data ke semua view menggunakan facades.View()
			facades.View().Share("auth_user", userData)
			facades.View().Share("is_authenticated", true)
		} else {
			// Jika tidak ada session, set default
			facades.View().Share("auth_user", nil)
			facades.View().Share("is_authenticated", false)
		}

		ctx.Request().Next()
	}
}
