package middlewares

import (
	"fmt"
	"service-user/model/dto/authDto"
	"service-user/model/dto/json"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	applicationName  = "incubation-golang"
	jwtSigningMethod = jwt.SigningMethodHS256
	jwtSignatureKey  = []byte("incubation-golang")
)

func GenerateTokenJwt(username string, expiredAt int64) (string, error) {
	loginExpDuration := time.Duration(expiredAt) * time.Minute
	myExpiresAt := time.Now().Add(loginExpDuration).Unix()
	claims := authDto.JwtClaim{
		StandardClaims: jwt.StandardClaims{
			Issuer:    applicationName,
			ExpiresAt: myExpiresAt,
		},
		Username: username,
	}

	token := jwt.NewWithClaims(
		jwtSigningMethod,
		claims,
	)

	signedToken, err := token.SignedString(jwtSignatureKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func BasicAuth(c *gin.Context) {
	user, password, ok := c.Request.BasicAuth()
	if !ok {
		fmt.Println("MASUK ERROR BASIC AUTH")
		json.NewAbortUnauthorized(c, "invalid token", "01", "01")
		return
	}

	// TODO use env
	if user != "nugroho" || password != "secretpass" {
		json.NewAbortUnauthorized(c, "unauthorized", "01", "01")
		return
	}

	c.Next()
}

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			json.NewAbortUnauthorized(c, "invalid token", "01", "01")
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", -1)
		claims := &authDto.JwtClaim{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSignatureKey, nil
		})
		if err != nil {
			fmt.Println("ERROR JWTAUTH TOKEN")
			json.NewAbortUnauthorized(c, "invalid token", "01", "01")
			return
		}
		if !token.Valid {
			fmt.Println("ERROR JWTAUTH FORBIDDEN")
			json.NewAbortForbidden(c, "access forbidden", "01", "01")
			return
		}
		c.Next()
	}
}
