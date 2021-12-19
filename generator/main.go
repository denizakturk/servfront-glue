package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"html/template"
	"os"
	"path"
	"regexp"
	"strings"
)

func main() {
	fmt.Println("asdasddsa")
	tmeplatePath := os.Args[1]
	outputPath := os.Args[2]
	schemaFile := os.Args[3]
	template := parseTemplate(tmeplatePath)
	generateByteCode := bytes.Buffer{}
	var parameters interface

	jsonFile, err := os.Open(schemaFile)
	defer jsonFile.Close()

	json.Unmarshal(jsonFile, parameters)

	template.ExecuteTemplate(template, generateByteCode, parameters)

	outputTemplate()
}

func parseTemplate(templatePath string) *template.Template {
	tmpl := template.New(path.Base(templatePath))
	tmpl.Funcs(getTemplateFuncMap())
	tmpl, tmplErr := tmpl.ParseFiles(templatePath)
	if nil != tmplErr {
		fmt.Println(tmplErr)
	}
	return tmpl
}

func execTemplate(template *template.Template, templateOutput *bytes.Buffer, templateParameters interface{}) {

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
	templateOutput.Reset()
	templateOutput.Write(outputTemplate)
}

func outputTemplate(outputPath string, byteCode []byte) {
	os.WriteFile(outputPath, byteCode, 0777)
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func to_snake_case(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func getTemplateFuncMap() template.FuncMap {
	funcMap := template.FuncMap{
		"ToLower":   strings.ToLower,
		"SnakeCase": to_snake_case,
	}

	return funcMap
}
