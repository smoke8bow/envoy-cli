package format

import (
	"fmt"
	"sort"
	"strings"
)

// Style represents an output format style.
type Style string

const (
	StyleTable Style = "table"
	StyleList  Style = "list"
	StyleCSV   Style = "csv"
)

// Supported returns all valid format styles.
func Supported() []Style {
	return []Style{StyleTable, StyleList, StyleCSV}
}

// IsSupported returns true if s is a known style.
func IsSupported(s Style) bool {
	for _, v := range Supported() {
		if v == s {
			return true
		}
	}
	return false
}

// Render formats a map of env vars according to the given style.
func Render(vars map[string]string, style Style) (string, error) {
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	switch style {
	case StyleTable:
		return renderTable(keys, vars), nil
	case StyleList:
		return renderList(keys, vars), nil
	case StyleCSV:
		return renderCSV(keys, vars), nil
	default:
		return "", fmt.Errorf("unsupported format style: %q", style)
	}
}

func renderTable(keys []string, vars map[string]string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s %s\n", "KEY", "VALUE"))
	sb.WriteString(strings.Repeat("-", 60) + "\n")
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("%-30s %s\n", k, vars[k]))
	}
	return sb.String()
}

func renderList(keys []string, vars map[string]string) string {
	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("%s=%s\n", k, vars[k]))
	}
	return sb.String()
}

func renderCSV(keys []string, vars map[string]string) string {
	var sb strings.Builder
	sb.WriteString("key,value\n")
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("%s,%s\n", k, vars[k]))
	}
	return sb.String()
}
