package common

type Career struct {
	School     string `json:"school"`
	CareerName string `json:"career_name"`
	CareerCode string `json:"career_code"`
	Campus     string `json:"campus"`
}

type User struct {
	FullName    string   `json:"full_name"`
	Rut         string   `json:"rut"`
	Avatar      string   `json:"avatar"`
	Careers     []Career `json:"careers"`
	Email       string   `json:"email,omitempty"`
	Username    string   `json:"username,omitempty"`
	StudentCode string   `json:"codAlumno,omitempty"` // It's probably an int, but well.
	StudentId   int      `json:"idAlumno,omitempty"`  // Why two ids (?) I don't know.
	//
	AccessToken           string `json:"access_token,omitempty"`
	AccessTokenExpiresIn  int    `json:"access_token_expires_in,omitempty"`
	RefreshToken          string `json:"refresh_token,omitempty"`
	RefreshTokenExpiresIn int    `json:"refresh_expires_in,omitempty"`
}
