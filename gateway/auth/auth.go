package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
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

	certData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "could not read cert data body")
	}

	err = a.parseCerts(certData)
	if err != nil {
		return errors.Wrap(err, "could not parse certs")
	}

	return nil
}

func (a *Authenticator) parseCerts(data []byte) error {
	a.publicKeys = make([]*rsa.PublicKey, 0)
	var block *pem.Block
	for len(data) != 0 {
		block, data = pem.Decode(data)
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return err
		}
		publicKey := cert.PublicKey.(*rsa.PublicKey)

		a.publicKeys = append(a.publicKeys, publicKey)
	}
	return nil
}

func (a *Authenticator) Validate(token string) (bool, error) {
	if a.publicKeys == nil {
		return false, errors.New("no public keys are available")
	}

	if token == "" {
		return false, nil
	}

	parsed, err := jwt.ParseSigned(token)
	if err != nil {
		return false, errors.Wrap(err, "could not parse JWT")
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
		return false, errors.New("no claim was able to be parsed")
	}

	expected := jwt.Expected{
		Issuer:  fmt.Sprintf(secureTokenEndpoint, a.ProjectID),
		Subject: a.ProjectID,
		Audience: jwt.Audience{
			a.ProjectID,
		},
		Time: time.Now(),
	}

	err = claims.Validate(expected)
	if err != nil {
		return false, nil
	}

	return true, nil
}
