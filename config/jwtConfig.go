package config

type JwtConfig struct {
	Issuer   string
	Audience string // no password set
}

func NewJwtConfig(issuer string, audience string) *JwtConfig {
	return &JwtConfig{Issuer: issuer, Audience: audience}
}
