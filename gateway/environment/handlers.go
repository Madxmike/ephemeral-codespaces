package environment

import "net/http"

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
