package auth

import (
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

var (
	expiredToken = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJkbiI6InRlc3RlciIsImV4IjoxNTUyNjcwMDM3LCJpZCI6InRlc3RJRCIsInRrIjpbInBsYXllciJdfQ.V_jHvEgCYnL18B4t0HZbtC76cSjLAjBRMvGAwf3fbRESWb60PbnSagclNc-5eKM20mS6PIOh4QOs3IBxkn0tQQ"
)

func TestJWT(t *testing.T) {
	credentials := &JWTUserCredentials{
		UserID:      "testID",
		DisplayName: "tester",
		Tokens:      []string{"player"},
	}

	tokenString, err := SignToken(credentials)
	assert.Nil(t, err)

	req := events.APIGatewayProxyRequest{
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", tokenString),
		},
	}
	credentials2, err := ValidateToken(req)
	assert.Nil(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, credentials.UserID, credentials2.UserID)

	assert.Equal(t, true, IsAuthorized(credentials2, "player"))
	assert.Equal(t, false, IsAuthorized(credentials2, "gameover"))
}

func TestJWTExpired(t *testing.T) {
	_, err := validateToken(expiredToken)
	assert.NotNil(t, err)
}
