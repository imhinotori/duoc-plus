package common

type Career struct {
	School     string `json:"school"`
	CareerName string `json:"career_name"`
	CareerCode string `json:"career_code"`
	Campus     string `json:"campus"`
}

type User struct {
	FullName string   `json:"full_name"`
	Rut      string   `json:"rut"`
	Avatar   string   `json:"avatar"`
	Careers  []Career `json:"careers"`
}
