package third

type Param struct {
	Name string
	Type string
}

type Api struct {
	Name        string
	Description string
	Path        string
	Request     string
	Method      string
	See         string
	FuncName    string
	GetParams   []Param
}

type ApiGroup struct {
	Name    string
	Apis    []Api
	Package string
}
