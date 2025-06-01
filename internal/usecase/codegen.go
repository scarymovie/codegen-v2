package usecase

import (
	"bytes"
	"generator/internal/model"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type CodeGenerator struct{}

func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{}
}

func (c *CodeGenerator) Generate(templateDir string, outputDir string, data model.ParsedYAML) error {
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
			tmplData := buildTemplateData(data, relPath)
			return c.generateFile(path, targetPath[:len(targetPath)-5], tmplData)
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

// buildTemplateData формирует структуру TemplateData для передачи в шаблон
func buildTemplateData(yaml model.ParsedYAML, relPath string) model.TemplateData {
	// Формируем namespace из относительного пути шаблона
	// Например: php84_symfony6/Controller/controller.php.tmpl -> php84_symfony6\Controller
	pathParts := strings.Split(relPath, string(os.PathSeparator))
	if len(pathParts) < 2 {
		pathParts = strings.Split(relPath, "/") // fallback для unix
	}
	namespace := strings.Join(pathParts[:len(pathParts)-1], "\\")

	// Получаем operationId (берём первый попавшийся)
	var operationId string
	if paths, ok := yaml.Content["paths"].(map[string]interface{}); ok {
		for _, v := range paths {
			if pathItem, ok := v.(map[string]interface{}); ok {
				for _, op := range pathItem {
					if opMap, ok := op.(map[string]interface{}); ok {
						if id, ok := opMap["operationId"].(string); ok {
							operationId = id
							break
						}
					}
				}
			}
			if operationId != "" {
				break
			}
		}
	}
	className := capitalize(operationId)
	methodName := operationId

	// Формируем use/namespace для зависимостей (пример, можно доработать под твои нужды)
	baseNs := namespace
	baseTypeOp := namespace

	return model.TemplateData{
		Namespace:              namespace,
		ClassName:              className,
		MethodName:             methodName,
		ActionNamespace:        baseNs + "\\" + className + "Action",
		Result200Namespace:     baseNs + "\\" + className + "Result200",
		ResultDefaultNamespace: baseNs + "\\" + className + "ResultDefault",
		RawProductNamespace:    baseNs + "\\NwkRawProduct",
		ErrorOpNamespace:       baseTypeOp + "\\NwkErrorOperations",
		StringOpNamespace:      baseTypeOp + "\\StringOperations",
		Result200Class:         className + "Result200",
		ResultDefaultClass:     className + "ResultDefault",
		RawValueObjectClass:    "NwkRawProduct",
		ErrorOpClass:           "NwkErrorOperations",
		StringOpClass:          "StringOperations",
		Content:                yaml.Content,
		FileName:               yaml.FileName,
	}
}

// capitalize делает первую букву строки заглавной
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
