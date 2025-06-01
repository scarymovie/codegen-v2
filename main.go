package main

import (
	"fmt"
	"generator/internal/usecase"
)

func main() {
	configDir := "./src"
	templateDir := "./templates"
	outputDir := "./out"

	parser := usecase.NewOpenAPIParser()
	codegen := usecase.NewCodeGenerator()

	useCase := usecase.NewGeneratorUseCase(parser, []usecase.CodeGenerator{*codegen})

	if err := useCase.Execute(configDir, templateDir, outputDir); err != nil {
		fmt.Println("ошибка:", err)
	}
}
