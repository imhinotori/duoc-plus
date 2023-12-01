package common

type Grades struct {
	Name     string    `json:"name"`
	Code     string    `json:"code"`
	Subjects []Subject `json:"subjects"`
}

type Subject struct {
	Code     string    `json:"code"`
	Name     string    `json:"name"`
	Average  float64   `json:"average"`
	Partials []float64 `json:"partials"`
	Exams    []float64 `json:"exams"`
}
