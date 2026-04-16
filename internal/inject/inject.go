package inject

import (
	"fmt"
	"os/exec"
	"strings"
)

// Injector runs a subprocess with a given set of environment variables
// merged on top of the current process environment.
type Injector struct {
	base []string // base env (usually os.Environ())
}

// NewInjector creates an Injector with the provided base environment.
func NewInjector(base []string) *Injector {
	return &Injector{base: base}
}

// Build constructs the merged environment slice.
// Values in overlay take precedence over base.
func (inj *Injector) Build(overlay map[string]string) []string {
	merged := make(map[string]string, len(inj.base)+len(overlay))
	for _, kv := range inj.base {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) == 2 {
			merged[parts[0]] = parts[1]
		}
	}
	for k, v := range overlay {
		merged[k] = v
	}
	out := make([]string, 0, len(merged))
	for k, v := range merged {
		out = append(out, fmt.Sprintf("%s=%s", k, v))
	}
	return out
}

// Command returns an *exec.Cmd configured with the merged environment.
func (inj *Injector) Command(overlay map[string]string, name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Env = inj.Build(overlay)
	return cmd
}
