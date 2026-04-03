package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/thinkingojha/go-htmx/internal/utils"
)

type ContactPageData struct {
	PageName     string
	Title        string
	Description  string
	CanonicalURL string
	OgImage      string
}

func ContactHandler(w http.ResponseWriter, r *http.Request) error {
	templates, err := utils.Templates.Templates.Clone()
	if err != nil {
		return err
	}
	templates, err = templates.ParseGlob(filepath.Join(utils.Templates.BasePath, "contact", "*.html"))
	if err != nil {
		return err
	}

	data := ContactPageData{
		PageName:     "contact",
		Title:        "Contact",
		Description:  "Get in touch with Ankush Ojha, an AI Platform Engineer based in New Delhi.",
		CanonicalURL: "https://ankush.fyi/contact",
	}

	if err = templates.ExecuteTemplate(w, "contact", data); err != nil {
		return err
	}
	return nil
}
