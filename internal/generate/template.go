package generate

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/gabemanfroi/go_crud/internal/utils"
	"html/template"
	"io/fs"
	"os"
	"strings"
)

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

const (
	templatesDir = "templates"
)

var (
	TemplateFs embed.FS
	templates  map[string]*template.Template
)

func Process(modelName string) {
	loadTemplates()
	processTemplates(utils.GetWorkingDirectory(), createTemplateData(modelName))
}

func loadTemplates() {
	templates = make(map[string]*template.Template)

	tmplFiles, err := fs.ReadDir(TemplateFs, templatesDir)

	utils.Check(err, "failed to load templates")

	initializeTemplatesMap(tmplFiles)
}

func initializeTemplatesMap(tmplFiles []fs.DirEntry) {
	for _, tmpl := range tmplFiles {
		if tmpl.IsDir() {
			continue
		}
		parsedTemplate := template.Must(template.New("").Funcs(sprig.FuncMap()).ParseFS(TemplateFs, templatesDir+"/"+tmpl.Name()))
		templates[tmpl.Name()] = parsedTemplate
	}
}

func createTemplateData(modelName string) TemplateData {
	var data TemplateData

	data.ModelProperties = getModelProperties(modelName)
	data.CreateDTOProperties = getCreateDTOProperties(modelName, data.ModelProperties)
	data.UpdateDTOProperties = getUpdateDTOProperties(modelName, data.ModelProperties)
	data.ReadDTOProperties = getReadDTOProperties(modelName, data.ModelProperties)
	data.CreateValidatorProperties = getCreateValidatorProperties(modelName, data.ModelProperties)
	data.UpdateValidatorProperties = getUpdateValidatorProperties(modelName, data.ModelProperties)
	data.PascalCasedModelName = GetPascalCasedModelName(modelName)
	data.CamelCasedModelName = GetCamelCasedModelName(modelName)
	data.SnakeCasedModelName = modelName
	data.ModelNameAbbreviation = modelName[0:1]
	data.ModelName = modelName

	return data
}

func processTemplates(directory string, data TemplateData) {
	processModel(directory, data)
	processDTOs(data)
	processValidators(data)
	processRepositories(directory, data)
	processServices(directory, data)
	processControllers(directory, data)
	processRoutes(directory, data)
}

func processValidators(data TemplateData) {
	validatorsDir := fmt.Sprintf("%s/infra/validators/%s", utils.GetWorkingDirectory(), data.ModelName)
	utils.CreateDirectoryIfNotExists(validatorsDir)

	if data.CreateValidatorProperties != nil {
		processTemplate("validator_create.tmpl", fmt.Sprintf("%s/create.go", validatorsDir), data)
	}
	if data.UpdateValidatorProperties != nil {
		processTemplate("validator_update.tmpl", fmt.Sprintf("%s/update.go", validatorsDir), data)
	}
}

func processModel(directory string, data TemplateData) {
	if data.ModelProperties != nil {
		processTemplate("model.tmpl", fmt.Sprintf("%s/domain/models/%s.go", directory, data.ModelName), data)
	}
}

func processDTOs(data TemplateData) {
	dtosDir := fmt.Sprintf("%s/domain/DTO/%s", utils.GetWorkingDirectory(), data.ModelName)
	utils.CreateDirectoryIfNotExists(dtosDir)
	if data.UpdateDTOProperties != nil {
		processTemplate("dto_update.tmpl", fmt.Sprintf("%s/update.go", dtosDir), data)
	}
	if data.ReadDTOProperties != nil {
		processTemplate("dto_read.tmpl", fmt.Sprintf("%s/read.go", dtosDir), data)
	}
	if data.CreateDTOProperties != nil {
		processTemplate("dto_create.tmpl", fmt.Sprintf("%s/create.go", dtosDir), data)
	}
}

func processRoutes(directory string, data TemplateData) {
	processTemplate("route.tmpl", fmt.Sprintf("%s/application/routes/%s_routes.go", directory, data.ModelName), data)
}

func processControllers(directory string, data TemplateData) {
	processTemplate("controller_interface.tmpl", fmt.Sprintf("%s/domain/interfaces/controllers/%s_controller_interface.go", directory, data.ModelName), data)
	processTemplate("controller.tmpl", fmt.Sprintf("%s/application/controllers/%s_controller.go", directory, data.ModelName), data)
}

func processServices(directory string, data TemplateData) {
	processTemplate("service_interface.tmpl", fmt.Sprintf("%s/domain/interfaces/services/%s_service_interface.go", directory, data.ModelName), data)
	processTemplate("service.tmpl", fmt.Sprintf("%s/domain/services/%s_service.go", directory, data.ModelName), data)
}

func processRepositories(directory string, data TemplateData) {
	processTemplate("repository_interface.tmpl", fmt.Sprintf("%s/domain/interfaces/repositories/%s_repository_interface.go", directory, data.ModelName), data)
	processTemplate("repository.tmpl", fmt.Sprintf("%s/infra/db/repositories/%s_repository.go", directory, data.ModelName), data)
}

func processTemplate(fileName string, outputFile string, data TemplateData) {
	tmpl, ok := templates[fileName]
	if !ok {
		fmt.Println(ok)
	}

	var processed bytes.Buffer

	err := tmpl.ExecuteTemplate(&processed, fileName, data)
	utils.Check(err, "Unable to parse data into template: ")

	fmt.Println("Writing file: ", outputFile)

	createOutputFile(outputFile, processed)
}

func createOutputFile(outputFile string, processed bytes.Buffer) {
	f, _ := os.Create(outputFile)
	w := bufio.NewWriter(f)
	w.WriteString(strings.ReplaceAll(processed.String(), "&#34;", `"`))
	w.Flush()
}
