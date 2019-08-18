package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"sbp.gitlab.schubergphilis.com/GCT-Finance/vault/generic/pkg/api"
)

// Errors
var (
	// ErrOauthInvalidRequest
	// The request is missing a required parameter, includes an
	// invalid parameter value, includes a parameter more than
	// once, or is otherwise malformed.
	ErrOauthInvalidRequest = errors.New("invalid_request")

	// ErrOauthUnauthorizedClient
	// The client is not authorized to request an authorization
	// code using this method.
	ErrOauthUnauthorizedClient = errors.New("unauthorized_client")

	// ErrOauthAccessDenied
	// The resource owner or authorization server denied the
	// request.
	ErrOauthAccessDenied = errors.New("access_denied")

	// ErrOauthUnsupportedResponseType
	// The authorization server does not support obtaining an
	// authorization code using this method.
	ErrOauthUnsupportedResponseType = errors.New("unsupported_response_type")

	// ErrOauthInvalidScope
	// The requested scope is invalid, unknown, or malformed.
	ErrOauthInvalidScope = errors.New("invalid_scope")

	// ErrOauthServerError
	// The authorization server encountered an unexpected
	// condition that prevented it from fulfilling the request.
	// (This error code is needed because a 500 Internal Server
	// Error HTTP status code cannot be returned to the client
	// via an HTTP redirect.)
	ErrOauthServerError = errors.New("server_error")

	// ErrOauthTemporarilyUnavailable
	// The authorization server is currently unable to handle
	// the request due to a temporary overloading or maintenance
	// of the server.  (This error code is needed because a 503
	// Service Unavailable HTTP status code cannot be returned
	// to the client via an HTTP redirect.)
	ErrOauthTemporarilyUnavailable = errors.New("temporarily_unavailable")
)

var (
	oauthErrorCode = map[error]int{
		ErrOauthInvalidRequest:          http.StatusBadRequest,
		ErrOauthUnauthorizedClient:      http.StatusUnauthorized,
		ErrOauthAccessDenied:            http.StatusForbidden,
		ErrOauthUnsupportedResponseType: http.StatusBadRequest,
		ErrOauthInvalidScope:            http.StatusBadRequest,
		ErrOauthServerError:             http.StatusInternalServerError,
		ErrOauthTemporarilyUnavailable:  http.StatusServiceUnavailable,
	}
)

// OauthError contains the error message details for oauth
type OauthError struct {
	ErrorMsg         error  `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorURL         string `json:"error_url"`
}

// OauthErrorMsg returns a error message for the client in json format
func OauthErrorMsg(req events.APIGatewayProxyRequest, oauthErr *OauthError) (events.APIGatewayProxyResponse, error) {
	headers := api.NewHeaders()
	headers.AddCORS(req)

	errMsg := fmt.Sprintf(`{"error":"%s","error_description":"%s","error_url":"%s"}`, oauthErr.ErrorMsg.Error(), url.QueryEscape(oauthErr.ErrorDescription), oauthErr.ErrorURL)
	log.Printf("error message: [%s]", errMsg)
	return events.APIGatewayProxyResponse{
		Headers:    headers.Get(),
		StatusCode: oauthErrorCode[oauthErr.ErrorMsg],
		Body:       string(errMsg),
	}, nil
}

// OauthRedirect returns a error message for the client in json format
func OauthRedirect(req events.APIGatewayProxyRequest, oauthErr *OauthError, location string) (events.APIGatewayProxyResponse, error) {
	headers := api.NewHeaders()
	headers.AddCORS(req)
	headers.Add("Location", fmt.Sprintf(`%s?error=%s&error_description=%s&error_url=%s`, location, oauthErr.ErrorMsg.Error(), url.QueryEscape(oauthErr.ErrorDescription), oauthErr.ErrorURL))

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusFound,
		/*	Headers: map[string]string{
			"Location": fmt.Sprintf(`%s?error=%s&error_description=%s&error_url=%s`, location, oauthErr.ErrorMsg.Error(), oauthErr.ErrorDescription, oauthErr.ErrorURL),
		},*/
		Headers: headers.Get(),
	}, nil
}

// OauthClients contains the client apps registered to use oauth
type OauthClients []OauthClient

// OauthClient contains the client app registered to use oauth
type OauthClient struct {
	ClientID    string `json:"client_id"`
	SecretID    string `json:"secret_id"`
	RedirectURL string `json:"redirect_url"`
}

// GetClientByID gets the client from OauthClients based on ID
func (clients *OauthClients) GetClientByID(id string) (OauthClient, error) {
	for _, client := range *clients {
		if client.ClientID == id {
			return client, nil
		}
	}
	return OauthClient{}, errors.New("invalid client id")
}

// GrantToken checks the credentials provided, and gives back user credentials if authorized
/*func GrantToken(client OauthClient, params url.Values) (*JWTUserCredentials, error) {
	switch params.Get("grant_type") {
	case "password":
		authenticatedUser, err := AuthenticateCredentials(params.Get("username"), params.Get("password"))
		if err != nil {
			return nil, ErrOauthUnauthorizedClient
		}
		return authenticatedUser, nil
	default:
		return nil, ErrOauthInvalidRequest
	}

}
*/
