package controllers

import (
	"fmt"
	"html"
	"log"
	"rowjak/website-inspektorat/app/helpers"
	"rowjak/website-inspektorat/app/models"
	"strings"

	"github.com/google/uuid"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/gosimple/slug"
)

type PostController struct {
	// Dependent services
}

func NewPostController() *PostController {
	return &PostController{
		// Inject services
	}
}

func (r *PostController) Index(ctx http.Context) http.Response {
	if ctx.Request().Header("X-Requested-With") == "XMLHttpRequest" {
		// query := facades.Orm().Query()
		query := facades.Orm().Query().With("PostImage")

		return helpers.RenderDataTable(ctx, helpers.DataTableConfig{
			Model: &models.Post{},
			Query: query,
			Search: []string{
				"judul", "isi",
			},
			FormatRow: func(index int, model any) map[string]any {
				post := model.(*models.Post)
				imageCount := len(post.PostImage)
				return map[string]any{
					"DT_RowIndex":   index + 1,
					"judul":         html.EscapeString(post.Judul),
					"slug":          html.EscapeString(post.Slug),
					"thumbnail":     html.EscapeString(post.ThumbnailName),
					"jumlah_gambar": imageCount,
					"action": fmt.Sprintf(`
						<button class="btn btn-sm btn-primary" onclick="ubah(%d)" data-toggle="tooltip" title="Edit">
							<i class="bx bxs-edit"></i>
						</button>
						<button class="btn btn-sm btn-danger" onclick="hapus(%d)" data-toggle="tooltip" title="Delete">
							<i class="bx bxs-trash"></i>
						</button>
					`, post.ID, post.ID), // &d ada 2, itu untuk format ID, format ID
				}
			},
		})
	}

	meta := helpers.DefaultMeta()
	meta.Title = "Manajemen User"
	meta.Description = "Manajemen User."
	meta.URL = "https://yourdomain.com/beranda"
	meta.Image = "assets/logo.avif"

	return ctx.Response().View().Make("post/index.tmpl", map[string]any{
		"Meta": meta,
	})
}

func (u *PostController) Store(ctx http.Context) http.Response {
	validator, err := ctx.Request().Validate(map[string]string{
		"judul":     "required|min:3",
		"isi":       "required|min:10",
		"thumbnail": "required|image", // max 2MB
		"gambar":    "array",
		"gambar.*":  "image|max:2048",
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
	judul := ctx.Request().Input("judul")
	isi := ctx.Request().Input("isi")

	// Generate slug dari judul
	slugText := slug.Make(judul)

	uuidString := uuid.New().String()

	var thumbnailOgName string
	var thumbnailName string
	thumbnailFile, err := ctx.Request().File("thumbnail")
	if err == nil && thumbnailFile != nil {
		thumbnailOgName = thumbnailFile.GetClientOriginalName()

		thumbnailName = uuidString + "." + thumbnailFile.GetClientOriginalExtension()
		thumbnailPath := "posts/thumbnails/"

		if _, err := thumbnailFile.StoreAs(thumbnailPath, thumbnailName); err != nil {
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{
				"status":  false,
				"message": "Error uploading thumbnail",
			})
		}
	}

	// Buat instance post
	post := models.Post{
		Judul:           judul,
		Slug:            slugText,
		Isi:             isi,
		ThumbnailOgName: thumbnailOgName,
		ThumbnailName:   thumbnailName,
	}

	// Simpan post ke database
	if err := facades.Orm().Query().Create(&post); err != nil {
		if facades.Config().GetString("app.env") == "local" {
			log.Println("Failed to create post:", err.Error())
		}
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  false,
			"message": "Error saving post",
		})
	}

	imageFiles, err := ctx.Request().Files("gambar")
	if err != nil || len(imageFiles) == 0 {
		log.Println("Failed to create post:", err.Error())
		// Tidak ada file dikirim atau field tidak ada — lanjutkan tanpa error
	} else {
		for _, file := range imageFiles {
			if file != nil && file.GetClientOriginalName() != "" {
				// Buat nama file unik
				imageName := uuid.New().String() + "." + file.GetClientOriginalExtension()
				imagePath := "posts/images/"

				// Simpan file
				if _, err := file.StoreAs(imagePath, imageName); err != nil {
					return ctx.Response().Json(http.StatusInternalServerError, http.Json{
						"status":  false,
						"message": "Error uploading image: " + file.GetClientOriginalName(),
					})
				}

				postImage := models.PostImage{
					PostID:      post.ID,
					ImageOgName: imageName,
					ImageName:   file.GetClientOriginalName(),
				}

				// Simpan post ke database
				if err := facades.Orm().Query().Create(&postImage); err != nil {
					if facades.Config().GetString("app.env") == "local" {
						log.Println("Failed to create post:", err.Error())
					}
					return ctx.Response().Json(http.StatusInternalServerError, http.Json{
						"status":  false,
						"message": "Error saving post",
					})
				}
			}
		}
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"status":  true,
		"message": "Data berhasil disimpan!",
	})
}

func (r *PostController) Show(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")

	// Query dengan FindOrFail
	var post models.Post
	err := facades.Orm().Query().
		With("PostImage").
		FindOrFail(&post, id)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]any{
			"status":  false,
			"message": "Data tidak ditemukan",
		})
	}

	return ctx.Response().Success().Json(http.Json{
		"status": true,
		"data":   post,
	})
}

func (u *PostController) Update(ctx http.Context) http.Response {
	// id := ctx.Request().Route("id")
	// // Cari user berdasarkan ID
	// var post models.Post
	// err := facades.Orm().Query().FindOrFail(&post, id)
	// if err != nil {
	// 	return ctx.Response().Json(http.StatusNotFound, map[string]any{
	// 		"status":  false,
	// 		"message": "Data tidak ditemukan",
	// 	})
	// }

	// // Validasi data
	// validator, err := ctx.Request().Validate(map[string]string{
	// 	"email":        "required|email",
	// 	"nama_lengkap": "required",
	// 	"role":         "required|in:staf,admin",
	// })

	// if err != nil {
	// 	return ctx.Response().Json(http.StatusInternalServerError, http.Json{
	// 		"status":  false,
	// 		"message": "Validation error",
	// 	})
	// }

	// if validator.Fails() {
	// 	var errorMessages []string
	// 	for _, messages := range validator.Errors().All() {
	// 		for _, message := range messages {
	// 			errorMessages = append(errorMessages, "<li>"+message+"</li>")
	// 		}
	// 	}

	// 	errorHtml := `<div class="alert alert-danger" role="alert"><div class="alert-message">` +
	// 		strings.Join(errorMessages, "") + `</div></div>`

	// 	return ctx.Response().Json(http.StatusOK, http.Json{
	// 		"status":  false,
	// 		"message": errorHtml,
	// 	})
	// }

	// // Ambil data dari request
	// email := ctx.Request().Input("email")
	// nama_lengkap := ctx.Request().Input("nama_lengkap")
	// role := ctx.Request().Input("role")

	// // Cek apakah email sudah digunakan user lain (kecuali user yang sedang diupdate)
	// var count int64
	// err = facades.Orm().Query().Model(&models.User{}).
	// 	Where("email = ?", email).
	// 	Where("id != ?", user.ID).
	// 	Count(&count)

	// if err != nil {
	// 	return ctx.Response().Json(http.StatusInternalServerError, http.Json{
	// 		"status":  false,
	// 		"message": "Error checking email duplication",
	// 	})
	// }

	// if count > 0 {
	// 	return ctx.Response().Json(http.StatusOK, map[string]interface{}{
	// 		"status": false,
	// 		"message": `<div class="alert alert-danger" role="alert">
	// 						<div class="alert-message"><li>Email sudah digunakan oleh user lain.</li></div>
	// 					</div>`,
	// 	})
	// }

	// // Update data user
	// user.NamaLengkap = nama_lengkap
	// user.Email = email
	// user.Role = role

	// // Simpan perubahan ke database
	// if err := facades.Orm().Query().Save(&user); err != nil {
	// 	if facades.Config().GetString("app.env") == "local" {
	// 		log.Println("Failed to update user:", err.Error())
	// 	}
	// 	return ctx.Response().Json(http.StatusInternalServerError, http.Json{
	// 		"status":  false,
	// 		"message": "Error updating user",
	// 	})
	// }

	// return ctx.Response().Json(http.StatusOK, http.Json{
	// 	"status":  true,
	// 	"message": "Data Berhasil Diperbarui!",
	// })
	return nil
}

func (u *PostController) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")

	// Query dengan FindOrFail
	var post models.Post
	err := facades.Orm().Query().
		With("PostImage").
		FindOrFail(&post, id)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]any{
			"status":  false,
			"message": "Data tidak ditemukan",
		})
	}

	// Hapus file gambar terkait
	for _, image := range post.PostImage {
		// Path relatif terhadap direktori 'storage' dari disk 'local'
		storagePath := fmt.Sprintf("posts/images/%s", image.ImageName)

		// Hapus file fisik
		if err := facades.Storage().Disk("local").Delete(storagePath); err != nil {
			facades.Log().Error(fmt.Sprintf("Gagal menghapus file: %s, error: %v", storagePath, err))
		}
	}

	if _, err := facades.Orm().Query().Where("post_id", post.ID).Delete(&models.PostImage{}); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "Gagal menghapus gambar terkait",
		})
	}

	storagePath := fmt.Sprintf("posts/thumbnails/%s", post.ThumbnailName)

	// Hapus file fisik
	if err := facades.Storage().Disk("local").Delete(storagePath); err != nil {
		facades.Log().Error(fmt.Sprintf("Gagal menghapus file: %s, error: %v", storagePath, err))
	}

	if _, err := facades.Orm().Query().Delete(&post); err != nil {
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
