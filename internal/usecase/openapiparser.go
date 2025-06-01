package usecase

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"generator/internal/model"
	"gopkg.in/yaml.v3"
)

type OpenAPIParser struct{}

func NewOpenAPIParser() *OpenAPIParser {
	return &OpenAPIParser{}
}

func (p *OpenAPIParser) ParseYAMLFiles(configDir string) ([]model.ParsedYAML, error) {
	if configDir == "" {
		return nil, errors.New("configDir is required")
	}

	var parsed []model.ParsedYAML

	err := filepath.WalkDir(configDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if isYAMLFile(path) {
			yamlData, err := parseYAML(path)
			if err != nil {
				return err
			}
			parsed = append(parsed, model.ParsedYAML{
				FileName: filepath.Base(path),
				Content:  yamlData,
			})
		}
		return nil
	})

	return parsed, err
}

func isYAMLFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".yaml" || ext == ".yml"
}

func parseYAML(path string) (map[string]interface{}, error) {
	data := make(map[string]interface{})

	fileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(fileContent, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
