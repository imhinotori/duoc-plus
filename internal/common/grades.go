package common

type GradesCourses struct {
	DegreeName string    `json:"nomCarrera"`
	DegreeCode string    `json:"codCarrera"`
	Subjects   []Subject `json:"asignaturas"`
}

type Subject struct {
	SubjectCode   string  `json:"codAsignatura"`
	Name          string  `json:"nombre"`
	Average       string  `json:"promedio"`
	PartialGrades []Grade `json:"parciales"`
	ExamsGrades   []Grade `json:"examenes"`
}

type Grade struct {
	Text  string `json:"texto"`
	Grade string `json:"nota"`
}
