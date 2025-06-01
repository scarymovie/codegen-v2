package usecase

type GeneratorUseCase struct {
	Parser   *OpenAPIParser
	CodeGens []CodeGenerator
}

func NewGeneratorUseCase(p *OpenAPIParser, gens []CodeGenerator) *GeneratorUseCase {
	return &GeneratorUseCase{
		Parser:   p,
		CodeGens: gens,
	}
}

func (g *GeneratorUseCase) Execute(configDir, templateDir, outputDir string) error {
	parsedData, err := g.Parser.ParseYAMLFiles(configDir)
	if err != nil {
		return err
	}

	for _, gen := range g.CodeGens {
		for _, d := range parsedData {
			if err := gen.Generate(templateDir, outputDir, d.Content); err != nil {
				return err
			}
		}
	}

	return nil
}
