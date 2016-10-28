package task

type Task struct {
	Id string `json: "id"`
	Type string `json: "type"`
	Stauts string `json: "status"`
}
