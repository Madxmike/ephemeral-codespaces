package auth

import (
	"crypto/rsa"
	"testing"

	"github.com/pkg/errors"
	"github.com/tierzer0/gateway/mock"
	"gopkg.in/square/go-jose.v2/jwt"
)

func TestAuthenticator_Validate(t *testing.T) {
	authenticator := Authenticator{
		ProjectID:  mock.ProjectID,
		publicKeys: []*rsa.PublicKey{&mock.PrivateKey.PublicKey},
	}

	tests := map[string]struct {
		input string
		want  error
	}{
		"valid token":          {input: mock.GenerateJWT(mock.ValidClaim, t), want: nil},
		"wrong issuer token":   {input: mock.GenerateJWT(mock.WrongIssuerClaim, t), want: jwt.ErrInvalidIssuer},
		"wrong subject token":  {input: mock.GenerateJWT(mock.WrongSubjectClaim, t), want: jwt.ErrInvalidSubject},
		"wrong audience token": {input: mock.GenerateJWT(mock.WrongAudienceClaim, t), want: jwt.ErrInvalidAudience},
		"future issued token":  {input: mock.GenerateJWT(mock.FutureIssuedClaim, t), want: jwt.ErrIssuedInTheFuture},
		"expired token":        {input: mock.GenerateJWT(mock.ExpiredClaim, t), want: jwt.ErrExpired},
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
