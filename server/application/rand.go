package application

import "sync"

type Rand struct {
	mu *sync.Mutex
	lastNum int
}

func (r *Rand) getNextID() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.lastNum++
	return r.lastNum
}
