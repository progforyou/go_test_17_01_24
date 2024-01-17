package data

type AgeRequest struct {
	Name string `json:"name"`
	Age  uint64 `json:"age"`
}

type GenderRequest struct {
	Name     string `json:"name"`
	Gender   Gender `json:"gender"`
	Probably string `json:"probably"`
}

type NationalizeRequest struct {
	Name    string        `json:"name"`
	Country []Nationalize `json:"country"`
}
