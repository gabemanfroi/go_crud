package generate

type Property struct {
	Name     string
	DataType string
}

type TemplateData struct {
	CamelCasedModelName       string
	ModelNameAbbreviation     string
	PascalCasedModelName      string
	SnakeCasedModelName       string
	ModelName                 string
	UpdateDTOProperties       []DTOProperty
	CreateDTOProperties       []DTOProperty
	ReadDTOProperties         []DTOProperty
	CreateValidatorProperties []ValidatorProperty
	UpdateValidatorProperties []ValidatorProperty
	ModelProperties           []ModelProperty
}
type ModelProperty struct {
	Property
	Nullable      bool
	MinimumValue  uint
	MinimumLength uint
	GormString    string
}

type DTOProperty struct {
	Property
	JsonString string
}

type ValidatorProperty struct {
	Property
	ValidateString string
}
