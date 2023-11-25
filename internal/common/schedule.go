package common

type Course struct {
	SubjectName string `json:"subject_name"`
	SubjectCode string `json:"subject_code"`
	Instructor  string `json:"instructor"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Classroom   string `json:"classroom"`
}

type Week struct {
	Monday    []Course `json:"monday"`
	Tuesday   []Course `json:"tuesday"`
	Wednesday []Course `json:"wednesday"`
	Thursday  []Course `json:"thursday"`
	Friday    []Course `json:"friday"`
	Saturday  []Course `json:"saturday"`
	Sunday    []Course `json:"sunday"`
}

type CareerSchedule struct {
	CareerName string `json:"career_name"`
	Schedule   Week   `json:"schedule"`
}
