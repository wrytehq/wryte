package templates

import (
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/wrytehq/wryte/web"
)

type Manager struct {
	templates map[string]*template.Template
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
	pages, err := fs.Glob(web.Files, "templates/*.html")
	if err != nil {
		return fmt.Errorf("error finding template files: %w", err)
	}

	for _, page := range pages {
		if page == "templates/layout.html" {
			continue // Skip layout.html, it's the base template
		}

		name := strings.TrimSuffix(filepath.Base(page), ".html")

		tmpl, err := template.ParseFS(web.Files, "templates/layout.html", page)
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
