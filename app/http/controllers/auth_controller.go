package controllers

import (
	"log"
	"rowjak/website-inspektorat/app/helpers"
	"rowjak/website-inspektorat/app/models"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	//Dependent services
}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (r *AuthController) ShowLoginForm(ctx http.Context) http.Response {
	// return ctx.Response().View().Make("auth/login.tmpl", map[string]any{})
	return ctx.Response().View().Make("auth/login.tmpl", map[string]any{
		"name":               "Goravel",
		"turnstile_site_key": helpers.GetTurnstileSiteKey(),
	})
}

func (r *AuthController) Login(ctx http.Context) http.Response {
	email := ctx.Request().Input("email")
	password := ctx.Request().Input("password")
	turnstileToken := ctx.Request().Input("cf-turnstile-response")

	// Validasi input
	if email == "" || password == "" {
		return ctx.Response().View().Make("auth/login.tmpl", map[string]any{
			"error":              "Email dan password wajib diisi",
			"turnstile_site_key": helpers.GetTurnstileSiteKey(),
		})
	}

	// Validasi Turnstile
	if turnstileToken == "" {
		return ctx.Response().View().Make("auth/login.tmpl", map[string]any{
			"error":              "Verifikasi Turnstile diperlukan",
			"turnstile_site_key": helpers.GetTurnstileSiteKey(),
		})
	}

	// Verifikasi Turnstile
	clientIP := getClientIP(ctx)
	turnstileValid, err := helpers.VerifyTurnstile(turnstileToken, clientIP)
	if err != nil {
		log.Printf("Turnstile verification error: %v", err)
		return ctx.Response().View().Make("auth/login.tmpl", map[string]any{
			"error":              "Verifikasi Turnstile gagal",
			"turnstile_site_key": helpers.GetTurnstileSiteKey(),
		})
	}

	if !turnstileValid {
		return ctx.Response().View().Make("auth/login.tmpl", map[string]any{
			"error":              "Verifikasi Turnstile tidak valid",
			"turnstile_site_key": helpers.GetTurnstileSiteKey(),
		})
	}

	// Cari user berdasarkan email
	var user models.User
	err = facades.Orm().Query().Where("email", email).First(&user)
	log.Println("ORM Error:", err)
	if err != nil {
		return ctx.Response().View().Make("auth/login.tmpl", map[string]any{
			"error":              "Email tidak ditemukan",
			"turnstile_site_key": helpers.GetTurnstileSiteKey(),
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return ctx.Response().View().Make("auth/login.tmpl", map[string]any{
			"error":              "password salah",
			"turnstile_site_key": helpers.GetTurnstileSiteKey(),
		})
	}

	// log.Println("USER DATA:", user)

	ctx.Request().Session().Put("user_id", user.ID)
	ctx.Request().Session().Put("user_name", user.NamaLengkap)
	ctx.Request().Session().Put("user_email", user.Email)
	ctx.Request().Session().Put("user_role", user.Role)

	return ctx.Response().Redirect(http.StatusFound, "/dashboard")
}

func (r *AuthController) Logout(ctx http.Context) http.Response {
	ctx.Request().Session().Flush()
	return ctx.Response().Redirect(http.StatusFound, "/login")
}

func (r *AuthController) Dashboard(ctx http.Context) http.Response {
	meta := helpers.DefaultMeta()
	meta.Title = "Dashboard - Goravel Starter App"
	meta.Description = "Selamat datang di aplikasi Goravel."
	meta.URL = facades.Config().GetString("app.url") + "/dashboard"
	meta.Image = "assets/logo.avif"

	if facades.Config().GetString("app.env") == "local" {
		log.Println("LOG TEST (using facades.Config):")
	}

	return ctx.Response().View().Make("dashboard.tmpl", map[string]any{
		"Meta": meta,
	})
}

func getClientIP(ctx http.Context) string {
	// Try to get IP from X-Forwarded-For header first
	forwarded := ctx.Request().Header("X-Forwarded-For")
	if forwarded != "" {
		// Take the first IP if multiple IPs are present
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// Try X-Real-IP header
	realIP := ctx.Request().Header("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fallback to remote address
	return ctx.Request().Ip()
}
