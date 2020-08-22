package environment

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

type AuthValidator interface {
	Validate(token string) (bool, error)
}

func Routes(v AuthValidator) *chi.Mux {
	r := chi.NewRouter()

	r.Use(RequiresLogin(v))

	//Any route that the user is not authorized to see should 404

	// GET allows the user to get the status of an environment from the id
	// only allowing them to see environments that they are authed to see.
	//r.Get()
	// PUT allows the user to create an environment.
	//r.Put()

	// DELETE allows the user to request an environment be deleted. Must be authed.
	//r.Delete()

	r.Route("/{env_id}", func(r chi.Router) {
		r.Use(RequiresAuth)
	})

	return r
}

// Requires that the user is logged into the website to perform the following actions.
// Does not perform any permission checking for the action.
// Returns a 404 if the user is not logged in.
func RequiresLogin(v AuthValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bearer := r.Header.Get("Authorization")
			token := strings.TrimLeft(bearer, "Bearer")
			valid, err := v.Validate(token)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			if !valid {
				http.NotFound(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Requires that the client making the request is authorized to make a request for the environment.
// If they are not authorized then a 404 is returned.
// NYI - Will pass through.
func RequiresAuth(next http.Handler) http.Handler {
	return next
}
