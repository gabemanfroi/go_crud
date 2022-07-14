package generate

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig"
	"go/format"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Data struct {
	ModelName       string
	ModelProperties []string
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
	files, _ := ioutil.ReadDir(directory + "/DTO/" + modelName)

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

	content, err := ioutil.ReadFile(fmt.Sprintf("%s/DTO/%s/update.go", directory, modelName))

	if err != nil {
		log.Fatalf(err.Error())
	}

	newContent := strings.Split(string(content), "{")[1]
	newContent = strings.Split(newContent, "}")[0]

	lines := strings.Split(newContent, "\n")

	var properties []string

	for _, line := range lines {
		properties = append(properties, strings.Replace(strings.Split(line, "*")[0], "\t", "", -1))
	}

	data := Data{
		ModelName:       modelName,
		ModelProperties: properties,
	}
	processTemplate("repository.tmpl", fmt.Sprintf("./repositories/%s_repository.go", data.ModelName), data)
	/*processTemplate("repository_interface.tmpl", fmt.Sprintf("./repositories/%s_repository_interface.go", data.ModelName), data)
	processTemplate("controller_interface.tmpl", fmt.Sprintf("./controllers/%s_controller_interface.go", data.ModelName), data)
	processTemplate("controller.tmpl", fmt.Sprintf("./controllers/%s_controller.go", data.ModelName), data)
	processTemplate("service_interface.tmpl", fmt.Sprintf("./services/%s_service_interface.go", data.ModelName), data)
	processTemplate("service.tmpl", fmt.Sprintf("./services/%s_service.go", data.ModelName), data)*/
}

func processTemplate(fileName string, outputFile string, data Data) {

	tmpl := template.Must(template.New("").Funcs(sprig.FuncMap()).ParseFiles("tmpl/" + fileName))

	var processed bytes.Buffer
	err := tmpl.ExecuteTemplate(&processed, fileName, data)
	if err != nil {
		log.Fatalf("Unable to parse data into template: %v\n", err)
	}

	formatted, err := format.Source(processed.Bytes())
	if err != nil {
		log.Fatalf("Could not format processed template: %v\n", err)
	}
	outputPath := outputFile
	fmt.Println("Writing file: ", outputPath)
	f, _ := os.Create(outputPath)
	w := bufio.NewWriter(f)
	w.WriteString(string(formatted))
	w.Flush()
}
