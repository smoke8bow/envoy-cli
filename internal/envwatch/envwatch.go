// Package envwatch monitors a profile for changes and notifies callers
// when the environment variable set has drifted from a recorded baseline.
package envwatch

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

// ErrProfileNotFound is returned when the watched profile does not exist.
var ErrProfileNotFound = errors.New("envwatch: profile not found")

// ProfileGetter is the minimal interface required to read a profile's vars.
type ProfileGetter interface {
	Get(name string) (map[string]string, error)
}

// ChangeEvent describes a detected change in a watched profile.
type ChangeEvent struct {
	Profile  string
	PrevHash string
	CurrHash string
	At       time.Time
}

// Watcher polls a profile at a fixed interval and sends ChangeEvents on C
// whenever the content hash changes.
type Watcher struct {
	C       <-chan ChangeEvent
	profile string
	getter  ProfileGetter
	interval time.Duration

	mu       sync.Mutex
	lastHash string
	stop     chan struct{}
	done     chan struct{}
}

// NewWatcher creates a Watcher for the named profile. Polling begins
// immediately; call Stop to release resources.
func NewWatcher(getter ProfileGetter, profile string, interval time.Duration) (*Watcher, error) {
	if profile == "" {
		return nil, errors.New("envwatch: profile name must not be empty")
	}
	if interval <= 0 {
		return nil, errors.New("envwatch: interval must be positive")
	}

	// Compute the initial hash so we have a baseline before the first tick.
	vars, err := getter.Get(profile)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrProfileNotFound, profile)
	}
	initial := hashVars(vars)

	ch := make(chan ChangeEvent, 8)
	w := &Watcher{
		C:        ch,
		profile:  profile,
		getter:   getter,
		interval: interval,
		lastHash: initial,
		stop:     make(chan struct{}),
		done:     make(chan struct{}),
	}

	go w.poll(ch)
	return w, nil
}

// Stop halts the polling goroutine and closes the event channel.
func (w *Watcher) Stop() {
	close(w.stop)
	<-w.done
}

// LastHash returns the most recently computed hash for the watched profile.
func (w *Watcher) LastHash() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.lastHash
}

func (w *Watcher) poll(ch chan<- ChangeEvent) {
	defer close(ch)
	defer close(w.done)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-w.stop:
			return
		case <-ticker.C:
			vars, err := w.getter.Get(w.profile)
			if err != nil {
				// Profile may have been deleted; skip silently.
				continue
			}
			curr := hashVars(vars)

			w.mu.Lock()
			prev := w.lastHash
			if curr != prev {
				w.lastHash = curr
				w.mu.Unlock()
				select {
				case ch <- ChangeEvent{
					Profile:  w.profile,
					PrevHash: prev,
					CurrHash: curr,
					At:       time.Now(),
				}:
				case <-w.stop:
					return
				}
			} else {
				w.mu.Unlock()
			}
		}
	}
}

// hashVars returns a stable SHA-256 hex digest of the key=value pairs.
func hashVars(vars map[string]string) string {
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintf(h, "%s=%s\n", k, vars[k])
	}
	return hex.EncodeToString(h.Sum(nil))
}
