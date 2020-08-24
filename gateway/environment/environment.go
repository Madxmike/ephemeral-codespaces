package environment

import "time"

// An evironment defines the vs-code environment that will be created
// in a form that is easily created on the client. It will be transformed into an
// actual k8s resource by an operator.
type Environment struct {
	// The id corresponding to the environment that will be created.
	ID string `json:"id, omitempty"`

	// The user that is requesting the environment be created.
	Owner string `json:"owner"`

	// The time the client requested the environment be made
	CreatedAt time.Time `json:"created_at"`

	// All the software that needs to exist in the deployed environment
	Requires []Software `json:"requires"`
}

// A piece of software that is required in the image. Names and versions are
// mapped to actual images that are deployed.
type Software struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
