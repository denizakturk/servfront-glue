package main

import (
	"os"
	"servfront-glue/codegen"
)

func main() {
	templatePath := os.Args[1]
	outputPath := os.Args[2]
	schemaFile := os.Args[3]
	codegen.GenerateFromTemplate(templatePath, outputPath, schemaFile)
}
