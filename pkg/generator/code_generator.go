package generator

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/yourusername/gofsm-gen/pkg/model"
)

// CodeGenerator generates Go code from FSM models
type CodeGenerator struct {
	templates *template.Template
}

// NewCodeGenerator creates a new code generator
func NewCodeGenerator() (*CodeGenerator, error) {
	return NewCodeGeneratorWithTemplateDir("")
}

// NewCodeGeneratorWithTemplateDir creates a new code generator with a custom template directory
func NewCodeGeneratorWithTemplateDir(templateDir string) (*CodeGenerator, error) {
	if templateDir == "" {
		// Find the templates directory relative to the current working directory
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get working directory: %w", err)
		}

		// Try to find templates directory
		possiblePaths := []string{
			filepath.Join(cwd, "templates"),
			filepath.Join(cwd, "..", "..", "templates"),
			filepath.Join(cwd, "..", "..", "..", "templates"),
		}

		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				templateDir = path
				break
			}
		}

		if templateDir == "" {
			return nil, fmt.Errorf("could not find templates directory")
		}
	}

	tmpl, err := template.New("").Funcs(TemplateFuncs()).ParseGlob(filepath.Join(templateDir, "*.tmpl"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates from %s: %w", templateDir, err)
	}

	return &CodeGenerator{
		templates: tmpl,
	}, nil
}

// Generate generates code for the given FSM model
func (g *CodeGenerator) Generate(model *model.FSMModel) ([]byte, error) {
	if model == nil {
		return nil, fmt.Errorf("model cannot be nil")
	}

	if model.Package == "" {
		model.Package = "main"
	}

	var buf bytes.Buffer
	if err := g.templates.ExecuteTemplate(&buf, "state_machine.tmpl", model); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.Bytes(), nil
}

// GenerateTo generates code and writes it to the given writer
func (g *CodeGenerator) GenerateTo(model *model.FSMModel, w io.Writer) error {
	code, err := g.Generate(model)
	if err != nil {
		return err
	}

	_, err = w.Write(code)
	return err
}
