package helpers

import (
	"fmt"
	"strings"
	"time"
)

var CurrentPath string

func FormatFileSize(size int64) string {
	if size == 0 {
		return "0 B"
	}
	const unit = 1024
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

	// Mengonversi ukuran ke float64 untuk pembagian desimal
	sizeFloat := float64(size)
	i := 0

	// Terus bagi dengan 1024 sampai angkanya di bawah 1024
	for sizeFloat >= unit && i < len(units)-1 {
		sizeFloat /= unit
		i++
	}

	// Format hasilnya dengan 2 angka desimal dan unit yang sesuai
	// Contoh: 206951 -> 202.10 KB
	return fmt.Sprintf("%.2f %s", sizeFloat, units[i])
}

var ViewHelpers = map[string]any{
	"today": func() string {
		return time.Now().Format("02-01-2006")
	},
	"currentPath": func() string {
		return CurrentPath
	},

	"isActive": func(menuPath string) string {
		if CurrentPath == menuPath || (menuPath != "/" && strings.HasPrefix(CurrentPath, menuPath+"/")) {
			return "active"
		}
		return ""
	},
	"formatFileSize": FormatFileSize,
}
