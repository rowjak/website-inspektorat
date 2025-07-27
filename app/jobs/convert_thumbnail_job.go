package jobs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"rowjak/website-inspektorat/app/helpers"
	"rowjak/website-inspektorat/app/models"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	// Optional: untuk format tambahan
	"github.com/goravel/framework/facades"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp" // untuk read .webp
)

type ConvertThumbnailJob struct {
}

// Signature The name and signature of the job.
func (receiver *ConvertThumbnailJob) Signature() string {
	return "convert_thumbnail_job"
}

// Handle Execute the job.
func (receiver *ConvertThumbnailJob) Handle(args ...any) error {
	log.Println(args...)
	postID := args[0].(int)
	tempFileName := args[1].(string)
	uuidString := args[2].(string)
	finalPath := args[3].(string)
	log.Printf("Processing thumbnail conversion for Post ID: %d", postID)

	webpFilename := uuidString + ".webp"
	avifFilename := uuidString + ".avif"
	tempPath := "temp"

	// Buka file dari path temp
	tempFullPath := filepath.Join("storage/app/public/", tempPath, tempFileName)
	srcFile, err := os.Open(tempFullPath)
	if err != nil {
		log.Printf("Failed to open temp file: %v", err)
		return fmt.Errorf("failed to open temp file: %w", err)
	}
	defer srcFile.Close()

	// DEBUG: Log file info
	fileInfo, _ := srcFile.Stat()
	log.Printf("Processing file: %s, Size: %d bytes", tempFileName, fileInfo.Size())

	// Konversi dan simpan sebagai .webp
	if err := helpers.ConvertAndSaveAsWebp(srcFile, webpFilename, finalPath); err != nil {
		log.Printf("Failed to convert to .webp: %v", err)
		return fmt.Errorf("failed to convert to .webp: %w", err)
	}

	// Reset file reader to the beginning
	if _, err := srcFile.Seek(0, 0); err != nil {
		log.Printf("Failed to reset file reader: %v", err)
		return fmt.Errorf("failed to reset file reader: %w", err)
	}

	// Konversi dan simpan sebagai .avif
	if err := helpers.ConvertAndSaveAsAvif(srcFile, avifFilename, finalPath); err != nil {
		log.Printf("Failed to convert to .avif: %v", err)
		return fmt.Errorf("failed to convert to .avif: %w", err)
	}

	// Update database dengan thumbnail yang sudah dikonversi
	var post models.Post
	if err := facades.Orm().Query().Where("id", postID).First(&post); err != nil {
		log.Printf("Failed to find post: %v", err)
		return fmt.Errorf("failed to find post: %w", err)
	}

	post.ThumbnailLg = webpFilename
	post.ThumbnailSm = avifFilename

	if err := facades.Orm().Query().Save(&post); err != nil {
		log.Printf("Failed to update post thumbnail: %v", err)
		return fmt.Errorf("failed to update post thumbnail: %w", err)
	}

	// Hapus file sementara
	if err := os.Remove(tempFullPath); err != nil {
		log.Printf("Warning: Failed to remove temp file %s: %v", tempFullPath, err)
	}

	log.Printf("Successfully converted thumbnail for Post ID: %d", postID)
	return nil
}
