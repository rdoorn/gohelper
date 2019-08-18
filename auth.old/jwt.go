package auth

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	jwt "github.com/dgrijalva/jwt-go"
)

// JWTUserCredentials is the token given to the user
type JWTUserCredentials struct {
	UserID      string   `json:"id"`    // user ID matching the user
	DisplayName string   `json:"dn"`    // display name of the user
	Tokens      []string `json:"tk"`    // authorization tokens of the user
	Nonce       float64  `json:"nonce"` // nonce
}

var (
	// JWTTokenSigningKey is key used to sign jtw tokens
	//JWTTokenSigningKey = rndKey()
	JWTTokenSigningKey = []byte("FT9#b=CLWjNE37iCJaWWzZRqMBhnRk&yhzAJ7FfM)d?UzyGoZW2RCwBG2w;o@bMEPfoq(VGTW8ory#KABxJjURNa3@dA?gnzWQmoWZTz3.A}DEfJVMqE?DFKtK87Yo2r")
	// JWTTokenDuration is how long the jwt token is valid
	JWTTokenDuration = 1 * time.Hour
)

// SignToken turns users credentials in to a JWT token string
func SignToken(credentials *JWTUserCredentials) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id": credentials.UserID,
		"dn": credentials.DisplayName,
		"tk": credentials.Tokens,
		"ex": time.Now().Add(JWTTokenDuration).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(JWTTokenSigningKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateToken validates the Authorization header provided by the client, and returns the JWT token
func validateToken(tokenStr string) (*JWTUserCredentials, error) {
	if tokenStr == "" {
		return nil, fmt.Errorf("Invalid token supplied")
	}

	// parse tokenStr and get its token details
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return JWTTokenSigningKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//log.Printf("Token Claims: %+v\n", claims)
		if time.Now().Unix() > int64(claims["ex"].(float64)) {
			return nil, fmt.Errorf("Token expired")
		}

		userCredentials := &JWTUserCredentials{
			UserID:      claims["id"].(string),
			DisplayName: claims["dn"].(string),
		}

		for _, token := range claims["tk"].([]interface{}) {
			userCredentials.Tokens = append(userCredentials.Tokens, token.(string))
		}

		log.Printf("Token Claims: %+v\n", claims)

		return userCredentials, nil
	}

	return nil, fmt.Errorf("Invalid token claims")
}

// ValidateToken validate an api user, and get its credentials
func ValidateToken(req events.APIGatewayProxyRequest) (*JWTUserCredentials, error) {
	authHeader, ok := req.Headers["Authorization"]
	if !ok {
		return nil, ErrOauthUnauthorizedClient
	}
	auth := strings.Split(authHeader, " ")
	switch auth[0] {
	case "Bearer":
		return validateToken(auth[1])
	default:
		return nil, ErrOauthUnauthorizedClient
	}
}

// IsAuthorized returns true if a users JWT credentials match any tokens
func IsAuthorized(user *JWTUserCredentials, requiredTokens ...string) bool {
	// if we require no token, then your authorized
	if len(requiredTokens) == 0 {
		return true
	}

	// 'any' token matches all
	if found("any", user.Tokens) {
		return true
	}

	// check individual tokens
	for _, requiredToken := range requiredTokens {
		if found(requiredToken, user.Tokens) {
			return true
		}
	}
	return false
}

// found returns true if a string is found in an array
func found(s string, a []string) bool {
	return a[sort.SearchStrings(a, s)] == s
}
