package jobs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"rowjak/website-inspektorat/app/helpers"
	"rowjak/website-inspektorat/app/models"

	"github.com/goravel/framework/facades"
)

type ConvertImageJob struct {
}

// Signature The name and signature of the job.
func (receiver *ConvertImageJob) Signature() string {
	return "convert_image_job"
}

// Handle Execute the job.
func (receiver *ConvertImageJob) Handle(args ...any) error {
	postID := args[0].(int)
	tempFileName := args[1].(string)
	finalPath := args[2].(string)
	fileNameOnly := args[3].(string)
	log.Printf("Processing image conversion for Post ID: %d, Original: %s", postID, tempFileName)

	webpFilename := fileNameOnly + ".webp"
	tempPath := "temp"

	// Buka file sementara
	tempFullPath := filepath.Join("storage/app/", tempPath, tempFileName)
	srcFile, err := os.Open(tempFullPath)
	if err != nil {
		log.Printf("Failed to open temp file: %v", err)
		return fmt.Errorf("failed to open temp file: %w", err)
	}
	defer srcFile.Close()

	// Konversi dan simpan sebagai .webp
	if err := helpers.ConvertAndSaveAsWebp(srcFile, webpFilename, finalPath); err != nil {
		log.Printf("Failed to convert to .webp: %v", err)
		return fmt.Errorf("failed to convert to .webp: %w", err)
	}

	// Simpan ke database
	postImage := models.PostImage{
		PostID:  uint(postID),
		ImageLg: webpFilename,
	}

	if err := facades.Orm().Query().Create(&postImage); err != nil {
		log.Printf("Failed to save image to database: %v", err)
		return fmt.Errorf("failed to save image to database: %w", err)
	}

	// Hapus file sementara
	if err := os.Remove(tempFullPath); err != nil {
		log.Printf("Warning: Failed to remove temp file %s: %v", tempFullPath, err)
	}

	log.Printf("Successfully converted image for Post ID: %d", postID)
	return nil
}
