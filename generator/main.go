package main

import (
	"os"

	"github.com/denizakturk/servfront-glue/codegen"
)

func main() {
	templatePath := os.Args[1]
	outputPath := os.Args[2]
	parameters := codegen.SchemaFileToParameter(os.Args[3])
	codegen.GenerateFromFileTemplate(templatePath, outputPath, parameters)
}
