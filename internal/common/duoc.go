package common

//This contains all the original Duoc structures

//ATTENDANCE

type DuocSubjectAttendanceDetail struct {
	Date       string `json:"fechaLarga"`
	Attendance string `json:"asistencia"`
}

type DuocSubjectAttendanceHeader struct {
	SubjectName     string `json:"nomAsignatura"`
	SubjectCode     string `json:"codAsignatura"`
	ClassesHeld     string `json:"clasesRealizadas"`
	AssistedClasses string `json:"clasesAsistente"`
	Percentage      string `json:"porcentaje"` //Why...?
}

type DuocSubjectAttendance struct {
	Header  DuocSubjectAttendanceHeader   `json:"cabecera"`
	Details []DuocSubjectAttendanceDetail `json:"detalle"`
}

type DuocAttendance struct {
	DegreeName        string                  `json:"nomCarrera"`
	DegreeCode        string                  `json:"codCarrera"`
	SubjectAttendance []DuocSubjectAttendance `json:"asistenciaAsignaturas"`
}

// GRADES

type DuocGradesCourses struct {
	DegreeName string        `json:"nomCarrera"`
	DegreeCode string        `json:"codCarrera"`
	Subjects   []DuocSubject `json:"asignaturas"`
}

type DuocSubject struct {
	SubjectCode   string      `json:"codAsignatura"`
	Name          string      `json:"nombre"`
	Average       string      `json:"promedio"`
	PartialGrades []DuocGrade `json:"parciales"`
	ExamsGrades   []DuocGrade `json:"examenes"`
}

type DuocGrade struct {
	Text  string `json:"texto"`
	Grade string `json:"nota"`
}

// SCHEDULE

type DuocDay struct {
	Day     string       `json:"dia"`
	Courses []DuocCourse `json:"ramos"`
}

type DuocCourse struct {
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

type DuocSchedule struct {
	CodeCareer string    `json:"codCarrera"`
	NameCareer string    `json:"nomCarrera"`
	Days       []DuocDay `json:"dias"`
}
