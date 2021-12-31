package servfront

type ControllerMethod struct {
	Name              string `json:"name"`
	RequestModelName  string `json:"request_model_name"`
	ResponseModelName string `json:"response_model_name"`
}

type ControllerTmplParams struct {
	ControllerName string             `json:"controller_name"`
	Method         []ControllerMethod `json:"method"`
}

type ModelTmplParamFields struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type ModelTmplParams struct {
	Name   string                 `json:"name"`
	Fields []ModelTmplParamFields `json:"fields"`
}
