package mock

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"testing"
	"time"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

const ProjectID = "test_project"
const secureTokenEndpoint = "https://securetoken.google.com/%s"

var ValidClaim = &jwt.Claims{
	Issuer:   fmt.Sprintf(secureTokenEndpoint, ProjectID),
	Subject:  ProjectID,
	ID:       ProjectID,
	Audience: jwt.Audience{ProjectID},
	IssuedAt: jwt.NewNumericDate(time.Now()),
	Expiry:   jwt.NewNumericDate(time.Now()),
}

var WrongIssuerClaim = &jwt.Claims{
	Issuer:   fmt.Sprintf(secureTokenEndpoint, "wrong"),
	Subject:  ProjectID,
	ID:       ProjectID,
	Audience: jwt.Audience{ProjectID},
	IssuedAt: jwt.NewNumericDate(time.Now()),
	Expiry:   jwt.NewNumericDate(time.Now()),
}

var WrongSubjectClaim = &jwt.Claims{
	Issuer:   fmt.Sprintf(secureTokenEndpoint, ProjectID),
	Subject:  "wrong",
	ID:       ProjectID,
	Audience: jwt.Audience{ProjectID},
	IssuedAt: jwt.NewNumericDate(time.Now()),
	Expiry:   jwt.NewNumericDate(time.Now()),
}

var WrongAudienceClaim = &jwt.Claims{
	Issuer:   fmt.Sprintf(secureTokenEndpoint, ProjectID),
	Subject:  ProjectID,
	ID:       ProjectID,
	Audience: jwt.Audience{"wrong"},
	IssuedAt: jwt.NewNumericDate(time.Now()),
	Expiry:   jwt.NewNumericDate(time.Now()),
}

var FutureIssuedClaim = &jwt.Claims{
	Issuer:   fmt.Sprintf(secureTokenEndpoint, ProjectID),
	Subject:  ProjectID,
	ID:       ProjectID,
	Audience: jwt.Audience{ProjectID},
	IssuedAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
	Expiry:   jwt.NewNumericDate(time.Now()),
}

var ExpiredClaim = &jwt.Claims{
	Issuer:   fmt.Sprintf(secureTokenEndpoint, ProjectID),
	Subject:  ProjectID,
	ID:       ProjectID,
	Audience: jwt.Audience{ProjectID},
	IssuedAt: jwt.NewNumericDate(time.Now()),
	Expiry:   jwt.NewNumericDate(time.Time{}.Add(1 * time.Hour)),
}

var PrivateKey *rsa.PrivateKey

func init() {
	var err error
	PrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}
}

func GenerateJWT(claim *jwt.Claims, t *testing.T) string {
	key := jose.SigningKey{Algorithm: jose.RS256, Key: PrivateKey}
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
