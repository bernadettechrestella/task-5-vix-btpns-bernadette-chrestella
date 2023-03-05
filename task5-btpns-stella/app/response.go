package app

type SignInResponse struct {
	AccessToken      string `json:"access_token" out:"accessToken"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token" out:"refreshToken"`
	TokenType        string `json:"token_type"`
	IDToken          string `json:"id_token"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

type ApiResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
