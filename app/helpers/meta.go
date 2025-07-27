package helpers

type MetaData struct {
	Title       string
	Description string
	Keywords    string
	Image       string
	URL         string
}

func DefaultMeta() MetaData {
	return MetaData{
		Title:       "Goravel Starter - My Web App",
		Description: "Goravel starter application with clean structure.",
		Keywords:    "goravel, golang, web, starter, seo",
		Image:       "assets/logo.avif",
		URL:         "https://inspektorat.pekalongannkabg.go.id",
	}
}
