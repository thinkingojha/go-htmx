package utils

import (
	"fmt"
	"html/template"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/russross/blackfriday/v2"
)

type TemplatesStruct struct {
	Templates *template.Template
	BasePath  string
}

var Templates TemplatesStruct

func ParseTemplates(basePath string) error {
	var err error

	// Create template with helper functions
	funcMap := template.FuncMap{
		"sub": func(a, b int) int {
			return a - b
		},
		"add": func(a, b int) int {
			return a + b
		},
		"eq": func(a, b interface{}) bool {
			return a == b
		},
		"ne": func(a, b interface{}) bool {
			return a != b
		},
		"gt": func(a, b int) bool {
			return a > b
		},
		"lt": func(a, b int) bool {
			return a < b
		},
		"divide": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"mod": func(a, b int) int {
			return a % b
		},
		"contains": func(s, substr string) bool {
			return strings.Contains(s, substr)
		},
		"hasPrefix": func(s, prefix string) bool {
			return strings.HasPrefix(s, prefix)
		},
		"hasSuffix": func(s, suffix string) bool {
			return strings.HasSuffix(s, suffix)
		},
		"replace": func(s, old, new string) string {
			return strings.ReplaceAll(s, old, new)
		},
		"lower": func(s string) string {
			return strings.ToLower(s)
		},
		"upper": func(s string) string {
			return strings.ToUpper(s)
		},
		"title": func(s string) string {
			return strings.Title(s)
		},
		"join": func(a []string, sep string) string {
			return strings.Join(a, sep)
		},
		"split": func(s, sep string) []string {
			return strings.Split(s, sep)
		},
		"date": func(t time.Time, layout string) string {
			return t.Format(layout)
		},
		"now": func() time.Time {
			return time.Now()
		},
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("dict: odd number of arguments")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict: key %v is not a string", values[i])
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"urlquery": func(s string) string {
			return url.QueryEscape(s)
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"safeCSS": func(s string) template.CSS {
			return template.CSS(s)
		},
		"safeJS": func(s string) template.JS {
			return template.JS(s)
		},
		"safeURL": func(s string) template.URL {
			return template.URL(s)
		},
		"markdownify": func(s string) template.HTML {
			output := blackfriday.Run([]byte(s))
			return template.HTML(output)
		},
		"truncate": func(text string, length int) string {
			if len(text) <= length {
				return text
			}
			return text[:length] + "..."
		},
		"substr": func(text string, start, end int) string {
			if start >= len(text) {
				return ""
			}
			if end > len(text) {
				end = len(text)
			}
			if start < 0 {
				start = 0
			}
			return text[start:end]
		},
		"slice": func() []interface{} {
			return make([]interface{}, 0)
		},
		"append": func(slice []interface{}, item interface{}) []interface{} {
			return append(slice, item)
		},
		"in": func(slice []interface{}, item interface{}) bool {
			for _, v := range slice {
				if v == item {
					return true
				}
			}
			return false
		},
		"first": func(slice []interface{}, count int) []interface{} {
			if count > len(slice) {
				count = len(slice)
			}
			return slice[:count]
		},
	}

	templates := template.New("t").Funcs(funcMap)
	if templates, err = templates.ParseGlob(filepath.Join(basePath, "common", "*.html")); err != nil {
		return err
	}

	Templates = TemplatesStruct{
		Templates: templates,
		BasePath:  basePath,
	}

	return nil
}

func (t *TemplatesStruct) AddTemplateFuncs() {
	t.Templates = t.Templates.Funcs(template.FuncMap{
		"eq":  func(a, b interface{}) bool { return a == b },
		"ne":  func(a, b interface{}) bool { return a != b },
		"gt":  func(a, b int) bool { return a > b },
		"lt":  func(a, b int) bool { return a < b },
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"divide": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"mod":       func(a, b int) int { return a % b },
		"contains":  func(s, substr string) bool { return strings.Contains(s, substr) },
		"hasPrefix": func(s, prefix string) bool { return strings.HasPrefix(s, prefix) },
		"hasSuffix": func(s, suffix string) bool { return strings.HasSuffix(s, suffix) },
		"replace":   func(s, old, new string) string { return strings.ReplaceAll(s, old, new) },
		"lower":     func(s string) string { return strings.ToLower(s) },
		"upper":     func(s string) string { return strings.ToUpper(s) },
		"title":     func(s string) string { return strings.Title(s) },
		"join":      func(a []string, sep string) string { return strings.Join(a, sep) },
		"split":     func(s, sep string) []string { return strings.Split(s, sep) },
		"truncate": func(text string, length int) string {
			if len(text) <= length {
				return text
			}
			return text[:length] + "..."
		},
		"urlquery": func(text string) string {
			return url.QueryEscape(text)
		},
		"substr": func(text string, start, end int) string {
			if start >= len(text) {
				return ""
			}
			if end > len(text) {
				end = len(text)
			}
			return text[start:end]
		},
		"date": func(t time.Time, layout string) string { return t.Format(layout) },
		"now":  func() time.Time { return time.Now() },
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("dict: odd number of arguments")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict: key %v is not a string", values[i])
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
		"safeCSS":  func(s string) template.CSS { return template.CSS(s) },
		"safeJS":   func(s string) template.JS { return template.JS(s) },
		"safeURL":  func(s string) template.URL { return template.URL(s) },
		"markdownify": func(text string) template.HTML {
			output := blackfriday.Run([]byte(text))
			return template.HTML(output)
		},
	})
}
