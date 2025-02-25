package models

type AuthConfig struct {
	JWTSecret string
}

func NewAuthConfig(jwtSecret string) *AuthConfig {
	return &AuthConfig{JWTSecret: jwtSecret}
}
