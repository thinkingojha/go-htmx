package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/thinkingojha/go-htmx/internal/logger"
	"github.com/thinkingojha/go-htmx/internal/utils"
	"gopkg.in/yaml.v3"
)

type Experience struct {
	Company     string     `json:"company"`
	Position    string     `json:"position"`
	Location    string     `json:"location"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"` // nil means current
	Description string     `json:"description"`
	Highlights  []string   `json:"highlights"`
	Skills      []string   `json:"skills"`
	Website     string     `json:"website,omitempty"`
}

type EducationItem struct {
	Institution string    `json:"institution"`
	Degree      string    `json:"degree"`
	Field       string    `json:"field"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Location    string    `json:"location"`
	Details     []string  `json:"details,omitempty"`
}

type Skill struct {
	Category string   `json:"category"`
	Items    []string `json:"items"`
}

type ExperienceData struct {
	Title       string          `json:"title"`
	Subtitle    string          `json:"subtitle"`
	Summary     string          `json:"summary"`
	Experiences []Experience    `json:"experiences"`
	Education   []EducationItem `json:"education"`
	Skills      []Skill         `json:"skills"`
}

type ExperienceYAML struct {
	Title       string               `yaml:"title"`
	Subtitle    string               `yaml:"subtitle"`
	Summary     string               `yaml:"summary"`
	Experiences []ExperienceItemYAML `yaml:"experiences"`
	Education   []EducationItemYAML  `yaml:"education"`
	Skills      []SkillYAML          `yaml:"skills"`
}

type ExperienceItemYAML struct {
	Company     string   `yaml:"company"`
	Position    string   `yaml:"position"`
	Location    string   `yaml:"location"`
	StartDate   string   `yaml:"start_date"`
	EndDate     *string  `yaml:"end_date"`
	Description string   `yaml:"description"`
	Highlights  []string `yaml:"highlights"`
	Skills      []string `yaml:"skills"`
	Website     string   `yaml:"website"`
}

type EducationItemYAML struct {
	Institution string   `yaml:"institution"`
	Degree      string   `yaml:"degree"`
	Field       string   `yaml:"field"`
	StartDate   string   `yaml:"start_date"`
	EndDate     string   `yaml:"end_date"`
	Location    string   `yaml:"location"`
	Details     []string `yaml:"details"`
}

type SkillYAML struct {
	Category string   `yaml:"category"`
	Items    []string `yaml:"items"`
}

func ExpHandler(w http.ResponseWriter, r *http.Request) error {
	templates, err := utils.Templates.Templates.Clone()
	if err != nil {
		return err
	}
	templates, err = templates.ParseGlob(filepath.Join(utils.Templates.BasePath, "info", "*.html"))
	if err != nil {
		return err
	}

	// Experience data - loads from YAML file or falls back to hardcoded data
	data := getExperienceData()

	if err = templates.ExecuteTemplate(w, "info", data); err != nil {
		return err
	}
	return nil
}

func getExperienceData() ExperienceData {
	// Try to load from YAML file first
	if data, err := loadExperienceFromYAML("experience.yaml"); err == nil {
		return data
	}

	// Fallback to hardcoded data if YAML file doesn't exist
	logger.Warn("Could not load experience.yaml, using fallback data")
	return getFallbackExperienceData()
}

func loadExperienceFromYAML(filename string) (ExperienceData, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return ExperienceData{}, err
	}

	var yamlData ExperienceYAML
	if err := yaml.Unmarshal(file, &yamlData); err != nil {
		return ExperienceData{}, err
	}

	// Convert YAML data to internal data structure
	data := ExperienceData{
		Title:    yamlData.Title,
		Subtitle: yamlData.Subtitle,
		Summary:  yamlData.Summary,
	}

	// Convert experiences
	for _, exp := range yamlData.Experiences {
		startDate, _ := time.Parse("2006-01-02", exp.StartDate)
		var endDate *time.Time
		if exp.EndDate != nil && *exp.EndDate != "" {
			if parsed, err := time.Parse("2006-01-02", *exp.EndDate); err == nil {
				endDate = &parsed
			}
		}

		data.Experiences = append(data.Experiences, Experience{
			Company:     exp.Company,
			Position:    exp.Position,
			Location:    exp.Location,
			StartDate:   startDate,
			EndDate:     endDate,
			Description: exp.Description,
			Highlights:  exp.Highlights,
			Skills:      exp.Skills,
			Website:     exp.Website,
		})
	}

	// Convert education
	for _, edu := range yamlData.Education {
		startDate, _ := time.Parse("2006-01-02", edu.StartDate)
		endDate, _ := time.Parse("2006-01-02", edu.EndDate)

		data.Education = append(data.Education, EducationItem{
			Institution: edu.Institution,
			Degree:      edu.Degree,
			Field:       edu.Field,
			StartDate:   startDate,
			EndDate:     endDate,
			Location:    edu.Location,
			Details:     edu.Details,
		})
	}

	// Convert skills
	for _, skill := range yamlData.Skills {
		data.Skills = append(data.Skills, Skill{
			Category: skill.Category,
			Items:    skill.Items,
		})
	}

	return data, nil
}

func getFallbackExperienceData() ExperienceData {
	// Fallback data if YAML file is not available
	return ExperienceData{
		Title:    "professional journey",
		Subtitle: "building impactful solutions across diverse technologies",
		Summary:  "with 3+ years of industry experience and 2+ years specializing in Go development, I've worked on transforming marketing analytics, digital procurement solutions, and scalable web applications.",

		Experiences: []Experience{
			{
				Company:     "Nielsen",
				Position:    "Software Engineer",
				Location:    "Remote, India",
				StartDate:   time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
				EndDate:     nil, // Current position
				Description: "transforming marketing analytics to make them smarter and more efficient, working on data-driven solutions that help brands understand consumer behavior.",
				Highlights: []string{
					"Built scalable microservices handling millions of data points daily",
					"Improved system performance by 40% through optimization techniques",
					"Led a team of 3 developers on a critical analytics platform",
					"Implemented real-time data processing pipelines using Go and Apache Kafka",
				},
				Skills:  []string{"Go", "Kubernetes", "Apache Kafka", "PostgreSQL", "Redis", "Docker"},
				Website: "https://www.nielsen.com",
			},
			{
				Company:     "Simfoni",
				Position:    "Backend Developer",
				Location:    "New Delhi, India",
				StartDate:   time.Date(2021, 8, 1, 0, 0, 0, 0, time.UTC),
				EndDate:     &[]time.Time{time.Date(2023, 5, 31, 0, 0, 0, 0, time.UTC)}[0],
				Description: "enhanced digital tail-spend management capabilities, developing solutions that helped organizations optimize their procurement processes.",
				Highlights: []string{
					"Developed RESTful APIs serving 10,000+ daily active users",
					"Reduced API response time by 60% through database optimization",
					"Implemented automated testing pipeline reducing bugs by 80%",
					"Built integration with 15+ third-party procurement systems",
				},
				Skills:  []string{"Go", "Python", "MySQL", "REST APIs", "AWS", "Jenkins"},
				Website: "https://simfoni.com",
			},
		},

		Education: []EducationItem{
			{
				Institution: "Delhi Technological University",
				Degree:      "Bachelor of Technology",
				Field:       "Computer Science Engineering",
				StartDate:   time.Date(2017, 8, 1, 0, 0, 0, 0, time.UTC),
				EndDate:     time.Date(2021, 6, 30, 0, 0, 0, 0, time.UTC),
				Location:    "New Delhi, India",
				Details: []string{
					"Graduated with First Class Honours",
					"Active member of coding society and tech clubs",
					"Led multiple hackathon teams to victory",
				},
			},
		},

		Skills: []Skill{
			{
				Category: "Programming Languages",
				Items:    []string{"Go", "Python", "JavaScript", "TypeScript", "Java", "C++"},
			},
			{
				Category: "Backend Technologies",
				Items:    []string{"REST APIs", "GraphQL", "gRPC", "Microservices", "Message Queues"},
			},
			{
				Category: "Databases",
				Items:    []string{"PostgreSQL", "MySQL", "Redis", "MongoDB", "InfluxDB"},
			},
			{
				Category: "Cloud & DevOps",
				Items:    []string{"AWS", "Docker", "Kubernetes", "Jenkins", "Terraform", "Prometheus"},
			},
			{
				Category: "Frontend",
				Items:    []string{"HTMX", "React", "TailwindCSS", "HTML5", "CSS3"},
			},
		},
	}
}
