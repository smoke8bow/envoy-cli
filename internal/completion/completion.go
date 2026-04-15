package completion

import (
	"fmt"
	"strings"
)

// Shell represents a supported shell type for completion scripts.
type Shell string

const (
	Bash Shell = "bash"
	Zsh  Shell = "zsh"
	Fish Shell = "fish"
)

// ProfileLister is implemented by any type that can list profile names.
type ProfileLister interface {
	List() ([]string, error)
}

// Generator generates shell completion scripts for envoy-cli commands.
type Generator struct {
	shell   Shell
	lister  ProfileLister
}

// NewGenerator returns a new Generator for the given shell.
func NewGenerator(shell Shell, lister ProfileLister) *Generator {
	return &Generator{shell: shell, lister: lister}
}

// ProfileNames returns a newline-separated list of profile names
// suitable for use in shell completion.
func (g *Generator) ProfileNames() (string, error) {
	names, err := g.lister.List()
	if err != nil {
		return "", fmt.Errorf("completion: listing profiles: %w", err)
	}
	return strings.Join(names, "\n"), nil
}

// Script returns a shell-specific completion script for envoy-cli.
func (g *Generator) Script(programName string) (string, error) {
	switch g.shell {
	case Bash:
		return bashScript(programName), nil
	case Zsh:
		return zshScript(programName), nil
	case Fish:
		return fishScript(programName), nil
	default:
		return "", fmt.Errorf("completion: unsupported shell %q", g.shell)
	}
}

func bashScript(prog string) string {
	return fmt.Sprintf(`_%s_completions() {
  local cur
  cur="${COMP_WORDS[COMP_CWORD]}"
  COMPREPLY=($(compgen -W "$(%s list)" -- "$cur"))
}
complete -F _%s_completions %s
`, prog, prog, prog, prog)
}

func zshScript(prog string) string {
	return fmt.Sprintf(`#compdef %s
_%s() {
  local -a profiles
  profiles=($(%s list))
  _describe 'profile' profiles
}
_%s "$@"
`, prog, prog, prog, prog)
}

func fishScript(prog string) string {
	return fmt.Sprintf(`complete -c %s -f -a "(%s list)" -d 'profile'
`, prog, prog)
}
