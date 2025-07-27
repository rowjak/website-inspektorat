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

type TagsController struct {
	// Dependent services
}

func NewTagsController() *TagsController {
	return &TagsController{
		// Inject services
	}
}

func (r *TagsController) Index(ctx http.Context) http.Response {
	if ctx.Request().Header("X-Requested-With") == "XMLHttpRequest" {
		query := facades.Orm().Query()

		return helpers.RenderDataTable(ctx, helpers.DataTableConfig{
			Model: &models.Tag{},
			Query: query,
			Search: []string{
				"nama",
			},
			FormatRow: func(index int, model any) map[string]any {
				tag := model.(*models.Tag)
				return map[string]any{
					"DT_RowIndex": index + 1,
					"slug": fmt.Sprintf(`
						<a href="/tags/%s"><span class="badge bg-primary">%s</span></a>
					`, html.EscapeString(tag.Slug), html.EscapeString(tag.Slug)),
					"nama": html.EscapeString(tag.Nama),
					"action": fmt.Sprintf(`
						<button class="btn btn-sm btn-primary" onclick="ubah(%d)" data-toggle="tooltip" title="Edit">
							<i class="bx bxs-edit"></i>
						</button>
						<button class="btn btn-sm btn-danger" onclick="hapus(%d)" data-toggle="tooltip" title="Delete">
							<i class="bx bxs-trash"></i>
						</button>
					`, tag.ID, tag.ID),
				}
			},
		})
	}

	meta := helpers.DefaultMeta()
	meta.Title = "Daftar Tag Berita"
	meta.Description = "Daftar Tag Berita."
	meta.URL = "https://inspektorat.pekalongankab.go.id/admin/tags"
	meta.Image = "assets/logo.avif"

	return ctx.Response().View().Make("tags/index.tmpl", map[string]any{
		"Meta": meta,
	})
}

func (r *TagsController) Show(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")

	// Query dengan FindOrFail
	var tags models.Tag
	err := facades.Orm().Query().FindOrFail(&tags, id)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]any{
			"status":  false,
			"message": "Data tidak ditemukan",
		})
	}

	return ctx.Response().Success().Json(http.Json{
		"status": true,
		"data":   tags,
	})
}

func (r *TagsController) Store(ctx http.Context) http.Response {
	validator, err := ctx.Request().Validate(map[string]string{
		"nama": "required",
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
	nama := ctx.Request().Input("nama")
	slug := slug.Make(nama)

	// Buat instance tag dengan data yang sudah divalidasi
	tags := models.Tag{
		Nama: nama,
		Slug: slug,
	}

	// simpan ke database
	if err := facades.Orm().Query().Create(&tags); err != nil {
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

func (r *TagsController) Update(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	// Cari user berdasarkan ID
	var tags models.Tag
	err := facades.Orm().Query().FindOrFail(&tags, id)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]any{
			"status":  false,
			"message": "Data tidak ditemukan",
		})
	}

	// Validasi data
	validator, err := ctx.Request().Validate(map[string]string{
		"nama": "required",
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
	nama := ctx.Request().Input("nama")
	slug := slug.Make(nama)

	var count int64
	count, err = facades.Orm().Query().Model(&models.Tag{}).
		Where("nama = ?", nama).
		Where("id != ?", tags.ID).
		Count()

	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  false,
			"message": "Error checking duplicate tag",
		})
	}

	if count > 0 {
		return ctx.Response().Json(http.StatusOK, map[string]interface{}{
			"status": false,
			"message": `<div class="alert alert-danger" role="alert">
							<div class="alert-message"><li>Tags tidak boleh sama.</li></div>
						</div>`,
		})
	}

	// Update data user
	tags.Nama = nama
	tags.Slug = slug

	// Simpan perubahan ke database
	if err := facades.Orm().Query().Save(&tags); err != nil {
		if facades.Config().GetString("app.env") == "local" {
			log.Println("Failed to update tags:", err.Error())
		}
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  false,
			"message": "Error updating tags",
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"status":  true,
		"message": "Data Berhasil Diperbarui!",
	})
}

func (r *TagsController) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")

	// Query dengan FindOrFail
	var tags models.Tag
	err := facades.Orm().Query().FindOrFail(&tags, id)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]any{
			"status":  false,
			"message": "Data tidak ditemukan",
		})
	}

	if _, err := facades.Orm().Query().Delete(&tags); err != nil {
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
