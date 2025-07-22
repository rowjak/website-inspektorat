// app/http/middleware/csrf.go
package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

func Csrf() http.Middleware {
	return func(ctx http.Context) {
		token := getOrCreateToken(ctx)

		// Inject langsung ke view share
		facades.View().Share("csrf_token", token)
		facades.View().Share("session_user_name", ctx.Request().Session().Get("user_name"))
		facades.View().Share("session_user_id", ctx.Request().Session().Get("user_id"))

		// Skip untuk GET, HEAD, OPTIONS
		if ctx.Request().Method() == "GET" ||
			ctx.Request().Method() == "HEAD" ||
			ctx.Request().Method() == "OPTIONS" {
			ctx.Request().Next()
			return
		}

		// Verifikasi token
		requestToken := ctx.Request().Input("_token")
		if requestToken == "" {
			requestToken = ctx.Request().Header("X-CSRF-Token")
		}

		if token != requestToken {
			ctx.Request().AbortWithStatusJson(419, http.Json{
				"error": "CSRF token mismatch",
			})
			return
		}

		ctx.Request().Next()
	}
}

func getOrCreateToken(ctx http.Context) string {
	if !ctx.Request().Session().Has("csrf_token") {
		token := generateToken()
		ctx.Request().Session().Put("csrf_token", token)
		return token
	}
	return ctx.Request().Session().Get("csrf_token").(string)
}

func generateToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
