package deployment

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Handler struct {
	Deployer Deployer
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var d Deployment
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, "invalid deployment", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	err = h.Deployer.Deploy(r.Context(), d)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println(fmt.Errorf("could not deploy %s: %w", d.ID, err))
		return
	}
}
