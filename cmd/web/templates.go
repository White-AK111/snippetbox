package main

import (
	//"White-AK111/snippetbox/pkg/forms"
	"White-AK111/snippetbox/pkg/models"
	"html/template"
	"path/filepath"
	"time"
)

// Объект шаблона.
type templateData struct {
	CurrentYear int
	Flash       string
	//Form        *forms.Form
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

// Функция преобразования даты в человекочитаемый формат.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

// Функция формирует map шаблонов (кэш).
func newTemplateCache(dir string) (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {

		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	// Возвращаем map шаблонов.
	return cache, nil
}
