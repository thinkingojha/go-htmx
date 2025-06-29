package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/thinkingojha/go-htmx/internal/logger"
	"github.com/thinkingojha/go-htmx/internal/utils"
	"gopkg.in/yaml.v3"
)

// Blog data structures
type BlogPost struct {
	ID          string     `json:"id" yaml:"id"`
	Title       string     `json:"title" yaml:"title"`
	Slug        string     `json:"slug" yaml:"slug"`
	Excerpt     string     `json:"excerpt" yaml:"excerpt"`
	Content     string     `json:"content" yaml:"content"`
	Author      string     `json:"author" yaml:"author"`
	PublishDate time.Time  `json:"publish_date" yaml:"-"`
	UpdatedDate *time.Time `json:"updated_date,omitempty" yaml:"-"`
	Category    string     `json:"category" yaml:"category"`
	Tags        []string   `json:"tags" yaml:"tags"`
	ReadingTime int        `json:"reading_time" yaml:"reading_time"`
	Featured    bool       `json:"featured" yaml:"featured"`
	Published   bool       `json:"published" yaml:"published"`
	Meta        PostMeta   `json:"meta" yaml:"meta"`
}

type BlogPostYAML struct {
	ID          string   `yaml:"id"`
	Title       string   `yaml:"title"`
	Slug        string   `yaml:"slug"`
	Excerpt     string   `yaml:"excerpt"`
	Content     string   `yaml:"content"`
	Author      string   `yaml:"author"`
	PublishDate string   `yaml:"publish_date"`
	UpdatedDate *string  `yaml:"updated_date"`
	Category    string   `yaml:"category"`
	Tags        []string `yaml:"tags"`
	ReadingTime int      `yaml:"reading_time"`
	Featured    bool     `yaml:"featured"`
	Published   bool     `yaml:"published"`
	Meta        PostMeta `yaml:"meta"`
}

type PostMeta struct {
	Description string   `json:"description" yaml:"description"`
	Keywords    []string `json:"keywords" yaml:"keywords"`
	OGImage     string   `json:"og_image,omitempty" yaml:"og_image,omitempty"`
}

type Category struct {
	Name        string `json:"name" yaml:"name"`
	Slug        string `json:"slug" yaml:"slug"`
	Description string `json:"description" yaml:"description"`
	Color       string `json:"color" yaml:"color"`
}

type BlogMeta struct {
	Keywords []string `json:"keywords" yaml:"keywords"`
	Author   string   `json:"author" yaml:"author"`
	SiteURL  string   `json:"site_url" yaml:"site_url"`
}

type BlogData struct {
	Title       string     `json:"title" yaml:"title"`
	Subtitle    string     `json:"subtitle" yaml:"subtitle"`
	Description string     `json:"description" yaml:"description"`
	Meta        BlogMeta   `json:"meta" yaml:"meta"`
	Categories  []Category `json:"categories" yaml:"categories"`
	Posts       []BlogPost `json:"posts" yaml:"-"`
}

type BlogDataYAML struct {
	Title       string     `yaml:"title"`
	Subtitle    string     `yaml:"subtitle"`
	Description string     `yaml:"description"`
	Meta        BlogMeta   `yaml:"meta"`
	Categories  []Category `yaml:"categories"`
}

type BlogPageData struct {
	BlogData
	Post             *BlogPost
	CurrentPage      int
	TotalPages       int
	PostsPerPage     int
	SelectedTag      string
	SelectedCategory string
	FeaturedPosts    []BlogPost
	RecentPosts      []BlogPost
	RelatedPosts     []BlogPost
	AllTags          []string
}

// Main blog listing handler
func WritingsHandler(w http.ResponseWriter, r *http.Request) error {
	templates, err := utils.Templates.Templates.Clone()
	if err != nil {
		return err
	}
	templates, err = templates.ParseGlob(filepath.Join(utils.Templates.BasePath, "blog", "*.html"))
	if err != nil {
		return err
	}

	// Load blog data
	blogData, err := loadBlogData()
	if err != nil {
		logger.Errorf("Failed to load blog data: %v", err)
		return err
	}

	// Parse query parameters
	page := parseIntParam(r, "page", 1)
	tag := r.URL.Query().Get("tag")
	category := r.URL.Query().Get("category")

	// Filter posts
	posts := filterPosts(blogData.Posts, tag, category)

	// Pagination
	postsPerPage := 10
	totalPages := (len(posts) + postsPerPage - 1) / postsPerPage
	if page > totalPages && totalPages > 0 {
		page = totalPages
	}

	start := (page - 1) * postsPerPage
	end := start + postsPerPage
	if end > len(posts) {
		end = len(posts)
	}

	var paginatedPosts []BlogPost
	if start < len(posts) {
		paginatedPosts = posts[start:end]
	}

	pageData := BlogPageData{
		BlogData:         *blogData,
		CurrentPage:      page,
		TotalPages:       totalPages,
		PostsPerPage:     postsPerPage,
		SelectedTag:      tag,
		SelectedCategory: category,
		FeaturedPosts:    getFeaturedPosts(blogData.Posts),
		RecentPosts:      getRecentPosts(blogData.Posts, 5),
		AllTags:          getAllTags(blogData.Posts),
	}
	pageData.Posts = paginatedPosts

	if r.Header.Get("HX-Request") == "true" {
		return templates.ExecuteTemplate(w, "posts-list", pageData)
	}

	return templates.ExecuteTemplate(w, "blog", pageData)
}

// Individual blog post handler
func BlogPostHandler(w http.ResponseWriter, r *http.Request) error {
	templates, err := utils.Templates.Templates.Clone()
	if err != nil {
		return err
	}
	templates, err = templates.ParseGlob(filepath.Join(utils.Templates.BasePath, "blog", "*.html"))
	if err != nil {
		return err
	}

	vars := mux.Vars(r)
	slug := vars["slug"]

	blogData, err := loadBlogData()
	if err != nil {
		return err
	}

	var post *BlogPost
	for i := range blogData.Posts {
		if blogData.Posts[i].Slug == slug && blogData.Posts[i].Published {
			post = &blogData.Posts[i]
			break
		}
	}

	if post == nil {
		http.NotFound(w, r)
		return nil
	}

	pageData := BlogPageData{
		BlogData:     *blogData,
		Post:         post,
		RelatedPosts: getRelatedPosts(blogData.Posts, *post, 3),
	}

	return templates.ExecuteTemplate(w, "blog", pageData)
}

// HTMX handler for filtering posts
func BlogFilterHandler(w http.ResponseWriter, r *http.Request) error {
	templates, err := utils.Templates.Templates.Clone()
	if err != nil {
		return err
	}
	templates, err = templates.ParseGlob(filepath.Join(utils.Templates.BasePath, "blog", "*.html"))
	if err != nil {
		return err
	}

	blogData, err := loadBlogData()
	if err != nil {
		return err
	}

	tag := r.URL.Query().Get("tag")
	category := r.URL.Query().Get("category")
	page := parseIntParam(r, "page", 1)

	posts := filterPosts(blogData.Posts, tag, category)

	postsPerPage := 10
	totalPages := (len(posts) + postsPerPage - 1) / postsPerPage
	start := (page - 1) * postsPerPage
	end := start + postsPerPage
	if end > len(posts) {
		end = len(posts)
	}

	var paginatedPosts []BlogPost
	if start < len(posts) {
		paginatedPosts = posts[start:end]
	}

	pageData := BlogPageData{
		BlogData:         *blogData,
		CurrentPage:      page,
		TotalPages:       totalPages,
		PostsPerPage:     postsPerPage,
		SelectedTag:      tag,
		SelectedCategory: category,
		AllTags:          getAllTags(blogData.Posts),
	}
	pageData.Posts = paginatedPosts

	return templates.ExecuteTemplate(w, "posts-list", pageData)
}

// RSS feed handler
func BlogRSSHandler(w http.ResponseWriter, r *http.Request) error {
	blogData, err := loadBlogData()
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/rss+xml")

	// Get latest 20 posts
	posts := getRecentPosts(blogData.Posts, 20)

	rss := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
<channel>
<title>%s</title>
<description>%s</description>
<link>%s</link>
<language>en-us</language>
<lastBuildDate>%s</lastBuildDate>`,
		blogData.Title,
		blogData.Description,
		blogData.Meta.SiteURL,
		time.Now().Format(time.RFC1123Z))

	for _, post := range posts {
		rss += fmt.Sprintf(`
<item>
<title>%s</title>
<description><![CDATA[%s]]></description>
<link>%s/blog/%s</link>
<guid>%s/blog/%s</guid>
<pubDate>%s</pubDate>
<author>%s</author>
</item>`,
			post.Title,
			post.Excerpt,
			blogData.Meta.SiteURL,
			post.Slug,
			blogData.Meta.SiteURL,
			post.Slug,
			post.PublishDate.Format(time.RFC1123Z),
			post.Author)
	}

	rss += `
</channel>
</rss>`

	fmt.Fprint(w, rss)
	return nil
}

// Helper functions
func loadBlogData() (*BlogData, error) {
	// Load main blog config
	configFile, err := os.ReadFile("blogs/blogs.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read blogs.yaml: %w", err)
	}

	var yamlData BlogDataYAML
	if err := yaml.Unmarshal(configFile, &yamlData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal blogs.yaml: %w", err)
	}

	blogData := &BlogData{
		Title:       yamlData.Title,
		Subtitle:    yamlData.Subtitle,
		Description: yamlData.Description,
		Meta:        yamlData.Meta,
		Categories:  yamlData.Categories,
	}

	// Load individual posts
	postFiles, err := filepath.Glob("blogs/posts/*.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to glob posts: %w", err)
	}

	for _, postFile := range postFiles {
		postData, err := os.ReadFile(postFile)
		if err != nil {
			logger.Warnf("failed to read post file %s: %v", postFile, err)
			continue
		}

		var postYAML BlogPostYAML
		if err := yaml.Unmarshal(postData, &postYAML); err != nil {
			logger.Warnf("failed to unmarshal post file %s: %v", postFile, err)
			continue
		}

		if !postYAML.Published {
			continue
		}

		publishDate, _ := time.Parse("2006-01-02", postYAML.PublishDate)
		var updatedDate *time.Time
		if postYAML.UpdatedDate != nil && *postYAML.UpdatedDate != "" {
			if parsed, err := time.Parse("2006-01-02", *postYAML.UpdatedDate); err == nil {
				updatedDate = &parsed
			}
		}

		post := BlogPost{
			ID:          postYAML.ID,
			Title:       postYAML.Title,
			Slug:        postYAML.Slug,
			Excerpt:     postYAML.Excerpt,
			Content:     postYAML.Content,
			Author:      postYAML.Author,
			PublishDate: publishDate,
			UpdatedDate: updatedDate,
			Category:    postYAML.Category,
			Tags:        postYAML.Tags,
			ReadingTime: postYAML.ReadingTime,
			Featured:    postYAML.Featured,
			Published:   postYAML.Published,
			Meta:        postYAML.Meta,
		}
		blogData.Posts = append(blogData.Posts, post)
	}

	// Sort posts by publish date (newest first)
	sort.Slice(blogData.Posts, func(i, j int) bool {
		return blogData.Posts[i].PublishDate.After(blogData.Posts[j].PublishDate)
	})

	return blogData, nil
}

func getAllTags(posts []BlogPost) []string {
	tagSet := make(map[string]struct{})
	for _, post := range posts {
		if post.Published {
			for _, tag := range post.Tags {
				tagSet[tag] = struct{}{}
			}
		}
	}
	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}

func filterPosts(posts []BlogPost, tag, category string) []BlogPost {
	var filtered []BlogPost

	for _, post := range posts {
		if !post.Published {
			continue
		}

		if category != "" && post.Category != category {
			continue
		}

		if tag != "" {
			hasTag := false
			for _, postTag := range post.Tags {
				if postTag == tag {
					hasTag = true
					break
				}
			}
			if !hasTag {
				continue
			}
		}

		filtered = append(filtered, post)
	}

	return filtered
}

func getFeaturedPosts(posts []BlogPost) []BlogPost {
	var featured []BlogPost
	for _, post := range posts {
		if post.Published && post.Featured {
			featured = append(featured, post)
		}
	}
	return featured
}

func getRecentPosts(posts []BlogPost, limit int) []BlogPost {
	var recent []BlogPost
	count := 0
	for _, post := range posts {
		if post.Published && count < limit {
			recent = append(recent, post)
			count++
		}
	}
	return recent
}

func getRelatedPosts(posts []BlogPost, currentPost BlogPost, limit int) []BlogPost {
	var related []BlogPost
	count := 0

	for _, post := range posts {
		if post.Published && post.ID != currentPost.ID && post.Category == currentPost.Category && count < limit {
			related = append(related, post)
			count++
		}
	}

	return related
}

func parseIntParam(r *http.Request, param string, defaultValue int) int {
	value := r.URL.Query().Get(param)
	if value == "" {
		return defaultValue
	}

	if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
		return parsed
	}

	return defaultValue
}

// Helper function to get category by slug
func getCategoryBySlug(categories []Category, slug string) *Category {
	for _, cat := range categories {
		if cat.Slug == slug {
			return &cat
		}
	}
	return nil
}

// Helper function to format publish date
func (p BlogPost) FormatDate() string {
	return p.PublishDate.Format("January 2, 2006")
}

// Helper function to get reading time text
func (p BlogPost) ReadingTimeText() string {
	return fmt.Sprintf("%d min read", p.ReadingTime)
}

// Helper function to get excerpt with fallback
func (p BlogPost) GetExcerpt() string {
	if p.Excerpt != "" {
		return p.Excerpt
	}
	// Simple truncation logic if no excerpt is provided
	words := strings.Fields(p.Content)
	if len(words) > 30 {
		return strings.Join(words[:30], " ") + "..."
	}
	return p.Content
}
