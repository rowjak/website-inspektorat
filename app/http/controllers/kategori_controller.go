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
	"github.com/gosimple/slug"
)

type KategoriController struct {
	// Dependent services
}

func NewKategoriController() *KategoriController {
	return &KategoriController{
		// Inject services
	}
}

func (r *KategoriController) Index(ctx http.Context) http.Response {
	if ctx.Request().Header("X-Requested-With") == "XMLHttpRequest" {
		query := facades.Orm().Query()

		return helpers.RenderDataTable(ctx, helpers.DataTableConfig{
			Model: &models.PostKategori{},
			Query: query,
			Search: []string{
				"kategori",
			},
			FormatRow: func(index int, model any) map[string]any {
				postKategori := model.(*models.PostKategori)
				return map[string]any{
					"DT_RowIndex": index + 1,
					"slug": fmt.Sprintf(`
						<a href="/kategori/%s"><span class="badge bg-info">%s</span></a>
					`, html.EscapeString(postKategori.Slug), html.EscapeString(postKategori.Slug)),
					"kategori": html.EscapeString(postKategori.Kategori),
					"action": fmt.Sprintf(`
						<button class="btn btn-sm btn-primary" onclick="ubah(%d)" data-toggle="tooltip" title="Edit">
							<i class="bx bxs-edit"></i>
						</button>
						<button class="btn btn-sm btn-danger" onclick="hapus(%d)" data-toggle="tooltip" title="Delete">
							<i class="bx bxs-trash"></i>
						</button>
					`, postKategori.ID, postKategori.ID),
				}
			},
		})
	}

	meta := helpers.DefaultMeta()
	meta.Title = "Daftar Kategori Berita"
	meta.Description = "Daftar Kategori Berita."
	meta.URL = "https://inspektorat.pekalongankab.go.id/admin/kategori"
	meta.Image = "assets/logo.avif"

	return ctx.Response().View().Make("kategori/index.tmpl", map[string]any{
		"Meta": meta,
	})
}

func (r *KategoriController) Show(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")

	// Query dengan FindOrFail
	var postKategori models.PostKategori
	err := facades.Orm().Query().FindOrFail(&postKategori, id)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]any{
			"status":  false,
			"message": "Data tidak ditemukan",
		})
	}

	return ctx.Response().Success().Json(http.Json{
		"status": true,
		"data":   postKategori,
	})
}

func (r *KategoriController) Store(ctx http.Context) http.Response {
	validator, err := ctx.Request().Validate(map[string]string{
		"kategori": "required",
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
	kategori := ctx.Request().Input("kategori")
	slug := slug.Make(kategori)

	// Buat instance tag dengan data yang sudah divalidasi
	postKategori := models.PostKategori{
		Kategori: kategori,
		Slug:     slug,
	}

	// simpan ke database
	if err := facades.Orm().Query().Create(&postKategori); err != nil {
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

func (r *KategoriController) Update(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	// Cari user berdasarkan ID
	var postKategori models.PostKategori
	err := facades.Orm().Query().FindOrFail(&postKategori, id)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]any{
			"status":  false,
			"message": "Data tidak ditemukan",
		})
	}

	// Validasi data
	validator, err := ctx.Request().Validate(map[string]string{
		"kategori": "required",
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
	kategori := ctx.Request().Input("kategori")
	slug := slug.Make(kategori)

	var count int64
	count, err = facades.Orm().Query().Model(&models.PostKategori{}).
		Where("kategori = ?", kategori).
		Where("id != ?", postKategori.ID).
		Count()

	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  false,
			"message": "Error checking duplicate kategori",
		})
	}

	if count > 0 {
		return ctx.Response().Json(http.StatusOK, map[string]interface{}{
			"status": false,
			"message": `<div class="alert alert-danger" role="alert">
							<div class="alert-message"><li>Kategori Berita tidak boleh sama.</li></div>
						</div>`,
		})
	}

	// Update data user
	postKategori.Kategori = kategori
	postKategori.Slug = slug

	// Simpan perubahan ke database
	if err := facades.Orm().Query().Save(&postKategori); err != nil {
		if facades.Config().GetString("app.env") == "local" {
			log.Println("Failed to update kategori:", err.Error())
		}
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  false,
			"message": "Error updating kategori",
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"status":  true,
		"message": "Data Berhasil Diperbarui!",
	})
}

func (r *KategoriController) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")

	// Query dengan FindOrFail
	var postKategori models.PostKategori
	err := facades.Orm().Query().FindOrFail(&postKategori, id)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]any{
			"status":  false,
			"message": "Data tidak ditemukan",
		})
	}

	if _, err := facades.Orm().Query().Delete(&postKategori); err != nil {
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
