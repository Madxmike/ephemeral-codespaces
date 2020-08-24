package environment

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"gopkg.in/square/go-jose.v2/jwt"
)

type TokenAuthorizer interface {
	Parse(token string) (*jwt.Claims, error)
	Validate(claims *jwt.Claims) error
}

func Routes(auth TokenAuthorizer, publisher Publisher) *chi.Mux {
	r := chi.NewRouter()

	r.Use(RequiresAuthorization(auth))
	r.Method("POST", "/", EnvironmentHandler{
		publisher: publisher,
	})

	//Any route that the user is not authorized to see should 404

	// GET allows the user to get the status of an environment from the id
	// only allowing them to see environments that they are authed to see.
	//r.Get()

	// DELETE allows the user to request an environment be deleted. Must be authed.
	//r.Delete()

	return r
}

// Requires that the user is authed to perform the following actions.
// Does not perform any permission checking for the action.
// Returns a 404 if the user is not authorized to perform the action.
func RequiresAuthorization(auth TokenAuthorizer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bearer := r.Header.Get("Authorization")
			token := strings.TrimLeft(bearer, "Bearer ")

			claims, err := auth.Parse(token)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			err = auth.Validate(claims)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "ID", claims.ID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
