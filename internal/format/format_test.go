package format

import (
	"strings"
	"testing"
)

func TestIsSupportedValid(t *testing.T) {
	for _, s := range Supported() {
		if !IsSupported(s) {
			t.Errorf("expected %q to be supported", s)
		}
	}
}

func TestIsSupportedInvalid(t *testing.T) {
	if IsSupported(Style("xml")) {
		t.Error("expected xml to be unsupported")
	}
}

func TestRenderList(t *testing.T) {
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := Render(vars, StyleList)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "FOO=bar") {
		t.Error("expected FOO=bar in list output")
	}
	if !strings.Contains(out, "BAZ=qux") {
		t.Error("expected BAZ=qux in list output")
	}
}

func TestRenderTable(t *testing.T) {
	vars := map[string]string{"HOST": "localhost"}
	out, err := Render(vars, StyleTable)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "KEY") {
		t.Error("expected KEY header in table output")
	}
	if !strings.Contains(out, "HOST") {
		t.Error("expected HOST in table output")
	}
	if !strings.Contains(out, "localhost") {
		t.Error("expected localhost in table output")
	}
}

func TestRenderCSV(t *testing.T) {
	vars := map[string]string{"PORT": "8080"}
	out, err := Render(vars, StyleCSV)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "key,value") {
		t.Error("expected csv header")
	}
	if !strings.Contains(out, "PORT,8080") {
		t.Error("expected PORT,8080 in csv output")
	}
}

func TestRenderSortedOutput(t *testing.T) {
	vars := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	out, err := Render(vars, StyleList)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "A_KEY") {
		t.Errorf("expected first line to be A_KEY, got %s", lines[0])
	}
}

func TestRenderUnsupportedStyle(t *testing.T) {
	_, err := Render(map[string]string{}, Style("yaml"))
	if err == nil {
		t.Error("expected error for unsupported style")
	}
}
