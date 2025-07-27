package controllers

import (
	"fmt"
	"html"
	"log"
	"os"
	"path/filepath"
	"rowjak/website-inspektorat/app/helpers"
	"rowjak/website-inspektorat/app/models"
	"strings"

	"github.com/google/uuid"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type CarouselsController struct {
	// Dependent services
}

func NewCarouselsController() *CarouselsController {
	return &CarouselsController{
		// Inject services
	}
}

func (r *CarouselsController) Index(ctx http.Context) http.Response {
	if ctx.Request().Header("X-Requested-With") == "XMLHttpRequest" {
		query := facades.Orm().Query()

		return helpers.RenderDataTable(ctx, helpers.DataTableConfig{
			Model: &models.Carousel{},
			Query: query,
			Search: []string{
				"nama",
			},
			FormatRow: func(index int, model any) map[string]any {
				carousel := model.(*models.Carousel)
				status := ""
				if carousel.Status == "Ditampilkan" {
					status = `<span class="badge bg-success">` + carousel.Status + `</span>`
				} else {
					status = `<span class="badge bg-warning">` + carousel.Status + `</span>`
				}

				url := ""
				if carousel.Link != "" {
					url = fmt.Sprintf(`
						<a target="_blank" href="%s" class="btn btn-sm btn-primary"><i class="bx bx-link"></i> Link</a>
					`, carousel.Link)
				}

				image := facades.Storage().Url("carousel/" + carousel.ImageSm)

				return map[string]any{
					"DT_RowIndex": index + 1,
					"keterangan":  html.EscapeString(carousel.Keterangan),
					"status":      status,
					"link":        url,
					"gambar": fmt.Sprintf(`
						<img src="%s" class="img-fluid" width="100px"/>
					`, image),
					"action": fmt.Sprintf(`
						<button class="btn btn-sm btn-primary" onclick="ubah(%d)" data-toggle="tooltip" title="Edit">
							<i class="bx bxs-edit"></i>
						</button>
						<button class="btn btn-sm btn-danger" onclick="hapus(%d)" data-toggle="tooltip" title="Delete">
							<i class="bx bxs-trash"></i>
						</button>
					`, carousel.ID, carousel.ID),
				}
			},
		})
	}

	meta := helpers.DefaultMeta()
	meta.Title = "Data Carousel Beranda"
	meta.Description = "Data Carousel Beranda."
	meta.URL = "https://inspektorat.pekalongankab.go.id/admin/carousel"
	meta.Image = "assets/logo.avif"

	return ctx.Response().View().Make("carousel/index.tmpl", map[string]any{
		"Meta": meta,
	})
}

func (r *CarouselsController) Show(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")

	// Query dengan FindOrFail
	var carousel models.Carousel
	err := facades.Orm().Query().FindOrFail(&carousel, id)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]any{
			"status":  false,
			"message": "Data tidak ditemukan",
		})
	}

	return ctx.Response().Success().Json(http.Json{
		"status": true,
		"data":   carousel,
	})
}

func (r *CarouselsController) Store(ctx http.Context) http.Response {
	validator, err := ctx.Request().Validate(map[string]string{
		"keterangan": "required",
		"status":     "required|in:Ditampilkan,Disembunyikan",
		"gambar":     "required|file",
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
	keterangan := ctx.Request().Input("keterangan")
	status := ctx.Request().Input("status")
	link := ctx.Request().Input("link")

	uuidString := uuid.New().String()

	carouselFile, err := ctx.Request().File("gambar")
	if err != nil || carouselFile == nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"status":  false,
			"message": "No file uploaded",
		})
	}

	originalExt := carouselFile.GetClientOriginalExtension()
	tempFilename := uuidString + "." + originalExt
	tempPath := "temp"

	// Simpan sementara di storage/app/temp/
	if _, err := carouselFile.StoreAs(tempPath, tempFilename); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  false,
			"message": "Failed to store temporary file",
		})
	}

	webpFilename := uuidString + ".webp"
	avifFilename := uuidString + ".avif"
	finalPath := "public/carousel"

	// Buka file dari path temp
	tempFullPath := filepath.Join("storage/app/public/", tempPath, tempFilename)
	srcFile, err := os.Open(tempFullPath)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  false,
			"message": "Failed to open temporary file",
		})
	}
	defer srcFile.Close()

	// DEBUG: Log file info
	fileInfo, _ := srcFile.Stat()
	log.Printf("Processing file: %s, Size: %d bytes", tempFilename, fileInfo.Size())

	// Konversi dan simpan sebagai .webp
	if err := helpers.ConvertAndSaveAsWebp(srcFile, webpFilename, finalPath); err != nil {
		log.Printf("Failed to convert to .webp: %v", err)
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  false,
			"message": "Failed to convert to .webp",
		})
	}

	// Reset file reader to the beginning
	if _, err := srcFile.Seek(0, 0); err != nil {
		log.Printf("Failed to reset file reader: %v", err)
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  false,
			"message": "Failed to open file",
		})
	}

	// Konversi dan simpan sebagai .avif
	if err := helpers.ConvertAndSaveAsAvif(srcFile, avifFilename, finalPath); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  false,
			"message": "Failed to convert to .avif",
		})
	}

	// Buat instance tag dengan data yang sudah divalidasi
	carousel := models.Carousel{
		Keterangan: keterangan,
		Status:     status,
		Link:       link,
		ImageLg:    webpFilename,
		ImageSm:    avifFilename,
	}

	// simpan ke database
	if err := facades.Orm().Query().Create(&carousel); err != nil {
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

func (r *CarouselsController) Update(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	// Cari user berdasarkan ID
	var carousel models.Carousel
	err := facades.Orm().Query().FindOrFail(&carousel, id)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]any{
			"status":  false,
			"message": "Data tidak ditemukan",
		})
	}

	// Validasi data
	validator, err := ctx.Request().Validate(map[string]string{
		"keterangan": "required",
		"status":     "required|in:Ditampilkan,Disembunyikan",
		"gambar":     "file",
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
	keterangan := ctx.Request().Input("keterangan")
	status := ctx.Request().Input("status")
	link := ctx.Request().Input("link")

	// Update data
	carousel.Keterangan = keterangan
	carousel.Status = status
	carousel.Link = link

	carouselFile, err := ctx.Request().File("gambar")
	if err != nil || carouselFile == nil {
		log.Println("Tidak ada images terkirim")
	} else {
		uuidString := uuid.New().String()

		originalExt := carouselFile.GetClientOriginalExtension()
		tempFilename := uuidString + "." + originalExt
		tempPath := "temp"

		// Simpan sementara di storage/app/temp/
		if _, err := carouselFile.StoreAs(tempPath, tempFilename); err != nil {
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{
				"status":  false,
				"message": "Failed to store temporary file",
			})
		}

		webpFilename := uuidString + ".webp"
		avifFilename := uuidString + ".avif"
		finalPath := "public/carousel"

		// Buka file dari path temp
		tempFullPath := filepath.Join("storage/app/public/", tempPath, tempFilename)
		srcFile, err := os.Open(tempFullPath)
		if err != nil {
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{
				"status":  false,
				"message": "Failed to open temporary file",
			})
		}
		defer srcFile.Close()

		// DEBUG: Log file info
		fileInfo, _ := srcFile.Stat()
		log.Printf("Processing file: %s, Size: %d bytes", tempFilename, fileInfo.Size())

		// Konversi dan simpan sebagai .webp
		if err := helpers.ConvertAndSaveAsWebp(srcFile, webpFilename, finalPath); err != nil {
			log.Printf("Failed to convert to .webp: %v", err)
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{
				"status":  false,
				"message": "Failed to convert to .webp",
			})
		}

		// Reset file reader to the beginning
		if _, err := srcFile.Seek(0, 0); err != nil {
			log.Printf("Failed to reset file reader: %v", err)
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{
				"status":  false,
				"message": "Failed to open file",
			})
		}

		// Konversi dan simpan sebagai .avif
		if err := helpers.ConvertAndSaveAsAvif(srcFile, avifFilename, finalPath); err != nil {
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{
				"status":  false,
				"message": "Failed to convert to .avif",
			})
		}

		storagePath := fmt.Sprintf("carousel/%s", carousel.ImageLg)

		// Hapus file fisik
		if err := facades.Storage().Disk("public").Delete(storagePath); err != nil {
			facades.Log().Error(fmt.Sprintf("Gagal menghapus file: %s, error: %v", storagePath, err))
		}

		storagePath = fmt.Sprintf("carousel/%s", carousel.ImageSm)

		// Hapus file fisik
		if err := facades.Storage().Disk("public").Delete(storagePath); err != nil {
			facades.Log().Error(fmt.Sprintf("Gagal menghapus file: %s, error: %v", storagePath, err))
		}

		carousel.ImageLg = webpFilename
		carousel.ImageSm = avifFilename
	}

	// Simpan perubahan ke database
	if err := facades.Orm().Query().Save(&carousel); err != nil {
		if facades.Config().GetString("app.env") == "local" {
			log.Println("Failed to update carousel:", err.Error())
		}
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  false,
			"message": "Error updating carousel",
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"status":  true,
		"message": "Data Berhasil Diperbarui!",
	})
}

func (r *CarouselsController) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")

	// Query dengan FindOrFail
	var carousel models.Carousel
	err := facades.Orm().Query().FindOrFail(&carousel, id)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]any{
			"status":  false,
			"message": "Data tidak ditemukan",
		})
	}

	storagePath := fmt.Sprintf("carousel/%s", carousel.ImageLg)

	// Hapus file fisik
	if err := facades.Storage().Disk("public").Delete(storagePath); err != nil {
		facades.Log().Error(fmt.Sprintf("Gagal menghapus file: %s, error: %v", storagePath, err))
	}

	storagePath = fmt.Sprintf("carousel/%s", carousel.ImageSm)

	// Hapus file fisik
	if err := facades.Storage().Disk("public").Delete(storagePath); err != nil {
		facades.Log().Error(fmt.Sprintf("Gagal menghapus file: %s, error: %v", storagePath, err))
	}

	if _, err := facades.Orm().Query().Delete(&carousel); err != nil {
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
