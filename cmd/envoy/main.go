// main is the entry point for the envoy-cli tool.
// It wires together all internal packages and registers
// the top-level cobra commands.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// newRootCmd builds the top-level cobra command and attaches all sub-commands.
func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "envoy",
		Short: "Manage and switch between named environment variable sets",
		Long: `envoy-cli lets you create, edit, and switch between named
profiles of environment variables across projects.

Examples:
  envoy profile create dev
  envoy profile set dev DB_HOST=localhost
  envoy switch dev
  envoy export --format dotenv dev`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version,
	}

	// Global flags
	var storePath string
	root.PersistentFlags().StringVar(
		&storePath, "store", "",
		"path to the envoy store directory (default: $HOME/.config/envoy)",
	)

	// Sub-command groups
	root.AddCommand(
		newProfileCmd(&storePath),
		newSwitchCmd(&storePath),
		newExportCmd(&storePath),
		newImportCmd(&storePath),
		newCompletionCmd(),
	)

	return root
}

// newProfileCmd returns the "profile" sub-command group.
func newProfileCmd(storePath *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage environment profiles",
	}
	cmd.AddCommand(
		newProfileCreateCmd(storePath),
		newProfileListCmd(storePath),
		newProfileDeleteCmd(storePath),
	)
	return cmd
}

func newProfileCreateCmd(storePath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new empty profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = storePath // resolved in future wiring
			fmt.Fprintf(cmd.OutOrStdout(), "created profile %q\n", args[0])
			return nil
		},
	}
}

func newProfileListCmd(storePath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = storePath
			fmt.Fprintln(cmd.OutOrStdout(), "(no profiles yet)")
			return nil
		},
	}
}

func newProfileDeleteCmd(storePath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = storePath
			fmt.Fprintf(cmd.OutOrStdout(), "deleted profile %q\n", args[0])
			return nil
		},
	}
}

func newSwitchCmd(storePath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "switch <profile>",
		Short: "Switch to a named profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = storePath
			fmt.Fprintf(cmd.OutOrStdout(), "switched to profile %q\n", args[0])
			return nil
		},
	}
}

func newExportCmd(storePath *string) *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "export <profile>",
		Short: "Export a profile in the requested format",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = storePath
			fmt.Fprintf(cmd.OutOrStdout(), "# export %s (format=%s)\n", args[0], format)
			return nil
		},
	}
	cmd.Flags().StringVarP(&format, "format", "f", "shell", "output format: shell, dotenv, json")
	return cmd
}

func newImportCmd(storePath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "import <file>",
		Short: "Import a .env or JSON file as a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = storePath
			fmt.Fprintf(cmd.OutOrStdout(), "imported from %q\n", args[0])
			return nil
		},
	}
}

func newCompletionCmd() *cobra.Command {
	return &cobra.Command{
		Use:       "completion <shell>",
		Short:     "Generate shell completion scripts",
		ValidArgs: []string{"bash", "zsh", "fish"},
		Args:      cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := cmd.Root()
			switch args[0] {
			case "bash":
				return root.GenBashCompletion(cmd.OutOrStdout())
			case "zsh":
				return root.GenZshCompletion(cmd.OutOrStdout())
			case "fish":
				return root.GenFishCompletion(cmd.OutOrStdout(), true)
			default:
				return fmt.Errorf("unsupported shell %q", args[0])
			}
		},
	}
}
