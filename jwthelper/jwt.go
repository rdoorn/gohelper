package jwthelper

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWTUserCredentials is the token given to the user
type Credentials struct {
	Id     string   `json:"id"`    // user ID matching the user
	Tokens []string `json:"tk"`    // authorization tokens of the user
	Nonce  float64  `json:"nonce"` // nonce
}

var (
	// JWTTokenSigningKey is key used to sign jtw tokens
	JWTTokenSigningKey = rndKey(32)
	// JWTTokenSigningKey set static for testing, but should be random..
	// JWTTokenDuration is how long the jwt token is valid
	JWTTokenDuration = 1 * time.Hour
)

// Errors returned
var (
	ErrTokenExpired              = errors.New("Token has expired")
	ErrTokenInvalidData          = errors.New("Token contains invalid data")
	ErrTokenInvalidSigningMethod = errors.New("Token has an invalid signing method")
	ErrTokenConversionError      = errors.New("Token conversion error")
)

func rndKey(i int) []byte {
	token := make([]byte, i)
	rand.Read(token)
	return token
}

// Sign turns users credentials in to a JWT token string
func Sign(i interface{}) (string, error) {

	// convert interface to jwt.MapClaims
	j, _ := json.Marshal(i)
	var claims jwt.MapClaims
	if err := json.Unmarshal(j, &claims); err != nil {
		return "", ErrTokenConversionError
	}

	// add expiry
	claims["ex"] = time.Now().Add(JWTTokenDuration).Unix()

	// Sign and get the complete encoded token as a string using the secret
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(JWTTokenSigningKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Validate validates the Authorization header provided by the client, and returns the JWT token details
func Validate(tokenStr string, i interface{}) error {
	if tokenStr == "" {
		return ErrTokenInvalidData
	}

	// parse tokenStr and get its token details
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalidSigningMethod
		}

		return JWTTokenSigningKey, nil
	})
	if err != nil {
		return ErrTokenInvalidData
	}

	// validate clain
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if time.Now().Unix() > int64(claims["ex"].(float64)) {
			return ErrTokenExpired
		}

		// convert claim to the requested interface
		j, _ := json.Marshal(token.Claims)
		if err := json.Unmarshal(j, i); err != nil {
			return ErrTokenConversionError
		}

		return nil
	}

	return ErrTokenInvalidData
}
