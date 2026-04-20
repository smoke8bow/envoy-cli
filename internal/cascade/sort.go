package cascade

import "sort"

// sortStrings sorts a string slice in-place (thin wrapper kept internal).
func sortStrings(s []string) { sort.Strings(s) }
