package authDto

import "github.com/dgrijalva/jwt-go"

type (
	JwtClaim struct {
		jwt.StandardClaims
		Username string `json:"username"`
	}
)
