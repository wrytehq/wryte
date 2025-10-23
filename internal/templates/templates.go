package templates

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"strings"

	"github.com/wrytehq/wryte/web"
)

type Manager struct {
	templates map[string]*template.Template
}

func templateFuncs() template.FuncMap {
	return template.FuncMap{
		// prettyJSON formats JSON bytes with indentation
		"prettyJSON": func(data []byte) (string, error) {
			var buf bytes.Buffer
			err := json.Indent(&buf, data, "", "  ")
			if err != nil {
				return string(data), err
			}
			return buf.String(), nil
		},
		// parseJSON parses JSON bytes into a Go interface
		"parseJSON": func(data []byte) (interface{}, error) {
			var result interface{}
			err := json.Unmarshal(data, &result)
			return result, err
		},
		// jsonString converts []byte to string
		"jsonString": func(data []byte) string {
			return string(data)
		},
		// jsonField extracts a field from JSON bytes
		"jsonField": func(data []byte, field string) (interface{}, error) {
			var result map[string]interface{}
			err := json.Unmarshal(data, &result)
			if err != nil {
				return nil, err
			}
			return result[field], nil
		},
		// dict creates a map from key-value pairs for passing to templates
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("dict expects even number of arguments")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
	}
}

func New() (*Manager, error) {
	m := &Manager{
		templates: make(map[string]*template.Template),
	}

	if err := m.loadTemplates(); err != nil {
		return nil, fmt.Errorf("error loading templates: %w", err)
	}

	return m, nil
}

func (m *Manager) loadTemplates() error {
	// Load component templates first
	components, err := fs.Glob(web.Files, "templates/components/*.html")
	if err != nil {
		return fmt.Errorf("error finding component files: %w", err)
	}

	pages, err := fs.Glob(web.Files, "templates/**/*.html")
	if err != nil {
		return fmt.Errorf("error finding template files: %w", err)
	}

	// Also find direct templates in the templates folder
	directPages, err := fs.Glob(web.Files, "templates/*.html")
	if err != nil {
		return fmt.Errorf("error finding direct template files: %w", err)
	}

	// Combine both lists (remove duplicates)
	allPages := make(map[string]bool)
	for _, page := range pages {
		allPages[page] = true
	}
	for _, page := range directPages {
		allPages[page] = true
	}

	for page := range allPages {
		if page == "templates/layout.html" {
			continue // Skip layout.html, it's the base template
		}

		// Skip component files as they are parsed separately
		if strings.HasPrefix(page, "templates/components/") {
			continue
		}

		name := strings.TrimPrefix(page, "templates/")
		name = strings.TrimSuffix(name, ".html")

		tmpl := template.New("layout.html").Funcs(templateFuncs())

		// Parse layout, components, and the page template
		filesToParse := []string{"templates/layout.html"}
		filesToParse = append(filesToParse, components...)
		filesToParse = append(filesToParse, page)

		tmpl, err := tmpl.ParseFS(web.Files, filesToParse...)
		if err != nil {
			return fmt.Errorf("error parsing template %s: %w", page, err)
		}

		m.templates[name] = tmpl
	}

	return nil
}

func (m *Manager) Render(name string, data interface{}) (*template.Template, error) {
	tmpl, ok := m.templates[name]
	if !ok {
		return nil, fmt.Errorf("template %s not found", name)
	}
	return tmpl, nil
}

func (m *Manager) MustRender(name string) *template.Template {
	tmpl, err := m.Render(name, nil)
	if err != nil {
		panic(err)
	}
	return tmpl
}

func (m *Manager) List() []string {
	names := make([]string, 0, len(m.templates))
	for name := range m.templates {
		names = append(names, name)
	}
	return names
}
