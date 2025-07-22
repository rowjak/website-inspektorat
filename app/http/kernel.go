package http

import (
	"rowjak/website-inspektorat/app/http/middleware"

	"github.com/goravel/framework/contracts/http"
	sessionMiddleware "github.com/goravel/framework/session/middleware"
)

type Kernel struct {
}

// The application's global HTTP middleware stack.
// These middleware are run during every request to your application.
func (kernel Kernel) Middleware() []http.Middleware {
	return []http.Middleware{
		sessionMiddleware.StartSession(),
		middleware.Csrf(),
		middleware.ActiveMenu(),
		middleware.SessionData(),
	}
}
