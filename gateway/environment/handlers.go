package environment

import (
	"context"
	"encoding/json"
	"net/http"
)

// Maintains a connection to the kubernetes cluser to easily obtain the status of any environment
type StatusHandler struct {
}

func (h *StatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

// Checks the current status of the environment pod on kubernetes.
// If the requester is not authorized to access the environment then an error will be returned.
func (h *StatusHandler) getStatus(id string, requester string) (status string, err error) {

	return "", nil
}

// A publisher allows a message to be published to some channel or topic
// E.G. Publish a message onto a Redis channel
type Publisher interface {
	Publish(ctx context.Context, channel string, message interface{}) error
}

type EnvironmentHandler struct {
	publisher Publisher
}

// POST
// Allows the user to start a new environment with the specified pieces of software installed.
// The client must submit a JWT from firebase auth to access the endpoint.
// Once parsing and verification of the environment is finished the request
// is sent to redis to be processed by a seperate operation.
func (h EnvironmentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var environment Environment

	err := json.NewDecoder(r.Body).Decode(&environment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ID, ok := r.Context().Value("ID").(string)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	environment.ID = ID

	h.publisher.Publish(r.Context(), "create", environment)
}
