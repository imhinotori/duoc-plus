package common

import "time"

type AttendanceDetail struct {
	Date       time.Time `json:"date"`
	Attendance int       `json:"attendance"`
}

type SubjectAttendance struct {
	Name            string             `json:"name"`
	Code            string             `json:"code"`
	ClassesHeld     int                `json:"classes_held"`
	AssistedClasses int                `json:"assisted_classes"`
	Percentage      float64            `json:"percentage"`
	Details         []AttendanceDetail `json:"details"`
}

type Attendance struct {
	DegreeName        string              `json:"degree_name"`
	DegreeCode        string              `json:"degree_code"`
	SubjectAttendance []SubjectAttendance `json:"attendance"`
}
