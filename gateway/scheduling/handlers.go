package scheduling

import "net/http"

type SchedulingHandler struct {
	scheduler Scheduler
}

// Post - Expects a Deployment in the body. Will schedule the deployment.
func (h *SchedulingHandler) ScheduleDeployment(w *http.ResponseWriter, r http.Request) {

}
