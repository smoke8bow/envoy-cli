# envoy-cli

A lightweight CLI for managing and switching between named environment variable sets across projects.

---

## Installation

```bash
go install github.com/yourusername/envoy-cli@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/envoy-cli/releases).

---

## Usage

```bash
# Create a new environment set
envoy create myproject-dev

# Set variables in the environment set
envoy set myproject-dev DB_HOST=localhost DB_PORT=5432 API_KEY=abc123

# List all saved environment sets
envoy list

# Switch to an environment set (exports variables into your shell)
eval $(envoy use myproject-dev)

# Show variables in an environment set
envoy show myproject-dev

# Remove an environment set
envoy remove myproject-dev
```

Environment sets are stored locally in `~/.config/envoy/envs.json`, keeping your variables organized and portable across projects.

---

## Commands

| Command | Description |
|---|---|
| `create <name>` | Create a new named environment set |
| `set <name> KEY=VAL` | Add or update variables in a set |
| `use <name>` | Export variables from a set to the shell |
| `list` | List all available environment sets |
| `show <name>` | Display variables in a set |
| `remove <name>` | Delete an environment set |

---

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.

---

## License

[MIT](LICENSE)