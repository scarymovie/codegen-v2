package model

type ParsedYAML struct {
	FileName string
	Content  map[string]interface{}
}

// TemplateData содержит все данные, которые будут переданы в шаблон
// для генерации PHP-контроллера и других файлов
type TemplateData struct {
	Namespace              string
	ClassName              string
	MethodName             string
	ActionNamespace        string
	Result200Namespace     string
	ResultDefaultNamespace string
	RawProductNamespace    string
	ErrorOpNamespace       string
	StringOpNamespace      string
	Result200Class         string
	ResultDefaultClass     string
	RawValueObjectClass    string
	ErrorOpClass           string
	StringOpClass          string
	Content                map[string]interface{}
	FileName               string
	Summary                string // описание/summary схемы
}
