package diff2

import "sort"

// Side indicates which profile a key belongs to exclusively.
type Side int

const (
	SideLeft  Side = iota // key only in left profile
	SideRight             // key only in right profile
	SideBoth              // key in both profiles
)

// Entry represents a single key comparison result.
type Entry struct {
	Key    string
	Left   string
	Right  string
	Side   Side
	Equal  bool
}

// Result holds the full comparison between two profiles.
type Result struct {
	Entries []Entry
}

// OnlyLeft returns entries exclusive to the left profile.
func (r Result) OnlyLeft() []Entry {
	return r.filter(func(e Entry) bool { return e.Side == SideLeft })
}

// OnlyRight returns entries exclusive to the right profile.
func (r Result) OnlyRight() []Entry {
	return r.filter(func(e Entry) bool { return e.Side == SideRight })
}

// Changed returns entries present in both profiles but with differing values.
func (r Result) Changed() []Entry {
	return r.filter(func(e Entry) bool { return e.Side == SideBoth && !e.Equal })
}

// Unchanged returns entries present in both profiles with identical values.
func (r Result) Unchanged() []Entry {
	return r.filter(func(e Entry) bool { return e.Side == SideBoth && e.Equal })
}

func (r Result) filter(fn func(Entry) bool) []Entry {
	out := []Entry{}
	for _, e := range r.Entries {
		if fn(e) {
			out = append(out, e)
		}
	}
	return out
}

// Compute performs a two-way diff between left and right env maps.
func Compute(left, right map[string]string) Result {
	seen := make(map[string]bool)
	var entries []Entry

	for k, lv := range left {
		seen[k] = true
		if rv, ok := right[k]; ok {
			entries = append(entries, Entry{Key: k, Left: lv, Right: rv, Side: SideBoth, Equal: lv == rv})
		} else {
			entries = append(entries, Entry{Key: k, Left: lv, Side: SideLeft})
		}
	}

	for k, rv := range right {
		if !seen[k] {
			entries = append(entries, Entry{Key: k, Right: rv, Side: SideRight})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return Result{Entries: entries}
}
