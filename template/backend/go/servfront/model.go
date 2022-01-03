package servfront
const MODEL_TEMPLATE string = "package model\n{{$modelName := .Name}}\n{{$modelFields := .Fields}}\n{{$uuidImport := false}}\n{{range $field := $modelFields }}{{if and (not $uuidImport) (eq $field.Type \"Uuid\")}}\nimport \"github.com/denizakturk/types/uuid/mysql\"{{$uuidImport = true}}\n{{end}}{{end}}\n\ntype {{$modelName}} struct {\n    {{range $field := $modelFields }}\n    {{$field.Name}} {{if eq $field.Type \"Uuid\"}}mysql.MyUuid{{else}}{{$field.Type}}{{end}} `json:\"{{$field.Name | SnakeCase}}\"`\n    {{end}}\n}"