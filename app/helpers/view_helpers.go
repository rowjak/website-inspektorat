package helpers

import (
	"fmt"
	"strings"
	"time"

	"github.com/goravel/framework/facades"
)

var CurrentPath string

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
	"isOpen": func(menuPath string) string {
		if CurrentPath == menuPath || (menuPath != "/" && strings.HasPrefix(CurrentPath, menuPath+"/")) {
			return "open"
		}
		return ""
	},
	"formatFileSize": FormatFileSize,
	"route": func(routeName string) string {
		// Cek apakah ada tanda titik dalam routeName
		if strings.Contains(routeName, ".") {
			parts := strings.Split(routeName, ".")
			if len(parts) == 2 {
				baseRouteName := parts[0]
				action := parts[1]

				// Cek route info untuk base route name
				routeInfo := facades.Route().Info(baseRouteName)

				// Jika method adalah RESOURCE dan action adalah index, store, atau delete
				if routeInfo.Method == "RESOURCE" && (action == "create" || action == "edit") {
					// Return path tanpa action (misal: /users)
					// return "/" + baseRouteName
					return routeInfo.Path
				}
			}
		}

		// Fallback ke logic original
		routeInfo := facades.Route().Info(routeName)
		if routeInfo.Path != "" {
			return routeInfo.Path
		}
		return ""
	},
	"routeParams": func(routeName string, param interface{}) string {
		var path string

		// Cek apakah ada tanda titik dalam routeName
		if strings.Contains(routeName, ".") {
			parts := strings.Split(routeName, ".")
			if len(parts) == 2 {
				baseRouteName := parts[0]
				action := parts[1]

				// Cek route info untuk base route name
				routeInfo := facades.Route().Info(baseRouteName)

				// Hanya untuk action .put atau .show yang memerlukan parameter
				if routeInfo.Method == "RESOURCE" && (action == "put" || action == "show") {
					// Gunakan path original dari route info dengan action
					routeInfoWithAction := facades.Route().Info(routeName)
					path = routeInfoWithAction.Path
				} else {
					// Untuk action lainnya, gunakan path original
					routeInfoWithAction := facades.Route().Info(routeName)
					path = routeInfoWithAction.Path
				}
			}
		} else {
			// Fallback ke logic original
			routeInfo := facades.Route().Info(routeName)
			path = routeInfo.Path
		}

		if path == "" {
			return ""
		}

		// Convert parameter to string
		var paramValue string
		switch v := param.(type) {
		case string:
			paramValue = v
		case int:
			paramValue = fmt.Sprintf("%d", v)
		case int64:
			paramValue = fmt.Sprintf("%d", v)
		case float64:
			paramValue = fmt.Sprintf("%.0f", v)
		default:
			paramValue = fmt.Sprintf("%v", v)
		}

		// Cari dan replace parameter placeholder pertama
		for i := 0; i < len(path); i++ {
			if path[i] == '{' {
				// Cari closing brace
				endIndex := i + 1
				for endIndex < len(path) && path[endIndex] != '}' {
					endIndex++
				}

				if endIndex < len(path) {
					// Extract parameter name
					paramName := path[i+1 : endIndex]
					placeholder := "{" + paramName + "}"

					// Replace placeholder dengan parameter value
					return strings.Replace(path, placeholder, paramValue, 1)
				}
				break
			}
		}

		// Jika tidak ada parameter placeholder, return path asli
		return path
	},
}

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
