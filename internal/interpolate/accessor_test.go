package interpolate

import (
	"errors"
	"testing"
)

type fakeGetterAccessor struct {
	profiles map[string]map[string]string
}

func (f *fakeGetterAccessor) Get(name string) (map[string]string, error) {
	p, ok := f.profiles[name]
	if !ok {
		return nil, errors.New("not found: " + name)
	}
	return p, nil
}

func TestInterpolateProfileSuccess(t *testing.T) {
	g := &fakeGetterAccessor{
		profiles: map[string]map[string]string{
			"prod": {"BASE": "/app", "LOG": "${BASE}/logs"},
		},
	}
	out, err := InterpolateProfile(g, "prod", DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if out["LOG"] != "/app/logs" {
		t.Fatalf("expected /app/logs, got %s", out["LOG"])
	}
}

func TestInterpolateProfileNotFound(t *testing.T) {
	g := &fakeGetterAccessor{profiles: map[string]map[string]string{}}
	_, err := InterpolateProfile(g, "missing", DefaultOptions())
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestInterpolateProfileNoRefs(t *testing.T) {
	g := &fakeGetterAccessor{
		profiles: map[string]map[string]string{
			"dev": {"PORT": "3000", "HOST": "localhost"},
		},
	}
	out, err := InterpolateProfile(g, "dev", DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if out["PORT"] != "3000" || out["HOST"] != "localhost" {
		t.Fatalf("unexpected output: %v", out)
	}
}
