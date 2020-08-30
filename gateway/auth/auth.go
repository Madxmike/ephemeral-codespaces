package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/square/go-jose.v2/jwt"
)

//TODO - Move this into an env variable
var firebasePubKeyEndpoint = "https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com"

const secureTokenEndpoint = "https://securetoken.google.com/%s"

type Authenticator struct {
	// The project id corresponding to the Firebase auth project
	ProjectID string

	publicKeys []*rsa.PublicKey
}

func (a *Authenticator) RetrievePublicKeys() error {
	req, err := http.NewRequest("GET", firebasePubKeyEndpoint, nil)
	if err != nil {
		return errors.Wrap(err, "could not create public keys request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "could not retrieve public keys")
	}

	defer resp.Body.Close()

	certs := make(map[string]string)
	err = json.NewDecoder(resp.Body).Decode(&certs)
	if err != nil {
		return errors.Wrap(err, "could not read cert data body")
	}

	for _, cert := range certs {
		err = a.parseCerts([]byte(cert))
		if err != nil {
			return errors.Wrap(err, "could not parse certs")
		}
	}

	return nil
}

func (a *Authenticator) parseCerts(data []byte) error {
	if a.publicKeys == nil {
		a.publicKeys = make([]*rsa.PublicKey, 0)
	}
	var block *pem.Block
	for len(data) != 0 {
		block, data = pem.Decode(data)
		if block == nil {
			return errors.New("no pem block was found")
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return err
		}
		publicKey := cert.PublicKey.(*rsa.PublicKey)

		a.publicKeys = append(a.publicKeys, publicKey)
	}
	return nil
}

func (a *Authenticator) Parse(token string) (*jwt.Claims, error) {
	if a.publicKeys == nil {
		return nil, errors.New("no public keys are available")
	}

	if token == "" {
		return nil, errors.New("token is empty")
	}

	parsed, err := jwt.ParseSigned(token)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse JWT")
	}
	claims := new(jwt.Claims)
	for _, key := range a.publicKeys {
		// We dont care too much about the error here.
		// We only care that a claim gets parsed or not
		_ = parsed.Claims(key, claims)
		if claims != nil {
			break
		}
	}

	if claims == nil {
		return nil, errors.New("no claim was able to be parsed")
	}

	return claims, err
}

func (a *Authenticator) Validate(claims *jwt.Claims) error {
	expected := jwt.Expected{
		Issuer:  fmt.Sprintf(secureTokenEndpoint, a.ProjectID),
		Subject: a.ProjectID,
		Audience: jwt.Audience{
			a.ProjectID,
		},
		Time: time.Now(),
	}

	return claims.Validate(expected)
}
