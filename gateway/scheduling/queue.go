package scheduling

import (
	"container/heap"
	"errors"
)

var ErrQueueEmpty = errors.New("queue is empty")

type deploymentQueue []Deployment

func newDeploymentQueue(initialValues ...Deployment) deploymentQueue {
	q := deploymentQueue(initialValues)
	heap.Init(&q)
	return q
}

// Adds a deployment to the queue
func (q *deploymentQueue) Add(d Deployment) {
	heap.Push(q, d)
}

// Pops the next deployment off the queue and returns it if available.
// If no deployments are in the queue then ErrEmptyQueue is returned.
func (q *deploymentQueue) Next() (Deployment, error) {
	if q.Len() == 0 {
		return Deployment{}, ErrQueueEmpty
	}
	next := heap.Pop(q).(Deployment)
	return next, nil
}

// Peeks at the next deployment off the queue, but does not remove the element.
// Returns the next deployment if available, ErrEmptyQueue if not.
func (q deploymentQueue) Peek() (Deployment, error) {
	if len(q) == 0 {
		return Deployment{}, ErrQueueEmpty
	}

	return q[len(q)-1], nil
}

// Internal interface implementation. Use DeploymentQueue.Add instead.
func (q *deploymentQueue) Push(x interface{}) {
	item := x.(Deployment)
	*q = append(*q, item)
}

// Internal interface implementation. Use DeploymentQueue.Next instead.
func (q *deploymentQueue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	*q = old[0 : n-1]
	return item
}

func (q deploymentQueue) Swap(i int, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q deploymentQueue) Less(i int, j int) bool {
	return q[i].StartAt.Before(q[j].StartAt)
}

func (q deploymentQueue) Len() int {
	return len(q)
}
