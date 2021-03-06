package chatkit

import (
	"fmt"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestNewChatkitSUTokenNoSub(t *testing.T) {
	token, expiry, err := NewChatkitSUToken("appID", "keyID", "keySecret", nil, time.Hour)
	assert.NoError(t, err, "expect no error")

	assert.False(t, time.Now().After(expiry), "expiry should be after now")

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing alg is HMAC-SHA256:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// return the key to be parsed with
		return []byte("keySecret"), nil
	})

	assert.NoError(t, err, "expect no error when parsing the token")
	assert.True(t, parsedToken.Valid, "token produced was invalid")

	claimMap, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fail()
	}

	_, present := claimMap["sub"]
	assert.False(t, present, "sub claim should not be present")
}

func TestNewChatkitSUTokenWithSub(t *testing.T) {
	sub := "jane"
	token, expiry, err := NewChatkitSUToken("appID", "keyID", "keySecret", &sub, time.Hour)
	assert.NoError(t, err, "expect no error")

	assert.False(t, time.Now().After(expiry), "expiry should be after now")

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing alg is HMAC-SHA256:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// return the key to be parsed with
		return []byte("keySecret"), nil
	})

	assert.NoError(t, err, "expect no error when parsing the token")
	assert.True(t, parsedToken.Valid, "token produced was invalid")

	claimMap, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fail()
	}

	assert.Equal(t, sub, claimMap["sub"], "token did not contain a sub claim with supplied user name")
}

func TestNewChatkitUserToken(t *testing.T) {
	token, expiry, err := NewChatkitUserToken("appID", "keyID", "keySecret", "bob", time.Hour)
	assert.NoError(t, err, "expect no error")

	assert.False(t, time.Now().After(expiry), "expiry should be after now")

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing alg is HMAC-SHA256:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// return the key to be parsed with
		return []byte("keySecret"), nil
	})

	assert.NoError(t, err, "expect no error when parsing the token")
	assert.True(t, parsedToken.Valid, "token produced was invalid")
}

func TestTokenManagerGetTokenNew(t *testing.T) {
	tokenManager := newTokenManager("testApp", "keyID", "keySecret")
	token, err := tokenManager.getToken()
	assert.NoError(t, err, "expect no error")

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing alg is HMAC-SHA256:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// return the key to be parsed with
		return []byte("keySecret"), nil
	})

	assert.NoError(t, err, "expect no error when parsing the token")
	assert.True(t, parsedToken.Valid, "token produced was invalid")
}

func TestTokenManagerGetTokenNotExpired(t *testing.T) {
	tokenManager := newTokenManager("testApp", "keyID", "keySecret")
	firstToken, err := tokenManager.getToken()
	assert.NoError(t, err, "expect no error")

	secondToken, err := tokenManager.getToken()
	assert.NoError(t, err, "expect no error")

	assert.Equal(t, firstToken, secondToken, "don't expect tokens to be regenerated if not expired")
}

func newMockTokenManager() tokenManager {
	return &mockTokenManager{}
}

type mockTokenManager struct{}

func (mtm *mockTokenManager) getToken() (string, error) { return "", nil }
