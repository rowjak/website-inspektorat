package controllers

import (
	"fmt"
	"html"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"rowjak/website-inspektorat/app/helpers"
	"rowjak/website-inspektorat/app/jobs"
	"rowjak/website-inspektorat/app/models"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/facades"
	"github.com/gosimple/slug"
)

type PostController struct {
	// Dependent services
}

type PostInput struct {
	Judul      string                `form:"judul" json:"judul"`
	Isi        string                `form:"isi" json:"isi"`
	Tags       string                `form:"tags" json:"tags"`
	Thumbnail  *multipart.FileHeader `form:"thumbnail" json:"thumbnail"`
	Gambar     *multipart.FileHeader `form:"gambar" json:"gambar"`
	Attachment *multipart.FileHeader `form:"attachment" json:"attachment"`
	Status     string                `form:"status" json:"status"`
	KategoriID uint                  `form:"kategori_id" json:"kategori_id"`
}

func NewPostController() *PostController {
	return &PostController{
		// Inject services
	}
}

func (r *PostController) Index(ctx http.Context) http.Response {
	if ctx.Request().Header("X-Requested-With") == "XMLHttpRequest" {
		// query := facades.Orm().Query()
		query := facades.Orm().Query().
			With("PostImage").
			With("PostKategori")

		route_berita := facades.Route().Info("home.index")

		return helpers.RenderDataTable(ctx, helpers.DataTableConfig{
			Model: &models.Post{},
			Query: query,
			Search: []string{
				"judul", "isi",
			},
			FormatRow: func(index int, model any) map[string]any {
				post := model.(*models.Post)
				var cleanedTags []string
				for _, tag := range strings.Split(post.Tags, ";") {
					tag = strings.TrimSpace(tag)
					if tag != "" {
						cleanedTags = append(cleanedTags, `<span class="badge bg-info">`+tag+`</span>`)
					}
				}
				status := ""
				if post.Status == "Ditampilkan" {
					status = `<span class="badge bg-success">` + post.Status + `</span>`
				} else {
					status = `<span class="badge bg-warning">` + post.Status + `</span>`
				}

				uri := route_berita.Path + "/" + post.Slug
				url := fmt.Sprintf(`
					<a target="_blank" href="%s" class="btn btn-sm btn-primary"><i class="bx bx-globe"></i> Lihat Berita</a>
				`, uri)

				items := []string{
					"Kategori : " + post.PostKategori.Kategori,
					"Tags : " + strings.Join(cleanedTags, " "),
				}

				kategori_tags := helpers.BuildDynamicList(items)

				thumbnail := facades.Storage().Url("berita/thumbnails/" + post.ThumbnailSm)

				return map[string]any{
					"DT_RowIndex":  index + 1,
					"judul":        html.EscapeString(helpers.GenerateCutWords(post.Judul, 20)),
					"slug":         url,
					"kategori_tag": kategori_tags,
					"kategori": fmt.Sprintf(`
						<span class="badge bg-primary">%s</span>
					`, post.PostKategori.Kategori),
					"status": status,
					"thumbnail": fmt.Sprintf(`
						<img src="%s" class="img-fluid" width="100px"/>
					`, thumbnail),
					"tanggal": helpers.TanggalStringToIndo(post.TanggalTampil),
					"action": fmt.Sprintf(`
						<button class="btn btn-sm btn-primary" onclick="ubah(%d)" data-toggle="tooltip" title="Edit">
							<i class="bx bxs-edit"></i>
						</button>
						<button class="btn btn-sm btn-danger" onclick="hapus(%d)" data-toggle="tooltip" title="Delete">
							<i class="bx bxs-trash"></i>
						</button>
					`, post.ID, post.ID),
				}
			},
		})
	}

	meta := helpers.DefaultMeta()
	meta.Title = "Data Berita"
	meta.Description = "Berita."

	return ctx.Response().View().Make("post/index.tmpl", map[string]any{
		"Meta": meta,
	})
}

func (r *PostController) Create(ctx http.Context) http.Response {
	var postKategori []models.PostKategori
	err := facades.Orm().Query().
		Select("id", "kategori").
		OrderBy("kategori", "asc").
		Get(&postKategori)

	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "Terjadi kesalahan pada server: " + err.Error(),
		})
	}
	var tag []models.Tag
	err = facades.Orm().Query().
		Select("slug", "nama").
		OrderBy("nama", "asc").
		Get(&tag)

	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "Terjadi kesalahan pada server: " + err.Error(),
		})
	}

	meta := helpers.DefaultMeta()
	meta.Title = "Publish Berita"
	meta.Description = "Publish Berita."

	return ctx.Response().View().Make("post/create.tmpl", map[string]any{
		"Meta":     meta,
		"Kategori": postKategori,
		"Tags":     tag,
	})
}

func (u *PostController) Store(ctx http.Context) http.Response {
	// panic("🚨 This should crash if executed")

	validator, err := ctx.Request().Validate(map[string]string{
		"isi":         "required|min:10",
		"tags":        "required",
		"kategori_id": "required",
		"status":      "required|in:Ditampilkan,Disembunyikan",
		"tanggal":     "required",

		"thumbnail":  "required|image",
		"gambar":     "array",
		"gambar.*":   "image",
		"attachment": "file",
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

	parsed, err := time.Parse("02-01-2006", ctx.Request().Input("tanggal"))
	if err != nil {
		errorHtml := `<div class="alert alert-danger" role="alert"><div class="alert-message">Format Tanggal Bukan dd-mm-yyyy</div></div>`

		return ctx.Response().Json(http.StatusOK, http.Json{
			"status":  false,
			"message": errorHtml,
		})
	}

	var input PostInput
	err = ctx.Request().Bind(&input)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  false,
			"message": "Error binding input data" + err.Error(),
		})
	}

	uuidString := uuid.New().String()

	thumbnailFile, err := ctx.Request().File("thumbnail")
	if err != nil || thumbnailFile == nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"status":  false,
			"message": "No file uploaded",
		})
	}

	// Buat UUID dan nama file untuk thumbnail
	originalExt := thumbnailFile.GetClientOriginalExtension()
	tempFilename := uuidString + "." + originalExt
	tempPath := "temp"

	// Simpan sementara di storage/app/temp/
	if _, err := thumbnailFile.StoreAs(tempPath, tempFilename); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  false,
			"message": "Failed to store temporary file",
		})
	}

	user_id, _ := helpers.StringToUint(ctx.Request().Session().Get("user_id"))
	slugText := slug.Make(input.Judul)

	attachmentOgName := ""
	attachmentName := ""
	attachmentMime := ""
	attachmentSize := int64(0)

	attachmentFile, err := ctx.Request().File("attachment")
	if err != nil || attachmentFile == nil {
		facades.Log().Info("Tidak ada attachment")
	} else {
		originalExt := attachmentFile.GetClientOriginalExtension()
		attachmentOgName = attachmentFile.GetClientOriginalName()
		attachmentName = uuidString + "." + originalExt
		path := "berita/attachments/"

		if _, err := attachmentFile.StoreAs(path, attachmentName); err != nil {
			facades.Log().Errorf("Gagal menyimpan file: %v", err)
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{
				"status":  false,
				"message": "Gagal menyimpan file attachment.",
			})
		}

		if attachmentMime, err = attachmentFile.MimeType(); err != nil {
			facades.Log().Errorf("Gagal mendapatkan mime type: %v", err)
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{
				"status":  false,
				"message": "Gagal memproses metadata file.",
			})
		}

		if attachmentSize, err = attachmentFile.Size(); err != nil {
			facades.Log().Errorf("Gagal mendapatkan ukuran file: %v", err)
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{
				"status":  false,
				"message": "Gagal memproses metadata file.",
			})
		}
	}

	// Simpan post ke database terlebih dahulu dengan thumbnail kosong
	post := models.Post{
		Judul:            input.Judul,
		Slug:             slugText,
		Readmore:         helpers.GenerateCutWords(input.Isi, 80),
		Isi:              input.Isi,
		KategoriID:       input.KategoriID,
		Tags:             input.Tags,
		Status:           input.Status,
		ThumbnailLg:      "", // akan diupdate oleh job
		ThumbnailSm:      "", // akan diupdate oleh job
		AttachmentOgName: attachmentOgName,
		AttachmentName:   attachmentName,
		AttachmentSize:   attachmentSize,
		AttachmentMime:   attachmentMime,
		UserID:           user_id,
		TanggalTampil:    parsed.Format("2006-01-02"),
	}

	if err := facades.Orm().Query().Create(&post); err != nil {
		if facades.Config().GetString("app.env") == "local" {
			log.Println("Failed to create post:", err.Error())
		}
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  false,
			"message": "Error saving post" + err.Error(),
		})
	}

	args := []queue.Arg{
		{Type: "int", Value: post.ID},
		{Type: "string", Value: tempFilename},
		{Type: "string", Value: uuidString},
		{Type: "string", Value: "public/berita/thumbnails/"},
	}

	if err := facades.Queue().Job(&jobs.ConvertThumbnailJob{}, args).Dispatch(); err != nil {
		log.Printf("Failed to dispatch thumbnail conversion job: %v", err)

		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  false,
			"message": "Error processing thumbnail",
		})
	}

	// Proses multiple images jika ada
	imageFiles, err := ctx.Request().Files("gambar")
	if err != nil || len(imageFiles) == 0 {
		log.Println("Tidak ada images terkirim")
	} else {
		for index, file := range imageFiles {
			if file != nil && file.GetClientOriginalName() != "" {
				// Buat nama file unik untuk setiap gambar
				originalExt := file.GetClientOriginalExtension()
				fileNameOnly := uuidString + "-" + strconv.Itoa(index+1)
				tempFilename := fileNameOnly + "." + originalExt

				// Simpan sementara
				if _, err := file.StoreAs(tempPath, tempFilename); err != nil {
					log.Printf("Gagal menyimpan file sementara: %s - %v", file.GetClientOriginalName(), err)
					continue
				}

				// Dispatch job untuk konversi gambar

				args := []queue.Arg{
					{Type: "int", Value: post.ID},
					{Type: "string", Value: tempFilename},
					{Type: "string", Value: "public/berita/images/"},
					{Type: "string", Value: fileNameOnly},
				}

				if err := facades.Queue().Job(&jobs.ConvertImageJob{}, args).Dispatch(); err != nil {
					log.Printf("Failed to dispatch image conversion job for %s: %v", file.GetClientOriginalName(), err)
					// Hapus file temp jika job gagal didispatch
					tempFullPath := filepath.Join("storage/app/", tempPath, tempFilename)
					_ = os.Remove(tempFullPath)
				}
			}
		}
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"status":  true,
		"message": "Data berhasil disimpan! Gambar sedang diproses di background.",
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
		storagePath := fmt.Sprintf("berita/images/%s", image.ImageLg)

		// Hapus file fisik
		if err := facades.Storage().Disk("public").Delete(storagePath); err != nil {
			facades.Log().Error(fmt.Sprintf("Gagal menghapus file: %s, error: %v", storagePath, err))
		}
	}

	if _, err := facades.Orm().Query().Where("post_id", post.ID).Delete(&models.PostImage{}); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "Gagal menghapus gambar terkait",
		})
	}

	storagePath := fmt.Sprintf("berita/thumbnails/%s", post.ThumbnailLg)

	// Hapus file fisik
	if err := facades.Storage().Disk("public").Delete(storagePath); err != nil {
		facades.Log().Error(fmt.Sprintf("Gagal menghapus file: %s, error: %v", storagePath, err))
	}

	storagePath = fmt.Sprintf("berita/thumbnails/%s", post.ThumbnailSm)

	// Hapus file fisik
	if err := facades.Storage().Disk("public").Delete(storagePath); err != nil {
		facades.Log().Error(fmt.Sprintf("Gagal menghapus file: %s, error: %v", storagePath, err))
	}

	if post.AttachmentName != "" {
		storagePath = fmt.Sprintf("berita/attachments/%s", post.AttachmentName)

		// Hapus file fisik
		if err := facades.Storage().Disk("public").Delete(storagePath); err != nil {
			facades.Log().Error(fmt.Sprintf("Gagal menghapus file: %s, error: %v", storagePath, err))
		}
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
