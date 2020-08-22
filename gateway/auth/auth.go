package auth

import (
	"context"
	"encoding/json"
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

	publicKeys certs
}

type certs struct {
	First  string `json:"12809dd239d24bd379c0ad191f8b0edcdb9d3914"`
	Second string `json:"49e88c53761996a73623f191d512d2b47df802a1"`
}

func (a *Authenticator) RetrievePublicKeys(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", firebasePubKeyEndpoint, nil)
	if err != nil {
		return errors.Wrap(err, "could not create public keys request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "could not retrieve public keys")
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&a.publicKeys)
	if err != nil {
		return errors.Wrap(err, "could not parse public keys")
	}

	return nil
}

func (a *Authenticator) Validate(token string) (bool, error) {
	parsed, err := jwt.ParseSigned(token)
	if err != nil {
		return false, errors.Wrap(err, "could not parse JWT")
	}
	var claims jwt.Claims
	err = parsed.Claims(a.publicKeys, &claims)
	if err != nil {
		return false, errors.Wrap(err, "could not parse JWT claims")
	}

	expected := jwt.Expected{
		Issuer: fmt.Sprintf(secureTokenEndpoint, a.ProjectID),
		Audience: jwt.Audience{
			a.ProjectID,
		},
		Time: time.Now(),
	}
	err = claims.Validate(expected)
	if err != nil {
		return false, errors.Wrap(err, "token is not valid")
	}

	return true, nil
}
