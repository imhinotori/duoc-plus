package common

type SubjectAttendanceDetail struct {
	Date       string `json:"fechaLarga"`
	Attendance string `json:"asistencia"`
}

type SubjectAttendanceHeader struct {
	SubjectName     string `json:"nomAsignatura"`
	SubjectCode     string `json:"codAsignatura"`
	ClassesHeld     string `json:"clasesRealizadas"`
	AssistedClasses string `json:"clasesAsistente"`
	Percentage      string `json:"porcentaje"` //Why...?
}

type SubjectAttendance struct {
	Header  SubjectAttendanceHeader   `json:"cabecera"`
	Details []SubjectAttendanceDetail `json:"detalle"`
}

type Attendance struct {
	DegreeName        string              `json:"nomCarrera"`
	DegreeCode        string              `json:"codCarrera"`
	SubjectAttendance []SubjectAttendance `json:"asistenciaAsignaturas"`
}
