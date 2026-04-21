package interpolate

import "os"

// InterpolateProfile resolves variable references within the named profile.
// It uses the profile's own variables as the primary source and falls back
// to OS environment variables when FallbackToOS is enabled.
func InterpolateProfile(g Getter, name string, opts Options) (map[string]string, error) {
	vars, err := g.Get(name)
	if err != nil {
		return nil, err
	}
	interp := New(g, opts)
	return interp.Apply(vars, os.Getenv)
}
