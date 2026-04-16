package resolve

import (
	"errors"
	"testing"
)

type fakeAccessor struct {
	profiles map[string]map[string]string
}

func (f *fakeAccessor) Get(name string) (map[string]string, error) {
	v, ok := f.profiles[name]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}

func TestResolveProfileSuccess(t *testing.T) {
	fa := &fakeAccessor{
		profiles: map[string]map[string]string{
			"dev": {"BASE": "/app", "LOG": "${BASE}/log"},
		},
	}
	out, err := ResolveProfile(fa, "dev", nil)
	if err != nil {
		t.Fatal(err)
	}
	if out["LOG"] != "/app/log" {
		t.Errorf("got %q", out["LOG"])
	}
}

func TestResolveProfileNotFound(t *testing.T) {
	fa := &fakeAccessor{profiles: map[string]map[string]string{}}
	_, err := ResolveProfile(fa, "missing", nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestResolveProfileWithAmbient(t *testing.T) {
	fa := &fakeAccessor{
		profiles: map[string]map[string]string{
			"prod": {"URL": "https://$DOMAIN/api"},
		},
	}
	ambient := map[string]string{"DOMAIN": "example.com"}
	out, err := ResolveProfile(fa, "prod", ambient)
	if err != nil {
		t.Fatal(err)
	}
	if out["URL"] != "https://example.com/api" {
		t.Errorf("got %q", out["URL"])
	}
}
