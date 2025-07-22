package controllers

import (
	"fmt"
	"html"
	"log"
	"path/filepath"
	"rowjak/website-inspektorat/app/helpers"
	"rowjak/website-inspektorat/app/models"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type FileEncryptController struct {
	// Dependent services
}

func NewFileEncryptController() *FileEncryptController {
	return &FileEncryptController{
		// Inject services
	}
}

func (r *FileEncryptController) Index(ctx http.Context) http.Response {
	if ctx.Request().Header("X-Requested-With") == "XMLHttpRequest" {

		query := facades.Orm().Query()
		return helpers.RenderDataTable(ctx, helpers.DataTableConfig{
			Model:  &models.FileEncrypt{},
			Query:  query,
			Search: []string{},
			FormatRow: func(index int, model any) map[string]any {
				fileEncrypt := model.(*models.FileEncrypt)
				encryptID, err := facades.Crypt().EncryptString(strconv.FormatUint(uint64(fileEncrypt.ID), 10))
				if err != nil {
					facades.Log().Errorf("gagal mengenkripsi ID: %v", err)
				}
				return map[string]any{
					"DT_RowIndex":  index + 1,
					"file_og_name": html.EscapeString(fileEncrypt.FileOgName),
					"file_size":    helpers.FormatFileSize(fileEncrypt.FileSize),
					"file_mime":    html.EscapeString(fileEncrypt.FileMime),
					"unduh": fmt.Sprintf(`
						<a class="btn btn-sm btn-success" href="/stream/file-encrypt/%s" data-toggle="tooltip" title="Edit">
							<i class="bx bx-download"></i> Unduh
						</a>
					`, encryptID),
					"action": fmt.Sprintf(`
						<button class="btn btn-sm btn-primary" onclick="ubah(%d)" data-toggle="tooltip" title="Edit">
							<i class="bx bxs-edit"></i>
						</button>
						<button class="btn btn-sm btn-danger" onclick="hapus(%d)" data-toggle="tooltip" title="Delete">
							<i class="bx bxs-trash"></i>
						</button>
					`, fileEncrypt.ID, fileEncrypt.ID), // &d ada 2, itu untuk format ID, format ID
				}
			},
		})
	}

	meta := helpers.DefaultMeta()
	meta.Title = "Upload dan Enkripsi File"
	meta.Description = "Upload dan Enkripsi File."
	meta.URL = "https://yourdomain.com/beranda"
	meta.Image = "assets/logo.avif"

	return ctx.Response().View().Make("file_encrypt/index.tmpl", map[string]any{
		"Meta": meta,
	})
}

func (r *FileEncryptController) UnduhEncryptedFileStream(ctx http.Context) http.Response {
	secretID := ctx.Request().Route("id")

	fileID, err := facades.Crypt().DecryptString(secretID)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": "ID file tidak valid."})
	}

	var fileEncrypt models.FileEncrypt
	err = facades.Orm().Query().FindOrFail(&fileEncrypt, fileID)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]any{
			"status":  false,
			"message": "Data tidak ditemukan",
		})
	}

	encryptedFilePath := facades.Storage().Path("encrypted/" + fileEncrypt.FileEncryptName)

	if !facades.Storage().Exists("encrypted/" + fileEncrypt.FileEncryptName) {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"error": "File tidak ditemukan."})
	}

	// 1. Dekripsi file langsung ke memori (dalam bentuk []byte)
	decryptedData, err := helpers.DecryptFileToMemory(encryptedFilePath)
	if err != nil {
		facades.Log().Errorf("gagal dekripsi file ke memori: %v", err)
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": "Gagal memproses file."})
	}

	// 1. Atur header pada response SEBELUM memanggil Stream.
	ctx.Response().Header("Content-Disposition", `attachment; filename="`+fileEncrypt.FileOgName+`"`)
	ctx.Response().Header("Content-Type", "application/octet-stream")
	ctx.Response().Header("Content-Length", strconv.Itoa(len(decryptedData)))

	// 2. Buat fungsi callback yang HANYA bertugas menulis data.
	streamCallback := func(w http.StreamWriter) error {
		// Tulis data ke stream
		_, err := w.Write(decryptedData)
		if err != nil {
			facades.Log().Errorf("gagal menulis stream ke response: %v", err)
			return err // Kembalikan error untuk menandakan kegagalan
		}
		return nil // Kembalikan nil untuk menandakan sukses
	}

	// 3. Panggil Stream dengan argumen yang benar: status dan callback.
	return ctx.Response().Stream(http.StatusOK, streamCallback)
}

func (u *FileEncryptController) Store(ctx http.Context) http.Response {
	validator, err := ctx.Request().Validate(map[string]string{
		"file": "required|file",
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

	var fileOgName string
	var fileEncryptName string

	// 1. Dapatkan file dari request
	uploadedFile, err := ctx.Request().File("file")
	if err != nil || uploadedFile == nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"error": "File tidak ditemukan dalam request. Pastikan menggunakan key 'file'.",
		})
	}

	// 2. Simpan nama file asli
	fileOgName = uploadedFile.GetClientOriginalName()

	// 3. Buat nama unik untuk file sementara (sebelum enkripsi)
	// Ini untuk menghindari konflik jika ada 2 user upload file bersamaan dengan nama sama
	tempUUID := uuid.New().String()
	tempFileName := tempUUID + filepath.Ext(fileOgName)

	// 4. Simpan file yang diunggah ke folder sementara
	// StoreAs mengembalikan path relatif dari 'storage/app/'
	tempRelativePath, err := uploadedFile.StoreAs("temp", tempFileName)
	if err != nil {
		facades.Log().Errorf("gagal menyimpan file sementara: %v", err)
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": "Gagal memproses file.",
		})
	}
	// Dapatkan path lengkap dari file sementara
	fullTempPath := facades.Storage().Path(tempRelativePath)

	// 5. Siapkan nama dan path untuk file terenkripsi
	fileEncryptName = tempUUID + ".enc"
	// Path lengkap tujuan file terenkripsi
	fullDestPath := facades.Storage().Path("encrypted/" + fileEncryptName)

	// 6. Panggil helper untuk mengenkripsi file
	// Helper akan membaca dari fullTempPath, mengenkripsi, menyimpan ke fullDestPath,
	// dan kemudian menghapus file di fullTempPath.
	_, err = helpers.EncryptFileLaravelCompat(fullTempPath, fullDestPath)
	if err != nil {
		facades.Log().Errorf("gagal enkripsi file: %v", err)
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": "Gagal mengenkripsi file." + err.Error(),
		})
	}

	fileMime, err := uploadedFile.MimeType()
	if err != nil {
		facades.Log().Errorf("gagal mendapatkan mime type: %v", err)
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": "Gagal memproses metadata file.",
		})
	}

	// Panggil Size() dan periksa error
	fileSize, err := uploadedFile.Size()
	if err != nil {
		facades.Log().Errorf("gagal mendapatkan ukuran file: %v", err)
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": "Gagal memproses metadata file.",
		})
	}

	// Buat instance post
	fileEncrypt := models.FileEncrypt{
		FileOgName:      fileOgName,
		FileEncryptName: fileEncryptName,
		FilePath:        "storage/app/encrypted/" + fileEncryptName,
		FileMime:        fileMime,
		FileSize:        fileSize,
	}

	// Simpan post ke database
	if err := facades.Orm().Query().Create(&fileEncrypt); err != nil {
		if facades.Config().GetString("app.env") == "local" {
			log.Println("Failed to create post:", err.Error())
		}
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  false,
			"message": "Error saving post",
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"status":  true,
		"message": "Data berhasil disimpan!",
	})
}

func (r *FileEncryptController) Show(ctx http.Context) http.Response {
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

func (u *FileEncryptController) Update(ctx http.Context) http.Response {
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

func (u *FileEncryptController) Destroy(ctx http.Context) http.Response {
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
