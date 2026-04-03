package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/thinkingojha/go-htmx/internal/utils"
)

// HomePageData wraps nil data with PageName for the base template
type HomePageData struct {
	PageName     string
	Title        string
	Description  string
	CanonicalURL string
	OgImage      string
}

func HomeHandler(w http.ResponseWriter, r *http.Request) error {
	templates, err := utils.Templates.Templates.Clone()
	if err != nil {
		return err
	}
	templates, err = templates.ParseGlob(filepath.Join(utils.Templates.BasePath, "home", "*.html"))
	if err != nil {
		return err
	}
	data := HomePageData{
		PageName:     "home",
		Title:        "AI Platform Engineer",
		Description:  "Ankush Ojha - AI Platform Engineer shipping production LLM systems and scalable microservices. Based in New Delhi.",
		CanonicalURL: "https://ankush.fyi/",
		OgImage:      "https://unavatar.io/twitter/fyiankush",
	}
	if err = templates.ExecuteTemplate(w, "home", data); err != nil {
		return err
	}
	return nil
}
