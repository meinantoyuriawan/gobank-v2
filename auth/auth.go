package auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/meinantoyuriawan/gobank-v2/config"
	"github.com/meinantoyuriawan/gobank-v2/helper"
)

func GenerateJWT(w http.ResponseWriter, r *http.Request, username string) (string, error) {
	// generate jwt token
	expTime := time.Now().Add(time.Minute * 2)
	claims := &config.JWTClaim{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-jwt-auth",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	// declare algorithm for sign in
	tokenAlg := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// signed token
	token, err := tokenAlg.SignedString(config.JWT_KEY)

	return token, err

}

func ValidateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				response := map[string]string{"message": "Unauthorized"}
				helper.ResponseJSON(w, http.StatusBadRequest, response)
				return
			}
		}

		// get token Value
		tokenString := c.Value

		claims := &config.JWTClaim{}

		// parsing token jwt
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return config.JWT_KEY, nil
		})

		// error validation
		if err != nil {
			v, _ := err.(*jwt.ValidationError)
			switch v.Errors {
			case jwt.ValidationErrorSignatureInvalid:
				// token invalid
				response := map[string]string{"message": "Unauthorized"}
				helper.ResponseJSON(w, http.StatusUnauthorized, response)
				return
			case jwt.ValidationErrorExpired:
				// token expired
				response := map[string]string{"message": "Unauthorized, Token Expired"}
				helper.ResponseJSON(w, http.StatusUnauthorized, response)
				return
			default:
				response := map[string]string{"message": "Unauthorized"}
				helper.ResponseJSON(w, http.StatusUnauthorized, response)
				return
			}
		}

		// isvalid
		if !token.Valid {
			response := map[string]string{"message": "Unauthorized, Token Expired"}
			helper.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}

		next.ServeHTTP(w, r)
	})
}
