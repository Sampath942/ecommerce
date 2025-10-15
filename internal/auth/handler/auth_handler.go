package handler

type AuthHandler interface {
	GenerateJWTToken(uid int, isAdmin bool) (string, error)
	ValidateJWTToken(tokenString string) (map[string]any, error)
}

type authHandler struct{}

func NewAuthHandler() AuthHandler {
	return &authHandler{}
}

func (*authHandler) GenerateJWTToken(uid int, isAdmin bool) (string, error) {
	return GenerateJWTToken(uid, isAdmin)
}

func (*authHandler) ValidateJWTToken(tokenString string) (map[string]any, error) {
	return ValidateJWTToken(tokenString)
}
