package generate

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

type Data struct {
	ModelName        string
	UpdateProperties []string
	CreateProperties []string
	ReadProperties   []string
}

func missing(a, b []string) string {
	ma := make(map[string]bool, len(a))
	for _, ka := range a {
		ma[ka] = true
	}
	for _, kb := range b {
		if !ma[kb] {
			return kb
		}
	}
	return ""
}

func Generate(modelName string) {
	directory, err := os.Getwd()
	if err != nil {
		log.Fatal("%s", err.Error())
	}
	files, _ := ioutil.ReadDir(directory + "/domain/DTO/" + modelName)

	necessaryDtos := []string{"create", "read", "update"}

	if len(files) == 0 {
		log.Fatalf("No DTOs found for [%s] , exiting.", modelName)
	}

	var foundFiles []string

	if len(files) < 3 {
		for _, f := range files {
			foundFiles = append(foundFiles, strings.Split(f.Name(), ".")[0])
		}

		log.Fatalf("Missing %sDTO for [%s], exiting.", cases.Title(language.BrazilianPortuguese).String(missing(foundFiles, necessaryDtos)), modelName)
	}

	data := Data{
		ModelName:        modelName,
		UpdateProperties: getUpdateProperties(modelName, err, directory),
		CreateProperties: getCreateProperties(modelName, err, directory),
		ReadProperties:   getReadProperties(modelName, err, directory),
	}

	processTemplate("repository.tmpl", fmt.Sprintf("%s/infra/db/repositories/%s_repository.go", directory, data.ModelName), data)
	processTemplate("repository_interface.tmpl", fmt.Sprintf("%s/domain/interfaces/repositories/%s_repository_interface.go", directory, data.ModelName), data)
	processTemplate("controller_interface.tmpl", fmt.Sprintf("%s/domain/interfaces/controllers/%s_controller_interface.go", directory, data.ModelName), data)
	processTemplate("controller.tmpl", fmt.Sprintf("%s/application/controllers/%s_controller.go", directory, data.ModelName), data)
	processTemplate("service_interface.tmpl", fmt.Sprintf("%s/domain/interfaces/services/%s_service_interface.go", directory, data.ModelName), data)
	processTemplate("service.tmpl", fmt.Sprintf("%s/domain/services/%s_service.go", directory, data.ModelName), data)
}

func getUpdateProperties(modelName string, err error, directory string) []string {
	content, err := ioutil.ReadFile(fmt.Sprintf("%s/domain/DTO/%s/update.go", directory, modelName))

	if err != nil {
		log.Fatalf(err.Error())
	}

	return getProperties(content)
}

func getReadProperties(modelName string, err error, directory string) []string {
	content, err := ioutil.ReadFile(fmt.Sprintf("%s/domain/DTO/%s/read.go", directory, modelName))

	if err != nil {
		log.Fatalf(err.Error())
	}

	return getProperties(content)
}

func getCreateProperties(modelName string, err error, directory string) []string {
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

	tmpl := template.Must(template.New("").Funcs(sprig.FuncMap()).ParseFiles("tmpl/" + fileName))

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
