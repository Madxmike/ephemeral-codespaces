package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"testing"
	"time"

	"github.com/square/go-jose/jwt"
	"gopkg.in/square/go-jose.v2"
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
		want  bool
	}{
		"empty input":          {input: "", want: false},
		"valid token":          {input: generateJWT(validClaim, privateKey, t), want: true},
		"wrong issuer token":   {input: generateJWT(wrongIssuerClaim, privateKey, t), want: false},
		"wrong subject token":  {input: generateJWT(wrongSubjectClaim, privateKey, t), want: false},
		"wrong audience token": {input: generateJWT(wrongAudienceClaim, privateKey, t), want: false},
		"future issued token":  {input: generateJWT(futureIssuedClaim, privateKey, t), want: false},
		"expired token":        {input: generateJWT(expiredClaim, privateKey, t), want: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := authenticator.Validate(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Fatalf("wanted: %t, got: %t", tc.want, got)
			}

		})
	}
}
