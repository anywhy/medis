package state

type AppDefinition struct {
	Id        string            `json:"id"`
	Cmd       string            `json:"cmd"`
	Cpus      float64           `json:"cpus"`
	Mem       float64           `json:"mem"`
	Ports     []int             `json:"ports"`
	Env       map[string]string `json:"env"`
	Instances int               `json:"instances"`
}
