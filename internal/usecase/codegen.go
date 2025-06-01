package usecase

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type CodeGenerator struct{}

func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{}
}

func (c *CodeGenerator) Generate(templateDir string, outputDir string, data any) error {
	return filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(templateDir, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(outputDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, os.ModePerm)
		}

		if strings.HasSuffix(info.Name(), ".tmpl") {
			return c.generateFile(path, targetPath[:len(targetPath)-5], data)
		}

		// просто копируем файл как есть (если не .tmpl)
		return copyRawFile(path, targetPath)
	})
}

func (c *CodeGenerator) generateFile(tmplPath, outPath string, data any) error {
	tmplContent, err := os.ReadFile(tmplPath)
	if err != nil {
		return err
	}

	tmpl, err := template.New(filepath.Base(tmplPath)).Parse(string(tmplContent))
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	// сохраняем файл
	if err := os.MkdirAll(filepath.Dir(outPath), os.ModePerm); err != nil {
		return err
	}

	return os.WriteFile(outPath, buf.Bytes(), 0644)
}

func copyRawFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}
