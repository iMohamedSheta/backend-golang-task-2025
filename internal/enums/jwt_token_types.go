package enums

type JwtTokenType string

const (
	AccessToken  JwtTokenType = "access_token"
	RefreshToken JwtTokenType = "refresh_token"
)
