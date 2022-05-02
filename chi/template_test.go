package chi

import (
	"embed"
	"html/template"
	"io/fs"
	"strings"
	"testing"
)

var (
	//go:embed testdata/partials/*.tmpl
	templatePartialsTest embed.FS

	//go:embed testdata/*.tmpl
	templateViewsTest embed.FS
)

// TestParseTemplates tests that parseTemplates compiles accurately with sub-templates matching
// paths from the original filesystem 1:1.
func TestParseTemplates(t *testing.T) {
	tmpl, err := parseTemplates(templatePartialsTest, templateViewsTest, template.FuncMap{})
	if err != nil {
		t.Errorf("parseTemplates() error: %v", err)
	}

	// Walk through templatesTest FS and record expected file names for assertion in test.
	origFiles := make(map[string]bool)
	if err := fs.WalkDir(templateViewsTest, ".", func(path string, info fs.DirEntry, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".tmpl") {
			return err
		}
		origFiles[path] = true
		return nil
	}); err != nil {
		t.Errorf("creating test file map: %v", err)
	}

	// Compare root sub-templates against filesystem templates. Check for unfulfilled expectations
	// afterwards.
	for vpath := range tmpl {
		if _, ok := origFiles[vpath]; !ok {
			t.Errorf("parseTemplates() got unexpected template \"%v\"", vpath)
		}
		origFiles[vpath] = true
	}
	for p, ok := range origFiles {
		if !ok {
			t.Errorf("parseTemplates() want template \"%v\", got nil template", p)
		}
	}
}
