package scheduling

import (
	"context"
	"time"

	"github.com/tierzer0/gateway/auth"
)

type Deployer interface {
	Deploy(context.Context, Deployment) error
}

type Deployment struct {
	ID string `json:"ID"`

	CreatedBy auth.User `json:"created_by"`

	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`

	Extensions []extension `json:"extensions"`
	Runtimes   []runtime   `json:"runtimes"`
}

//A VSCode extension to be installed during deployment
type extension struct {
	ID      string `json:"ID"`
	Enabled bool   `json:"enabled"`
}

//A runtime needed to be isntalled during deployment
type runtime struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func NewDeployment(ID string, createdBy auth.User, startAt time.Time, endAt time.Time) Deployment {
	return Deployment{
		ID:         ID,
		CreatedBy:  createdBy,
		StartAt:    startAt,
		EndAt:      endAt,
		Extensions: make([]extension, 0),
		Runtimes:   make([]runtime, 0),
	}
}

func (d *Deployment) AddExtension(ID string, enabled bool) {
	d.Extensions = append(d.Extensions, extension{
		ID:      ID,
		Enabled: enabled,
	})
}

func (d *Deployment) AddRuntime(name string, version string) {
	d.Runtimes = append(d.Runtimes, runtime{
		Name:    name,
		Version: version,
	})
}

// Does the deployment start within the current time + lead time
func (d Deployment) Ready(leadTime time.Duration) bool {
	return d.StartAt.Before(time.Now().Add(leadTime))
}
