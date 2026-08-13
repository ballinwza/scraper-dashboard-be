package config

type contextKey string

const (
	USERNAME_KEY             contextKey = "username"
	REFRESH_HASH_TOKEN_KEY   contextKey = "refresh_hash_token"
	COOKIE_REFRESH_TOKEN_KEY string     = "refresh_token"
)
