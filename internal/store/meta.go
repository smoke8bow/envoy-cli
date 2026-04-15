package store

import "fmt"

// SetMeta stores an arbitrary key-value metadata string for a named profile.
// The profile must already exist in the store.
func (s *Store) SetMeta(profile, key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data.Profiles[profile]; !ok {
		return fmt.Errorf("profile %q not found", profile)
	}
	if s.data.Meta == nil {
		s.data.Meta = make(map[string]map[string]string)
	}
	if s.data.Meta[profile] == nil {
		s.data.Meta[profile] = make(map[string]string)
	}
	s.data.Meta[profile][key] = value
	return s.save()
}

// GetMeta retrieves a metadata value for a profile key.
// Returns an empty string (and no error) when the key is absent.
func (s *Store) GetMeta(profile, key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.data.Profiles[profile]; !ok {
		return "", fmt.Errorf("profile %q not found", profile)
	}
	if s.data.Meta == nil {
		return "", nil
	}
	if m, ok := s.data.Meta[profile]; ok {
		return m[key], nil
	}
	return "", nil
}

// DeleteMeta removes a single metadata key for a profile.
func (s *Store) DeleteMeta(profile, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data.Meta == nil {
		return nil
	}
	if m, ok := s.data.Meta[profile]; ok {
		delete(m, key)
		if len(m) == 0 {
			delete(s.data.Meta, profile)
		}
	}
	return s.save()
}
