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
		Image:       "https://yourdomain.com/assets/default-og.png",
		URL:         "https://yourdomain.com",
	}
}
