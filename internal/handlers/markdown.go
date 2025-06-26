package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/thinkingojha/go-htmx/internal/logger"
	"github.com/thinkingojha/go-htmx/internal/utils"
)

var markdownRenderer *utils.Markdown

func init() {
	markdownRenderer = &utils.Markdown{}
	markdownRenderer.InitMarkdownRenderer()
}

func MarkdownHandler(w http.ResponseWriter, r *http.Request) error {
	// Get the markdown content from query parameter or form
	markdownContent := r.FormValue("content")
	if markdownContent == "" {
		markdownContent = r.URL.Query().Get("content")
	}

	// If no content provided, show the markdown editor form
	if markdownContent == "" {
		return renderMarkdownEditor(w, r)
	}

	// Render markdown to HTML
	htmlContent, err := markdownRenderer.RenderHTML([]byte(markdownContent))
	if err != nil {
		logger.Errorf("Failed to render markdown: %v", err)
		return err
	}

	// Check if this is an HTMX request (partial update)
	if r.Header.Get("HX-Request") == "true" {
		// Return just the rendered HTML for HTMX
		w.Header().Set("Content-Type", "text/html")
		w.Write(htmlContent)
		return nil
	}

	// Return full page with rendered markdown
	return renderMarkdownPage(w, r, string(htmlContent))
}

func renderMarkdownEditor(w http.ResponseWriter, r *http.Request) error {
	templates, err := utils.Templates.Templates.Clone()
	if err != nil {
		return err
	}

	templates, err = templates.ParseGlob(filepath.Join(utils.Templates.BasePath, "markdown", "*.html"))
	if err != nil {
		// If markdown templates don't exist, create a simple editor
		return renderSimpleMarkdownEditor(w, r)
	}

	return templates.ExecuteTemplate(w, "markdown-editor", nil)
}

func renderMarkdownPage(w http.ResponseWriter, r *http.Request, htmlContent string) error {
	templates, err := utils.Templates.Templates.Clone()
	if err != nil {
		return err
	}

	data := struct {
		Content string
	}{
		Content: htmlContent,
	}

	return templates.ExecuteTemplate(w, "base", data)
}

func renderSimpleMarkdownEditor(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
    <title>Markdown Editor</title>
    <script src="https://unpkg.com/htmx.org@1.9.8"></script>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .container { max-width: 1200px; margin: 0 auto; }
        .editor { display: flex; gap: 20px; }
        .input-section, .output-section { flex: 1; }
        textarea { width: 100%; height: 400px; padding: 10px; }
        .output { border: 1px solid #ccc; padding: 10px; height: 400px; overflow-y: auto; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Markdown Editor</h1>
        <div class="editor">
            <div class="input-section">
                <h3>Markdown Input</h3>
                <textarea 
                    name="content" 
                    hx-post="/write" 
                    hx-target="#preview"
                    hx-trigger="keyup changed delay:500ms"
                    placeholder="Enter your markdown here..."></textarea>
            </div>
            <div class="output-section">
                <h3>Preview</h3>
                <div id="preview" class="output">
                    <p>Start typing to see the preview...</p>
                </div>
            </div>
        </div>
    </div>
</body>
</html>
`))
	return nil
}
