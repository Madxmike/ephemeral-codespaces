package scheduling

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

type Handler struct {
	scheduler Scheduler
}

func (h *Handler) AddRoutes(r *chi.Mux) {
	r.Post("/", h.ScheduleNewDeployment)
}

// Schedule a new deployment
// Deleted a scheduled deployment
// Get a sheduled deployment by id
// Get all scheduled deployments for user_id
// Delete all scheduled deployments for user_id
// Reschedule a deployment with new details

func (h *Handler) ScheduleNewDeployment(w http.ResponseWriter, r *http.Request) {
	var d Deployment
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, "invalid deployment", http.StatusBadRequest)
		return
	}

	h.scheduler.Requests <- d
}

func (h *Handler) RemoveDeployment(w http.ResponseWriter, r *http.Request) {

}
