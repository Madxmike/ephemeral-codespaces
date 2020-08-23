package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

const projectID = "test_project"

var validClaim = &jwt.Claims{
	Issuer:   fmt.Sprintf(secureTokenEndpoint, projectID),
	Subject:  projectID,
	ID:       projectID,
	Audience: jwt.Audience{projectID},
	IssuedAt: jwt.NewNumericDate(time.Now()),
	Expiry:   jwt.NewNumericDate(time.Now()),
}

var wrongIssuerClaim = &jwt.Claims{
	Issuer:   fmt.Sprintf(secureTokenEndpoint, "wrong"),
	Subject:  projectID,
	ID:       projectID,
	Audience: jwt.Audience{projectID},
	IssuedAt: jwt.NewNumericDate(time.Now()),
	Expiry:   jwt.NewNumericDate(time.Now()),
}

var wrongSubjectClaim = &jwt.Claims{
	Issuer:   fmt.Sprintf(secureTokenEndpoint, projectID),
	Subject:  "wrong",
	ID:       projectID,
	Audience: jwt.Audience{projectID},
	IssuedAt: jwt.NewNumericDate(time.Now()),
	Expiry:   jwt.NewNumericDate(time.Now()),
}

var wrongAudienceClaim = &jwt.Claims{
	Issuer:   fmt.Sprintf(secureTokenEndpoint, projectID),
	Subject:  projectID,
	ID:       projectID,
	Audience: jwt.Audience{"wrong"},
	IssuedAt: jwt.NewNumericDate(time.Now()),
	Expiry:   jwt.NewNumericDate(time.Now()),
}

var futureIssuedClaim = &jwt.Claims{
	Issuer:   fmt.Sprintf(secureTokenEndpoint, projectID),
	Subject:  projectID,
	ID:       projectID,
	Audience: jwt.Audience{projectID},
	IssuedAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
	Expiry:   jwt.NewNumericDate(time.Now()),
}

var expiredClaim = &jwt.Claims{
	Issuer:   fmt.Sprintf(secureTokenEndpoint, projectID),
	Subject:  projectID,
	ID:       projectID,
	Audience: jwt.Audience{projectID},
	IssuedAt: jwt.NewNumericDate(time.Now()),
	Expiry:   jwt.NewNumericDate(time.Time{}.Add(1 * time.Hour)),
}

func generateJWT(claim *jwt.Claims, privateKey *rsa.PrivateKey, t *testing.T) string {
	key := jose.SigningKey{Algorithm: jose.RS256, Key: privateKey}
	var options = jose.SignerOptions{}
	options.WithType("JWT")
	signer, err := jose.NewSigner(key, &options)
	if err != nil {
		t.Error(err)
	}

	builder := jwt.Signed(signer)
	token, err := builder.Claims(claim).CompactSerialize()
	if err != nil {
		t.Error(err)
	}

	return token
}

func TestAuthenticator_Validate(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Error(err)
	}

	authenticator := Authenticator{
		ProjectID:  projectID,
		publicKeys: []*rsa.PublicKey{&privateKey.PublicKey},
	}

	tests := map[string]struct {
		input string
		want  error
	}{
		"valid token":          {input: generateJWT(validClaim, privateKey, t), want: nil},
		"wrong issuer token":   {input: generateJWT(wrongIssuerClaim, privateKey, t), want: jwt.ErrInvalidIssuer},
		"wrong subject token":  {input: generateJWT(wrongSubjectClaim, privateKey, t), want: jwt.ErrInvalidSubject},
		"wrong audience token": {input: generateJWT(wrongAudienceClaim, privateKey, t), want: jwt.ErrInvalidAudience},
		"future issued token":  {input: generateJWT(futureIssuedClaim, privateKey, t), want: jwt.ErrIssuedInTheFuture},
		"expired token":        {input: generateJWT(expiredClaim, privateKey, t), want: jwt.ErrExpired},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			claims, err := authenticator.Parse(tc.input)
			if err != nil {
				t.Error(err)
			}

			got := authenticator.Validate(claims)
			if tc.want != nil && !errors.Is(got, tc.want) {
				t.Fatalf("wanted: %s, got: %s", tc.want, err)
			} else if tc.want == nil && got != nil {
				t.Fatalf("wanted: nil, got: %s", err)
			}
		})
	}
}
