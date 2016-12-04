package classic

import (
	"sync"
)

// LockingSemaphore provides the functionality to limit bandwidth.
// https://en.wikipedia.org/wiki/Semaphore_(programming)#Operation_names
type LockingSemaphore interface {
	// P decrements the value of semaphore variable by n.
	// It must be safe to call P concurrently on a single semaphore.
	P(n int)
	// V increments the value of semaphore variable by n.
	// It must be safe to call V concurrently on a single semaphore.
	V(n int)
}

// NewLocking constructs a new LockingSemaphore with the given number of places.
func NewLocking(size int) LockingSemaphore {
	return make(semaphore, size)
}

func (sem semaphore) P(n int) {
	for i := 0; i < n; i++ {
		sem <- struct{}{}
	}
}

func (sem semaphore) V(n int) {
	for i := 0; i < n; i++ {
		<-sem
	}
}

// SyncingSemaphore provides the functionality to synchronize multiple gorutines.
// https://en.wikipedia.org/wiki/Semaphore_(programming)#Semantics_and_implementation
type SyncingSemaphore interface {
	// Signal reports on the completion of gorutine work.
	// It must be safe to call Signal concurrently on a single semaphore.
	Signal()
	// Wait starts to wait n gorutines.
	// It must be safe to call Wait concurrently on a single semaphore.
	Wait(n int)
}

// NewSyncing constructs a new SyncingSemaphore with the given number of places.
func NewSyncing(size int) SyncingSemaphore {
	sem := make(semaphore, size)
	sem.P(size)
	return sem
}

func (sem semaphore) Signal() {
	sem.V(1)
}

func (sem semaphore) Wait(n int) {
	sem.P(n)
}

// BinarySemaphore represents the classic binary semaphore with mutex-like interface.
type BinarySemaphore interface {
	sync.Locker
}

// NewBinary constructs a new BinarySemaphore with one place.
func NewBinary() BinarySemaphore {
	return make(semaphore, 1)
}

func (sem semaphore) Lock() {
	sem.P(1)
}

func (sem semaphore) Unlock() {
	sem.V(1)
}

type semaphore chan struct{}