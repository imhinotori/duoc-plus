package common

type Day struct {
	Day     string   `json:"dia"`
	Courses []Course `json:"ramos"`
}

type Course struct {
	Day        string `json:"dia"`
	Name       string `json:"nombre"`
	StartTime  string `json:"horaInicio"`
	EndTime    string `json:"horaFin"`
	Code       string `json:"sigla"`
	Instructor string `json:"profesor"`
	Classroom  string `json:"sala"`
	Campus     string `json:"sede"`
	PlanCode   string `json:"codPlan"`
	PlanName   string `json:"nombrePlan"`
	Section    string `json:"seccion"`
}

type Schedule struct {
	CodeCareer string `json:"codCarrera"`
	NameCareer string `json:"nomCarrera"`
	Days       []Day  `json:"dias"`
}
