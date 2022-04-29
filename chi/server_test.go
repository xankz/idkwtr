package chi

import (
	"embed"
	"html/template"
	"io/fs"
	"strings"
	"testing"
)

//go:embed testdata/*.tmpl
var templatesTest embed.FS

// TestParseTemplates tests that parseTemplates compiles accurately with sub-templates matching
// paths from the original filesystem 1:1.
func TestParseTemplates(t *testing.T) {
	tmpl, err := parseTemplates(templatesTest, template.FuncMap{})
	if err != nil {
		t.Errorf("parseTemplates() error: %v", err)
	}

	// Walk through templatesTest FS and record expected file names for assertion in test.
	origFiles := make(map[string]bool)
	if err := fs.WalkDir(templatesTest, ".", func(path string, info fs.DirEntry, err error) error {
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
	for _, t2 := range tmpl.Templates() {
		// Skip root template.
		if t2.Name() == "" {
			continue
		}

		if _, ok := origFiles[t2.Name()]; !ok {
			t.Errorf("parseTemplates() got unexpected template \"%v\"", t2.Name())
		}
		origFiles[t2.Name()] = true
	}
	for f, ok := range origFiles {
		if !ok {
			t.Errorf("parseTemplates() want template \"%v\", got nil template", f)
		}
	}
}
