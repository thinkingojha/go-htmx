package utils

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// This renders markdown blogs into html and saves the files in internal/template/blog/pages

type Markdown struct {
	renderer *html.Renderer
	parser   *parser.Parser
}

func (m *Markdown) InitMarkdownRenderer() {

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	m.parser = parser.NewWithExtensions(extensions)

	htmlWithFlags := html.CommonFlags | html.HrefTargetBlank | html.LazyLoadImages | html.Safelink
	opts := html.RendererOptions{
		Flags: htmlWithFlags,
	}
	m.renderer = html.NewRenderer(opts)
}

func (m *Markdown) RenderHTML(doc []byte) ([]byte, error) {
	document := m.parser.Parse(doc)
	return markdown.Render(document, m.renderer), nil
}
