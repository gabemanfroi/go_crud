package generate

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/gabemanfroi/go_crud/internal/generate/validations"
	"html/template"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

type Data struct {
	CamelCasedModelName   string
	ModelNameAbbreviation string
	PascalCasedModelName  string
	SnakeCasedModelName   string
	ModelName             string
	UpdateProperties      []string
	CreateProperties      []string
	ReadProperties        []string
}

const (
	templatesDir = "templates"
)

var (
	TemplateFs embed.FS
	templates  map[string]*template.Template
)

func Process(modelName string) {
	directory := getWorkingDirectory()
	validations.CheckIfModelExists(modelName, directory)
	validations.CheckIfModelHasAllDtos(modelName, directory)
	validations.CheckIfModelHasAllValidators(modelName, directory)

	loadTemplates()
	data := createData(modelName, directory)
	processTemplates(directory, data)
}

func getWorkingDirectory() string {
	directory, err := os.Getwd()
	if err != nil {
		log.Fatal("%s", err.Error())
	}
	return directory
}

func loadTemplates() error {
	templates = make(map[string]*template.Template)

	tmplFiles, err := fs.ReadDir(TemplateFs, templatesDir)
	if err != nil {
		return err
	}

	initializeTemplatesMap(tmplFiles)
	return nil
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

func processTemplates(directory string, data Data) {
	processTemplate("repository_interface.tmpl", fmt.Sprintf("%s/domain/interfaces/repositories/%s_repository_interface.go", directory, data.ModelName), data)
	processTemplate("repository.tmpl", fmt.Sprintf("%s/infra/db/repositories/%s_repository.go", directory, data.ModelName), data)
	processTemplate("service_interface.tmpl", fmt.Sprintf("%s/domain/interfaces/services/%s_service_interface.go", directory, data.ModelName), data)
	processTemplate("service.tmpl", fmt.Sprintf("%s/domain/services/%s_service.go", directory, data.ModelName), data)
	processTemplate("controller_interface.tmpl", fmt.Sprintf("%s/domain/interfaces/controllers/%s_controller_interface.go", directory, data.ModelName), data)
	processTemplate("controller.tmpl", fmt.Sprintf("%s/application/controllers/%s_controller.go", directory, data.ModelName), data)
	processTemplate("route.tmpl", fmt.Sprintf("%s/application/routes/%s_routes.go", directory, data.ModelName), data)
}

func createData(modelName string, directory string) Data {
	return Data{
		PascalCasedModelName:  validations.PascalizeModelName(modelName),
		CamelCasedModelName:   validations.CamelizeModelName(modelName),
		SnakeCasedModelName:   modelName,
		ModelNameAbbreviation: modelName[0:1],
		ModelName:             modelName,
		UpdateProperties:      getUpdateProperties(modelName, directory),
		CreateProperties:      getCreateProperties(modelName, directory),
		ReadProperties:        getReadProperties(modelName, directory),
	}
}

func getUpdateProperties(modelName string, directory string) []string {

	content, err := ioutil.ReadFile(fmt.Sprintf("%s/domain/DTO/%s/update.go", directory, modelName))

	if err != nil {
		log.Fatalf(err.Error())
	}

	return getProperties(content)
}

func getReadProperties(modelName string, directory string) []string {
	content, err := ioutil.ReadFile(fmt.Sprintf("%s/domain/DTO/%s/read.go", directory, modelName))

	if err != nil {
		log.Fatalf(err.Error())
	}

	return getProperties(content)
}

func getCreateProperties(modelName string, directory string) []string {
	content, err := ioutil.ReadFile(fmt.Sprintf("%s/domain/DTO/%s/create.go", directory, modelName))

	if err != nil {
		log.Fatalf(err.Error())
	}

	return getProperties(content)
}

func getProperties(content []byte) []string {
	newContent := removeBracketsFromContent(content)
	return getPropertiesFromLines(getLines(newContent))
}

func getPropertiesFromLines(lines []string) []string {
	var properties []string
	for _, line := range lines {
		word := extractPropertyFromLine(line)
		properties = getOnlyNonEmptyWords(word, properties)
	}
	return properties
}

func getOnlyNonEmptyWords(word string, properties []string) []string {
	if word != "" {
		properties = append(properties, word)
	}
	return properties
}

func extractPropertyFromLine(line string) string {
	re := regexp.MustCompile("\t[^\\s]+")

	match := re.FindStringSubmatch(line)
	if match != nil {
		return strings.Replace(match[0], "\t", "", -1)
	}
	return ""
}

func getLines(newContent string) []string {
	lines := strings.Split(newContent, "\n")
	return lines
}

func removeBracketsFromContent(content []byte) string {
	newContent := strings.Split(string(content), "{")[1]
	newContent = strings.Split(newContent, "}")[0]
	return newContent
}

func processTemplate(fileName string, outputFile string, data Data) {
	tmpl, ok := templates[fileName]
	if !ok {
		fmt.Println(ok)
	}

	var processed bytes.Buffer
	err := tmpl.ExecuteTemplate(&processed, fileName, data)
	if err != nil {
		log.Fatalf("Unable to parse data into template: %v\n", err)
	}

	outputPath := outputFile
	fmt.Println("Writing file: ", outputPath)
	f, _ := os.Create(outputPath)
	w := bufio.NewWriter(f)
	w.WriteString(processed.String())
	w.Flush()
}
