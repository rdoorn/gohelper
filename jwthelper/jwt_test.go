package jwthelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	invalidToken = "eyJrbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJkbiI6InRlc3RlciIsImV4IjoxNTUyNjcwMDM3LCJpZCI6InRlc3RJRCIsInRrIjpbInBsYXllciJdfQ.V_jHvEgCYnL18B4t0HZbtC76cSjLAjBRMvGAwf3fbRESWb60PbnSagclNc-5eKM20mS6PIOh4QOs3IBxkn0tQQ"
	expiredToken = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJBcnJheVN0cmluZyI6WyJhcnJheXN0cm9uZyJdLCJJbnQiOjEsIlN0cmluZyI6InN0cmluZyIsImV4IjoxNTYzNTMyODYwLCJpZCI6InVzZXIxIiwibm9uY2UiOjAsInRrIjpbImFjIiwiZGMiLCJhZCJdfQ.uFT1r-1fKokvIhaI-Mho1jEA34HhxXRQv47mULOqjUgneC2ecNmmljrcdPTpIuaIwDud08lib2GgYzlhMaJiVQ"
)

type Token struct {
	*Credentials
	String      string
	Int         int
	ArrayString []string
}

func TestJWT(t *testing.T) {
	JWTTokenSigningKey = []byte("FT9#b=CLWjNE37iCJaWWzZRqMBhnRk&yhzAJ7FfM)d?UzyGoZW2RCwBG2w;o@bMEPfoq(VGTW8ory#KABxJjURNa3@dA?gnzWQmoWZTz3.A}DEfJVMqE?DFKtK87Yo2r")

	token := &Token{
		Credentials: &Credentials{
			Id:     "user1",
			Tokens: []string{"ac", "dc", "ad"},
		},
		String:      "string",
		Int:         1,
		ArrayString: []string{"arraystrong"},
	}

	signedToken, err := Sign(token)
	assert.Nil(t, err)

	//log.Printf("signed: %s", signedToken)
	token2 := &Token{}
	err = Validate(signedToken, token2)
	assert.Nil(t, err)

	assert.Equal(t, token.Credentials.Id, token2.Credentials.Id)
	assert.Equal(t, token.String, token2.String)

}

func TestJWTExpired(t *testing.T) {
	JWTTokenSigningKey = []byte("FT9#b=CLWjNE37iCJaWWzZRqMBhnRk&yhzAJ7FfM)d?UzyGoZW2RCwBG2w;o@bMEPfoq(VGTW8ory#KABxJjURNa3@dA?gnzWQmoWZTz3.A}DEfJVMqE?DFKtK87Yo2r")
	token := &Token{}
	err := Validate(expiredToken, token)
	assert.Equal(t, ErrTokenExpired, err)
}

func TestJWTInvalid(t *testing.T) {
	JWTTokenSigningKey = []byte("FT9#b=CLWjNE37iCJaWWzZRqMBhnRk&yhzAJ7FfM)d?UzyGoZW2RCwBG2w;o@bMEPfoq(VGTW8ory#KABxJjURNa3@dA?gnzWQmoWZTz3.A}DEfJVMqE?DFKtK87Yo2r")
	token := &Token{}
	err := Validate(invalidToken, token)
	assert.Equal(t, ErrTokenInvalidData, err)
}
