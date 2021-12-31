package model

type ControllerMetod struct {
	Name              string `json:"name"`
	RequestModelName  string `json:"request_model_name"`
	ResponseModelName string `json:"response_model_name"`
}

type Model struct {
	Name  string       `json:"name"`
	Field []ModelField `json:"field"`
}

type ModelField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type AutogenInit struct {
	PackageName    string            `json:"package_name"`
	ControllerName string            `json:"controller_name"`
	ControllerPath string            `json:"controller_path"`
	ModelPath      string            `json:"model_path"`
	MiddlewarePath string            `json:"middleware_path"`
	FrontendPath   string            `json:"frontend_path"`
	Method         []ControllerMetod `json:"method"`
	Model          []Model           `json:"model"`
}
