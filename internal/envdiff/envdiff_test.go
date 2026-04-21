package envdiff

import (
	"strings"
	"testing"
)

func TestOnlyInProfile(t *testing.T) {
	t.Setenv("EXISTING_KEY", "os_value")

	profile := map[string]string{
		"NEW_KEY": "new_value",
	}

	r := Compare(profile)

	if len(r.OnlyInProfile) != 1 || r.OnlyInProfile[0].Key != "NEW_KEY" {
		t.Fatalf("expected NEW_KEY in OnlyInProfile, got %v", r.OnlyInProfile)
	}
	if r.OnlyInProfile[0].Source != SourceProfile {
		t.Errorf("expected source=profile, got %s", r.OnlyInProfile[0].Source)
	}
}

func TestInBoth(t *testing.T) {
	t.Setenv("SHARED_KEY", "os_value")

	profile := map[string]string{
		"SHARED_KEY": "profile_value",
	}

	r := Compare(profile)

	var found *Entry
	for i := range r.InBoth {
		if r.InBoth[i].Key == "SHARED_KEY" {
			found = &r.InBoth[i]
			break
		}
	}
	if found == nil {
		t.Fatal("expected SHARED_KEY in InBoth")
	}
	if found.ProfileValue != "profile_value" {
		t.Errorf("unexpected profile value: %s", found.ProfileValue)
	}
	if found.OSValue != "os_value" {
		t.Errorf("unexpected os value: %s", found.OSValue)
	}
	if found.Source != SourceBoth {
		t.Errorf("expected source=both, got %s", found.Source)
	}
}

func TestOnlyInOS(t *testing.T) {
	t.Setenv("OS_ONLY_KEY", "some_value")

	profile := map[string]string{}

	r := Compare(profile)

	var found bool
	for _, e := range r.OnlyInOS {
		if e.Key == "OS_ONLY_KEY" {
			found = true
			if e.Source != SourceOS {
				t.Errorf("expected source=os, got %s", e.Source)
			}
		}
	}
	if !found {
		t.Error("expected OS_ONLY_KEY in OnlyInOS")
	}
}

func TestResultIsSorted(t *testing.T) {
	profile := map[string]string{
		"ZEBRA": "z",
		"ALPHA": "a",
		"MANGO": "m",
	}

	r := Compare(profile)

	keys := make([]string, len(r.OnlyInProfile))
	for i, e := range r.OnlyInProfile {
		keys[i] = e.Key
	}
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("entries not sorted: %v", keys)
		}
	}
}

func TestFormatContainsPlusForProfileOnly(t *testing.T) {
	profile := map[string]string{"BRAND_NEW": "val"}
	r := Compare(profile)
	out := Format(r)
	if !strings.Contains(out, "+ BRAND_NEW=val") {
		t.Errorf("expected '+ BRAND_NEW=val' in output:\n%s", out)
	}
}

func TestFormatTildeForMismatch(t *testing.T) {
	t.Setenv("DIFF_KEY", "os_val")
	profile := map[string]string{"DIFF_KEY": "profile_val"}
	r := Compare(profile)
	out := Format(r)
	if !strings.Contains(out, "~ DIFF_KEY") {
		t.Errorf("expected '~ DIFF_KEY' in output:\n%s", out)
	}
}
