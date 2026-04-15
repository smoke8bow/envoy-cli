package tag

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/envoy-cli/envoy-cli/internal/store"
)

// ErrTagNotFound is returned when a tag does not exist.
var ErrTagNotFound = errors.New("tag not found")

// Manager handles tagging of profiles.
type Manager struct {
	store *store.Store
}

// NewManager returns a new tag Manager backed by the given store.
func NewManager(s *store.Store) *Manager {
	return &Manager{store: s}
}

// Add attaches a tag to a profile. Duplicate tags are silently ignored.
func (m *Manager) Add(profile, tag string) error {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return errors.New("tag must not be empty")
	}
	tags, err := m.List(profile)
	if err != nil {
		return err
	}
	for _, t := range tags {
		if t == tag {
			return nil
		}
	}
	tags = append(tags, tag)
	sort.Strings(tags)
	return m.store.SetMeta(profile, "tags", strings.Join(tags, ","))
}

// Remove detaches a tag from a profile.
func (m *Manager) Remove(profile, tag string) error {
	tags, err := m.List(profile)
	if err != nil {
		return err
	}
	filtered := tags[:0]
	for _, t := range tags {
		if t != tag {
			filtered = append(filtered, t)
		}
	}
	if len(filtered) == len(tags) {
		return fmt.Errorf("%w: %s", ErrTagNotFound, tag)
	}
	return m.store.SetMeta(profile, "tags", strings.Join(filtered, ","))
}

// List returns all tags attached to a profile, sorted alphabetically.
func (m *Manager) List(profile string) ([]string, error) {
	raw, err := m.store.GetMeta(profile, "tags")
	if err != nil {
		return nil, err
	}
	if raw == "" {
		return []string{}, nil
	}
	parts := strings.Split(raw, ",")
	sort.Strings(parts)
	return parts, nil
}

// ProfilesWithTag returns all profile names that carry the given tag.
func (m *Manager) ProfilesWithTag(tag string) ([]string, error) {
	profiles, err := m.store.List()
	if err != nil {
		return nil, err
	}
	var matched []string
	for _, p := range profiles {
		tags, err := m.List(p)
		if err != nil {
			return nil, err
		}
		for _, t := range tags {
			if t == tag {
				matched = append(matched, p)
				break
			}
		}
	}
	return matched, nil
}
