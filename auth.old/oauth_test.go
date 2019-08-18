package auth

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestOauthErrorMsg(t *testing.T) {
	req := events.APIGatewayProxyRequest{}
	res, err := OauthErrorMsg(req, &OauthError{ErrorMsg: ErrOauthUnauthorizedClient})
	assert.Nil(t, err)
	assert.Equal(t, `{"error":"unauthorized_client","error_description":"","error_url":""}`, res.Body)
	assert.Equal(t, 401, res.StatusCode)
}

func TestOauthRedirect(t *testing.T) {
	req := events.APIGatewayProxyRequest{}
	res, err := OauthRedirect(req, &OauthError{ErrorMsg: ErrOauthUnauthorizedClient}, "/redirect")
	assert.Nil(t, err)
	//assert.Equal(t, `{"error":"unauthorized_client","error_description":"","error_url":""}`, res.Body)
	assert.Equal(t, 302, res.StatusCode)
	assert.Equal(t, "/redirect?error=unauthorized_client&error_description=&error_url=", res.Headers["Location"])
}
