package environment

import (
	"net/http"

	"github.com/go-chi/chi"
)

func Routes() *chi.Mux {
	r := chi.NewRouter()

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

// Requires that the client making the request is authorized to make a request for the environment.
// If they are not authorized then a 404 is returned.
// NYI - Will pass through.
func RequiresAuth(next http.Handler) http.Handler {
	return next
}
