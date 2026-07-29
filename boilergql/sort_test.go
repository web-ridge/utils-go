package boilergql

import (
	"testing"
	"time"
)

func ptr[T any](v T) *T { return &v }

func TestToCursorValuePlainValues(t *testing.T) {
	if got := ToCursorValue("ID", 42); got != "ID:42" {
		t.Errorf("int: got %q, want %q", got, "ID:42")
	}
	if got := ToCursorValue("NAME", "hello"); got != "NAME:hello" {
		t.Errorf("string: got %q, want %q", got, "NAME:hello")
	}
}

func TestToCursorValueDereferencesPointers(t *testing.T) {
	if got := ToCursorValue("ID", ptr(42)); got != "ID:42" {
		t.Errorf("*int: got %q, want %q", got, "ID:42")
	}
	if got := ToCursorValue("NAME", ptr("hello")); got != "NAME:hello" {
		t.Errorf("*string: got %q, want %q", got, "NAME:hello")
	}
	if got := ToCursorValue("ID", ptr(ptr(7))); got != "ID:7" {
		t.Errorf("**int: got %q, want %q", got, "ID:7")
	}
}

func TestToCursorValueNilPointer(t *testing.T) {
	var nilInt *int
	if got := ToCursorValue("ID", nilInt); got != "ID:" {
		t.Errorf("nil *int: got %q, want %q", got, "ID:")
	}
	var nilTime *time.Time
	if got := ToCursorValue("CREATED_AT", nilTime); got != "CREATED_AT:" {
		t.Errorf("nil *time.Time: got %q, want %q", got, "CREATED_AT:")
	}
}

func TestToCursorValueTime(t *testing.T) {
	// The wall-clock reading must be preserved in the value's own location,
	// matching what the SQL driver read from the column.
	utc := time.Date(2026, 7, 30, 13, 37, 1, 0, time.UTC)
	if got := ToCursorValue("CREATED_AT", utc); got != "CREATED_AT:2026-07-30 13:37:01" {
		t.Errorf("time.Time UTC: got %q", got)
	}

	loc := time.FixedZone("CEST", 2*60*60)
	local := time.Date(2026, 7, 30, 13, 37, 1, 0, loc)
	if got := ToCursorValue("CREATED_AT", local); got != "CREATED_AT:2026-07-30 13:37:01" {
		t.Errorf("time.Time with zone: got %q", got)
	}

	if got := ToCursorValue("UPDATED_AT", ptr(utc)); got != "UPDATED_AT:2026-07-30 13:37:01" {
		t.Errorf("*time.Time: got %q", got)
	}

	// time.Time formatting does not depend on TimeSortKeyMatcher.
	if got := ToCursorValue("SOMETHING_ELSE", utc); got != "SOMETHING_ELSE:2026-07-30 13:37:01" {
		t.Errorf("time.Time with non-matching key: got %q", got)
	}
}

func TestToCursorValueTimeRoundTripsThroughFromCursorValue(t *testing.T) {
	cursor := ToCursorValue("CREATED_AT", time.Date(2026, 7, 30, 13, 37, 1, 0, time.UTC))
	k, v := FromCursorValue(cursor)
	if k != "CREATED_AT" || v != "2026-07-30 13:37:01" {
		t.Errorf("round trip: got (%q, %q)", k, v)
	}
}

func TestToCursorValueTimeSortKeyMatcherDisabledByDefault(t *testing.T) {
	if TimeSortKeyMatcher != nil {
		t.Fatal("TimeSortKeyMatcher should default to nil")
	}
	// Without a matcher, integers pass through untouched.
	if got := ToCursorValue("CREATED_AT", 1753875421); got != "CREATED_AT:1753875421" {
		t.Errorf("matcher off: got %q", got)
	}
	if got := ToCursorValue("CREATED_AT", ptr(1753875421)); got != "CREATED_AT:1753875421" {
		t.Errorf("matcher off, pointer: got %q", got)
	}
}

func TestToCursorValueTimeSortKeyMatcherEnabled(t *testing.T) {
	TimeSortKeyMatcher = SuffixTimeSortKeyMatcher("_AT")
	t.Cleanup(func() { TimeSortKeyMatcher = nil })

	want := "CREATED_AT:" + time.Unix(1753875421, 0).UTC().Format(CursorTimeLayout)
	if got := ToCursorValue("CREATED_AT", 1753875421); got != want {
		t.Errorf("matching key: got %q, want %q", got, want)
	}
	// Pointer to unix seconds dereferences first, then matches.
	if got := ToCursorValue("CREATED_AT", ptr(int64(1753875421))); got != want {
		t.Errorf("matching key, *int64: got %q, want %q", got, want)
	}
	// Non-matching keys keep the bare integer.
	if got := ToCursorValue("ID", 1753875421); got != "ID:1753875421" {
		t.Errorf("non-matching key: got %q", got)
	}
	// Non-integer values under a matching key are untouched.
	if got := ToCursorValue("CREATED_AT", "already-a-string"); got != "CREATED_AT:already-a-string" {
		t.Errorf("matching key, string value: got %q", got)
	}
}

func TestSuffixTimeSortKeyMatcher(t *testing.T) {
	m := SuffixTimeSortKeyMatcher("_AT", "_ON")
	for key, want := range map[string]bool{
		"CREATED_AT":  true,
		"UPDATED_AT":  true,
		"PLANNED_ON":  true,
		"ID":          false,
		"NAME":        false,
		"ATTACHMENTS": false,
	} {
		if got := m(key); got != want {
			t.Errorf("matcher(%q) = %v, want %v", key, got, want)
		}
	}
}
