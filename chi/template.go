package chi

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"io/ioutil"
	"strings"
)

// baseTemplateName refers to the name of the defined "base" template that will be rendered.
const baseTemplateName string = "base"

var (
	//go:embed templates/partials/*.tmpl
	templatePartials embed.FS

	//go:embed templates/*.tmpl
	templateViews embed.FS

	templateFuncMap template.FuncMap = make(template.FuncMap)
)

// templates is a wrapper over a map of templates with function execute() for convenience.
type templates map[string]*template.Template

// execute is similar to (template.Template).ExecuteTemplate() but automatically renders the named
// base template internally (which holds the final result). The consumer therefore only needs to
// pass the filepath of the view to be rendered.
func (t templates) execute(wr io.Writer, view string, data interface{}) error {
	return t[view].ExecuteTemplate(wr, baseTemplateName, data)
}

// parseTemplates creates a map of templates identified by the filenames recorded in views. Each
// template is cloned from a root template holding definitions from partials, before being appended
// with the final view template.
func parseTemplates(
	partials fs.FS,
	views fs.FS,
	funcMap template.FuncMap,
) (map[string]*template.Template, error) {
	// Walk through partials FS and record every template into root template.
	root := template.New("")
	if err := walkDirTemplate(partials, ".", func(path, text string) error {
		if _, err := root.Parse(text); err != nil {
			return fmt.Errorf("parsing %v: %w", path, err)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("scanning partials: %w", err)
	}

	// Walk through views FS and create a new template consisting of definitions cloned from the
	// partials root template + final view.
	tmap := make(map[string]*template.Template)
	if err := walkDirTemplate(views, ".", func(path, text string) error {
		r, err := root.Clone()
		if err != nil {
			return err
		}
		if _, err := r.Parse(text); err != nil {
			return fmt.Errorf("parsing %v: %w", path, err)
		}
		tmap[path] = r
		return nil

	}); err != nil {
		return nil, fmt.Errorf("scanning views: %w", err)
	}

	return tmap, nil
}

// walkDirTemplate is a simplified version of fs.WalkDir that skips directories and non-template
// files, returning the path and text contents of each record.
func walkDirTemplate(dir fs.FS, root string, fn func(path string, text string) error) (fErr error) {
	return fs.WalkDir(dir, root, func(path string, d fs.DirEntry, err error) error {
		// Skip directories and non-HTML templates.
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".tmpl") {
			return err
		}

		f, err := dir.Open(path)
		if err != nil {
			return fmt.Errorf("opening %v: %w", path, err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				fErr = err
			}
		}()

		fd, err := ioutil.ReadAll(f)
		if err != nil {
			return fmt.Errorf("reading %v: %w", path, err)
		}

		if err := fn(path, string(fd)); err != nil {
			return err
		}

		return nil
	})
}
