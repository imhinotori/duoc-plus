package common

type AuthenticationResponse struct {
	AccessToken           string `json:"access_token"`
	AccessTokenExpiresIn  int    `json:"access_token_expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_expires_in"`
}
