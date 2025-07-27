package controllers

import (
	"fmt"
	"html"
	"log"
	"rowjak/website-inspektorat/app/helpers"
	"rowjak/website-inspektorat/app/models"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type DokumenController struct {
	// Dependent services
}

func NewDokumenController() *DokumenController {
	return &DokumenController{
		// Inject services
	}
}

func (r *DokumenController) Index(ctx http.Context) http.Response {
	return nil
}

func (r *DokumenController) Download(ctx http.Context) http.Response {
	secretID := ctx.Request().Route("id")

	id, err := facades.Crypt().DecryptString(secretID)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": "ID file tidak valid."})
	}

	var unduhan models.Unduhan
	err = facades.Orm().Query().FindOrFail(&unduhan, id)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]any{
			"status":  false,
			"message": "Data tidak ditemukan",
		})
	}

	facades.Orm().Query().Model(&unduhan).Update("downloaded_count", unduhan.DownloadedCount+1)

	filePath := facades.Storage().Path("unduhan/" + unduhan.Jenis + "/" + unduhan.FileName)

	return ctx.Response().Download(filePath, unduhan.FileOgName+unduhan.FileMime)
}

func (r *DokumenController) Show(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")

	if ctx.Request().Header("X-Requested-With") == "XMLHttpRequest" {
		query := facades.Orm().Query().
			Where("jenis", id)

		return helpers.RenderDataTable(ctx, helpers.DataTableConfig{
			Model: &models.Unduhan{},
			Query: query,
			Search: []string{
				"file_og_name",
			},
			FormatRow: func(index int, model any) map[string]any {
				unduhan := model.(*models.Unduhan)
				status := ""
				if unduhan.Status == "Ditampilkan" {
					status = `<span class="badge bg-info">` + unduhan.Status + `</span>`
				} else {
					status = `<span class="badge bg-warning">` + unduhan.Status + `</span>`
				}

				encrypted_id, _ := facades.Crypt().EncryptString(strconv.FormatUint(uint64(unduhan.ID), 10))

				uri := "/dokumen/unduh/" + encrypted_id

				url := fmt.Sprintf(`
					<a href="%s" class="btn btn-sm btn-success"><i class="bx bx-download"></i> Unduh</a>
				`, uri)

				return map[string]any{
					"DT_RowIndex":    index + 1,
					"nama_file":      html.EscapeString(unduhan.FileOgName),
					"ukuran":         helpers.FormatFileSize(unduhan.FileSize),
					"mime":           html.EscapeString(unduhan.FileMime),
					"status":         status,
					"jumlah_unduhan": unduhan.DownloadedCount,
					"link":           url,
					"tanggal":        helpers.TanggalStringToIndo(unduhan.Tanggal),
					"action": fmt.Sprintf(`
						<button class="btn btn-sm btn-primary" onclick="ubah(%d)" data-toggle="tooltip" title="Edit">
							<i class="bx bxs-edit"></i>
						</button>
						<button class="btn btn-sm btn-danger" onclick="hapus(%d)" data-toggle="tooltip" title="Delete">
							<i class="bx bxs-trash"></i>
						</button>
					`, unduhan.ID, unduhan.ID),
				}
			},
		})
	}

	menuTipe := helpers.CapitalizeEachWords(id)
	meta := helpers.DefaultMeta()
	meta.Title = "Data " + menuTipe + " Inspektorat Daerah"
	meta.Description = "Data " + menuTipe + " Inspektorat Daerah"

	return ctx.Response().View().Make("dokumen/index.tmpl", map[string]any{
		"Meta":     meta,
		"MenuTipe": menuTipe,
		"Slug":     id,
	})
}

func (r *DokumenController) Store(ctx http.Context) http.Response {
	validator, err := ctx.Request().Validate(map[string]string{
		"jenis":     "required|in:renja,renstra,lkjip,iku,perjanjian-kinerja",
		"nama_file": "required",
		"status":    "required|in:Ditampilkan,Disembunyikan",
		"file":      "required|file",
		"tanggal":   "required",
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
	nama_file := ctx.Request().Input("nama_file")
	status := ctx.Request().Input("status")
	jenis := ctx.Request().Input("jenis")

	uuidString := uuid.New().String()

	fileOgName := ""
	fileName := ""
	fileSize := int64(0)
	originalExt := ""

	unduhanFile, err := ctx.Request().File("file")
	if err != nil || unduhanFile == nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"status":  false,
			"message": "No file uploaded",
		})
	} else {
		originalExt = unduhanFile.GetClientOriginalExtension()
		fileOgName = nama_file
		fileName = uuidString + "." + originalExt
		path := "unduhan/" + jenis

		if _, err := unduhanFile.StoreAs(path, fileName); err != nil {
			facades.Log().Errorf("Gagal menyimpan file: %v", err)
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{
				"status":  false,
				"message": "Gagal menyimpan file attachment.",
			})
		}

		if fileSize, err = unduhanFile.Size(); err != nil {
			facades.Log().Errorf("Gagal mendapatkan ukuran file: %v", err)
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{
				"status":  false,
				"message": "Gagal memproses metadata file.",
			})
		}
	}

	parsed, err := time.Parse("02-01-2006", ctx.Request().Input("tanggal"))
	if err != nil {
		errorHtml := `<div class="alert alert-danger" role="alert"><div class="alert-message">Format Tanggal Bukan dd-mm-yyyy</div></div>`

		return ctx.Response().Json(http.StatusOK, http.Json{
			"status":  false,
			"message": errorHtml,
		})
	}

	// Buat instance tag dengan data yang sudah divalidasi
	unduhan := models.Unduhan{
		FileOgName: fileOgName,
		Status:     status,
		Jenis:      jenis,
		FileName:   fileName,
		FileMime:   "." + originalExt,
		FileSize:   fileSize,
		Tanggal:    parsed.Format("2006-01-02"),
	}

	// simpan ke database
	if err := facades.Orm().Query().Create(&unduhan); err != nil {
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

func (r *DokumenController) Update(ctx http.Context) http.Response {
	return nil
}

func (r *DokumenController) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")

	// Query dengan FindOrFail
	var unduhan models.Unduhan
	err := facades.Orm().Query().FindOrFail(&unduhan, id)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]any{
			"status":  false,
			"message": "Data tidak ditemukan",
		})
	}

	storagePath := fmt.Sprintf("unduhan/%s/%s", unduhan.Jenis, unduhan.FileName)

	// Hapus file fisik
	if err := facades.Storage().Disk("public").Delete(storagePath); err != nil {
		facades.Log().Error(fmt.Sprintf("Gagal menghapus file: %s, error: %v", storagePath, err))
	}

	if _, err := facades.Orm().Query().Delete(&unduhan); err != nil {
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
