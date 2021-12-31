package codegen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"
)

func ParseTemplate(templatePath string) *template.Template {
	tmpl := template.New(path.Base(templatePath))
	tmpl.Funcs(GetTemplateFuncMap())
	tmpl, tmplErr := tmpl.ParseFiles(templatePath)
	if nil != tmplErr {
		fmt.Println(tmplErr)
	}
	return tmpl
}

func ExecTemplate(template *template.Template, templateParameters interface{}) []byte {
	templateOutput := &bytes.Buffer{}
	err := template.Execute(templateOutput, templateParameters)
	if nil != err {
		fmt.Println(err)
	}

	formatedTemplate, formatErr := format.Source(templateOutput.Bytes())
	var outputTemplate []byte
	if nil != formatErr {
		outputTemplate = templateOutput.Bytes()
		fmt.Println(formatErr)
	} else {
		outputTemplate = formatedTemplate
	}
	return outputTemplate
}

func OutputTemplate(outputPath string, byteCode []byte) {
	os.WriteFile(outputPath, byteCode, 0777)
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func To_snake_case(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func GetTemplateFuncMap() template.FuncMap {
	funcMap := template.FuncMap{
		"ToLower":   strings.ToLower,
		"SnakeCase": To_snake_case,
	}

	return funcMap
}

func GenerateFromTemplate(templatePath string, outputPath string, parameters interface{}) {
	tmpl := ParseTemplate(templatePath)

	generateByteCode := ExecTemplate(tmpl, parameters)

	OutputTemplate(outputPath, generateByteCode)
}

func SchemaFileToParameter(schemaFile string) interface{} {
	jsonFile, err := os.Open(schemaFile)
	if nil != err {
		panic(err)
	}
	defer jsonFile.Close()

	jsonString, _ := ioutil.ReadAll(jsonFile)
	var parameters interface{}
	json.Unmarshal(jsonString, &parameters)

	return parameters
}
