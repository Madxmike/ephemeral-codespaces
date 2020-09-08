package scheduling

import (
	"context"
	"sync"
	"time"
)

type Scheduler struct {
	sync.Mutex
	Requests chan Deployment

	queue deploymentQueue

	deployer Deployer

	update   *time.Ticker
	leadTime time.Duration
}

func NewScheduler(deployer Deployer, updateInterval time.Duration, leadTime time.Duration) *Scheduler {
	return &Scheduler{
		Requests: make(chan Deployment),
		queue:    newDeploymentQueue(),
		deployer: deployer,
		update:   time.NewTicker(updateInterval),
		leadTime: leadTime,
	}
}

// Start processing incoming deployment requests and begin deploying.
// Should be ran concurrently.
func (s *Scheduler) Run(ctx context.Context) {
	go s.processRequests(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.update.C:
			for _, d := range s.readyDeployments() {
				s.deployer.Deploy(ctx, d)
			}
		}
	}
}

func (s *Scheduler) readyDeployments() []Deployment {
	s.Lock()
	defer s.Unlock()

	var ready []Deployment

	var err error
	var next Deployment
	for err != nil {
		// Check
		next, err = s.queue.Peek()
		if err != nil && next.Ready(s.leadTime) {
			ready = append(ready, next)
			_, err = s.queue.Next()
		}
	}

	return ready
}

func (s *Scheduler) processRequests(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case req := <-s.Requests:
			s.Lock()
			s.queue.Add(req)
			s.Unlock()
		}
	}
}
