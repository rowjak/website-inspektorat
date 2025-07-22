package controllers

import (
	"fmt"
	"html"
	"log"
	"rowjak/website-inspektorat/app/helpers"
	"rowjak/website-inspektorat/app/models"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	// Dependent services
}

func NewUserController() *UserController {
	return &UserController{
		// Inject services
	}
}

func (r *UserController) Index(ctx http.Context) http.Response {
	if ctx.Request().Header("X-Requested-With") == "XMLHttpRequest" {
		query := facades.Orm().Query()

		return helpers.RenderDataTable(ctx, helpers.DataTableConfig{
			Model: &models.User{},
			Query: query,
			Search: []string{
				"nama_lengkap", "email", "role",
			},
			FormatRow: func(index int, model any) map[string]any {
				user := model.(*models.User)
				return map[string]any{
					"DT_RowIndex":  index + 1,
					"nama_lengkap": html.EscapeString(user.NamaLengkap),
					"email":        html.EscapeString(user.Email),
					"role":         html.EscapeString(user.Role),
					"action": fmt.Sprintf(`
						<button class="btn btn-sm btn-primary" onclick="ubah(%d)" data-toggle="tooltip" title="Edit">
							<i class="bx bxs-edit"></i>
						</button>
						<button class="btn btn-sm btn-danger" onclick="hapus(%d)" data-toggle="tooltip" title="Delete">
							<i class="bx bxs-trash"></i>
						</button>
					`, user.ID, user.ID),
				}
			},
		})
	}

	meta := helpers.DefaultMeta()
	meta.Title = "Manajemen User"
	meta.Description = "Manajemen User."
	meta.URL = "https://yourdomain.com/beranda"
	meta.Image = "assets/logo.avif"

	return ctx.Response().View().Make("user/index.tmpl", map[string]any{
		"Meta": meta,
	})
}

func (u *UserController) Store(ctx http.Context) http.Response {
	// Validasi data
	validator, err := ctx.Request().Validate(map[string]string{
		"email":                 "required|email",
		"nama_lengkap":          "required",
		"role":                  "required|in:staf,admin",
		"password":              "required",
		"password_confirmation": "required",
	})

	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  false,
			"message": "Validation error",
		})
	}

	if validator.Fails() {
		var errorMessages []string
		for _, messages := range validator.Errors().All() {
			for _, message := range messages {
				errorMessages = append(errorMessages, "<li>"+message+"</li>")
			}
		}

		errorHtml := `<div class="alert alert-danger" role="alert"><div class="alert-message">` +
			strings.Join(errorMessages, "") + `</div></div>`

		return ctx.Response().Json(http.StatusOK, http.Json{
			"status":  false,
			"message": errorHtml,
		})
	}

	// Ambil data dari request yang sudah divalidasi
	email := ctx.Request().Input("email")
	nama_lengkap := ctx.Request().Input("nama_lengkap")
	role := ctx.Request().Input("role")
	password := ctx.Request().Input("password")
	password_confirmation := ctx.Request().Input("password_confirmation")

	// Manual check untuk confirmed
	if password != password_confirmation {
		return ctx.Response().Json(http.StatusOK, map[string]interface{}{
			"status": false,
			"message": `<div class="alert alert-danger" role="alert">
							<div class="alert-message"><li>Password tidak cocok dengan konfirmasi.</li></div>
						</div>`,
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  false,
			"message": "Error hashing password",
		})
	}

	// Buat instance user dengan data yang sudah divalidasi
	user := models.User{
		NamaLengkap: nama_lengkap,
		Email:       email,
		Role:        role,
		Password:    string(hashedPassword),
	}

	// simpan ke database
	if err := facades.Orm().Query().Create(&user); err != nil {
		if facades.Config().GetString("app.env") == "local" {
			log.Println("Failed to create user:", err.Error())
		}
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  false,
			"message": "Error saving user: ",
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"status":  true,
		"message": "Data Berhasil Disimpan!",
	})
}

func (r *UserController) Show(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")

	// Query dengan FindOrFail
	var user models.User
	err := facades.Orm().Query().FindOrFail(&user, id)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]any{
			"status":  false,
			"message": "Data tidak ditemukan",
		})
	}

	return ctx.Response().Success().Json(http.Json{
		"status": true,
		"data":   user,
	})
}

func (u *UserController) Update(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	// Cari user berdasarkan ID
	var user models.User
	err := facades.Orm().Query().FindOrFail(&user, id)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]any{
			"status":  false,
			"message": "Data tidak ditemukan",
		})
	}

	// Validasi data
	validator, err := ctx.Request().Validate(map[string]string{
		"email":        "required|email",
		"nama_lengkap": "required",
		"role":         "required|in:staf,admin",
	})

	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  false,
			"message": "Validation error",
		})
	}

	if validator.Fails() {
		var errorMessages []string
		for _, messages := range validator.Errors().All() {
			for _, message := range messages {
				errorMessages = append(errorMessages, "<li>"+message+"</li>")
			}
		}

		errorHtml := `<div class="alert alert-danger" role="alert"><div class="alert-message">` +
			strings.Join(errorMessages, "") + `</div></div>`

		return ctx.Response().Json(http.StatusOK, http.Json{
			"status":  false,
			"message": errorHtml,
		})
	}

	// Ambil data dari request
	email := ctx.Request().Input("email")
	nama_lengkap := ctx.Request().Input("nama_lengkap")
	role := ctx.Request().Input("role")
	password := ctx.Request().Input("password")
	password_confirmation := ctx.Request().Input("password_confirmation")

	var count int64
	count, err = facades.Orm().Query().Model(&models.User{}).
		Where("email = ?", email).
		Where("id != ?", user.ID).
		Count()

	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  false,
			"message": "Error checking email duplication",
		})
	}

	if count > 0 {
		return ctx.Response().Json(http.StatusOK, map[string]interface{}{
			"status": false,
			"message": `<div class="alert alert-danger" role="alert">
							<div class="alert-message"><li>Email sudah digunakan oleh user lain.</li></div>
						</div>`,
		})
	}

	// Update data user
	user.NamaLengkap = nama_lengkap
	user.Email = email
	user.Role = role

	// Jika password diisi, lakukan validasi dan hash
	if password != "" {
		// Manual check untuk password confirmation
		if password != password_confirmation {
			return ctx.Response().Json(http.StatusOK, map[string]interface{}{
				"status": false,
				"message": `<div class="alert alert-danger" role="alert">
							<div class="alert-message"><li>Password tidak cocok dengan konfirmasi.</li></div>
						</div>`,
			})
		}

		// Hash password baru
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{
				"status":  false,
				"message": "Error hashing password",
			})
		}
		user.Password = string(hashedPassword)
	}

	// Simpan perubahan ke database
	if err := facades.Orm().Query().Save(&user); err != nil {
		if facades.Config().GetString("app.env") == "local" {
			log.Println("Failed to update user:", err.Error())
		}
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  false,
			"message": "Error updating user",
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"status":  true,
		"message": "Data Berhasil Diperbarui!",
	})
}

func (u *UserController) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")

	// Query dengan FindOrFail
	var user models.User
	err := facades.Orm().Query().FindOrFail(&user, id)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]any{
			"status":  false,
			"message": "Data tidak ditemukan",
		})
	}

	if _, err := facades.Orm().Query().Delete(&user); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  false,
			"message": "Terjadi kesalahan saat menghapus data",
		})
	}

	return ctx.Response().Json(http.StatusOK, map[string]any{
		"status":  true,
		"message": "Data Berhasil Dihapus",
	})
}
