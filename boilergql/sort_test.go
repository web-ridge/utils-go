package boilergql

import (
	"testing"
	"time"

	"github.com/aarondl/null/v8"
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

type cursorColumnModel struct {
	ID         uint       `boil:"id"`
	LinkedDate null.Time  `boil:"linked_date"`
	CreatedAt  time.Time  `boil:"created_at"`
	DeletedAt  *time.Time `boil:"deleted_at"`
	Minutes    null.Int   `boil:"minutes"`
	R          *struct{}  `boil:"-"`
}

func TestCursorValueForColumn(t *testing.T) {
	for _, test := range []struct {
		name   string
		column string
		value  string
		want   interface{}
	}{
		{name: "nullable time column from unix seconds", column: "linked_date", value: "1753875421", want: time.Unix(1753875421, 0)},
		{name: "time column from unix seconds", column: "created_at", value: "1753875421", want: time.Unix(1753875421, 0)},
		{name: "pointer time column from unix seconds", column: "deleted_at", value: "1753875421", want: time.Unix(1753875421, 0)},
		{name: "time column already in cursor layout", column: "created_at", value: "2026-07-30 13:37:01", want: "2026-07-30 13:37:01"},
		{name: "time column without a comparable value", column: "linked_date", value: "not-a-time", want: nil},
		{name: "null time column", column: "linked_date", value: "", want: nil},
		{name: "int column", column: "minutes", value: "60", want: "60"},
		{name: "unknown column", column: "nope", value: "60", want: "60"},
		{name: "empty value", column: "minutes", value: "", want: nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := CursorValueForColumn(cursorColumnModel{}, test.column, test.value); got != test.want {
				t.Errorf("CursorValueForColumn(%q, %q) = %v, want %v", test.column, test.value, got, test.want)
			}
		})
	}
}

func TestCursorValueForColumnWithoutModel(t *testing.T) {
	if got := CursorValueForColumn(nil, "created_at", "1753875421"); got != "1753875421" {
		t.Errorf("nil model: got %v", got)
	}
	if got := CursorValueForColumn(&cursorColumnModel{}, "created_at", "1753875421"); got != "1753875421" {
		t.Errorf("pointer model: got %v", got)
	}
}

func TestCursorValueForColumnRoundTripsTimeCursors(t *testing.T) {
	_, value := FromCursorValue(ToCursorValue("CREATED_AT", time.Date(2026, 7, 30, 13, 37, 1, 0, time.UTC)))
	if got := CursorValueForColumn(cursorColumnModel{}, "created_at", value); got != "2026-07-30 13:37:01" {
		t.Errorf("time.Time cursor: got %v", got)
	}

	// A GraphQL model exposing the column as unix seconds round trips into a
	// time.Time the SQL driver can bind against a DATETIME column.
	_, value = FromCursorValue(ToCursorValue("LINKED_DATE", ptr(1753875421)))
	if got := CursorValueForColumn(cursorColumnModel{}, "linked_date", value); got != time.Unix(1753875421, 0) {
		t.Errorf("unix seconds cursor: got %v", got)
	}
}
