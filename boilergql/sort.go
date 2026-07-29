package boilergql

import (
	"encoding/base64"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

type CursorType string

const (
	CursorTypeOffset CursorType = "OFFSET"
	CursorTypeCursor CursorType = "CURSOR"
)

func GetIDFromCursor(id string) interface{} {
	splitID := strings.SplitN(id, IDSeparator, 2)
	if len(splitID) != 2 {
		return 0
	}
	return splitID[1]
}

func GetModelFromCursor(id string) string {
	splitID := strings.SplitN(id, IDSeparator, 2)
	if len(splitID) != 2 {
		return ""
	}
	return splitID[0]
}

func ZeroOrMore(limit int) int {
	if limit < 0 {
		return 0
	}
	return limit
}

type ComparisonSign string

const (
	ComparisonSignBiggerThan         ComparisonSign = ">"
	ComparisonSignBiggerThanOrEqual  ComparisonSign = ">="
	ComparisonSignSmallerThan        ComparisonSign = "<"
	ComparisonSignSmallerThanOrEqual ComparisonSign = "<="
)

func GetComparison(
	forward *ConnectionForwardPagination,
	backward *ConnectionBackwardPagination, reverse bool,
	direction SortDirection,
) ComparisonSign {
	if forward != nil {
		if direction == SortDirectionDesc {
			return getForwardComparisonDesc(reverse)
		}
		return getForwardComparison(reverse)
	}
	if backward != nil {
		if direction == SortDirectionAsc {
			return getBackwardComparisonAsc(reverse)
		}
		return getBackwardComparison(reverse)
	}
	return ""
}

func getForwardComparison(reverse bool) ComparisonSign {
	if reverse {
		return ComparisonSignSmallerThanOrEqual
	}
	return ComparisonSignBiggerThan
}

func getForwardComparisonDesc(reverse bool) ComparisonSign {
	if reverse {
		return ComparisonSignBiggerThanOrEqual
	}
	return ComparisonSignSmallerThan
}

func getBackwardComparison(reverse bool) ComparisonSign {
	if reverse {
		return ComparisonSignBiggerThan
	}
	return ComparisonSignSmallerThanOrEqual
}

func getBackwardComparisonAsc(reverse bool) ComparisonSign {
	if reverse {
		return ComparisonSignSmallerThan
	}
	return ComparisonSignBiggerThanOrEqual
}

func GetCursor(forward *ConnectionForwardPagination, backward *ConnectionBackwardPagination) *string {
	if forward != nil {
		return forward.After
	}
	if backward != nil {
		return backward.Before
	}
	return nil
}

func GetLimit(forward *ConnectionForwardPagination, backward *ConnectionBackwardPagination) int {
	if forward != nil {
		return ZeroOrMore(forward.First + 1)
	}
	if backward != nil {
		return ZeroOrMore(backward.Last + 1)
	}
	return 0
}

func GetOffsetFromCursor(cursor *string) int {
	if cursor == nil {
		return 0
	}
	i, _ := strconv.Atoi(*cursor) //nolint:errcheck
	return i
}

func FromOffsetCursor(cursor string) []qm.QueryMod {
	offset, _ := strconv.Atoi(cursor) //nolint:errcheck

	if offset > 0 {
		return []qm.QueryMod{
			qm.Offset(offset),
		}
	}
	return nil
}

func GetDirection(direction SortDirection, reverse bool) SortDirection {
	if reverse {
		if direction == SortDirectionAsc {
			return SortDirectionDesc
		}
		return SortDirectionAsc
	}
	return direction
}

func GetOrderBy(dbColumn string, direction SortDirection) string {
	return dbColumn + " " + string(direction)
}

func CursorTypeCounter() (func(SortDirection), func() CursorType) {
	var asc, desc int

	return func(d SortDirection) {
			switch d {
			case SortDirectionDesc:
				desc++
			case SortDirectionAsc:
				asc++
			}
		}, func() CursorType {
			oneDirectionOnly := asc == 0 || desc == 0
			if oneDirectionOnly {
				return CursorTypeCursor
			}
			return CursorTypeOffset
		}
}

func HasReversePage(
	cursor *string,
	pagination ConnectionPagination,
	cursorType CursorType,
	countFunc func() (int64, error),
) (bool, error) {
	if cursor != nil {
		if cursorType == CursorTypeCursor {
			reverseCount, err := countFunc()
			if err != nil {
				return false, err
			}
			return reverseCount > 0, nil
		}
		return true, nil
	}
	return false, nil
}

const (
	cursorValueSeparator = ":"
	cursorSliceSeparator = "--##%$(_)$%##--"
)

// CursorTimeLayout is the layout used to write time values into pagination
// cursors: wall-clock at second precision, the textual form SQL databases
// compare correctly against DATETIME/TIMESTAMP columns. Override it (before
// serving requests) if your columns need sub-second precision, e.g.
// "2006-01-02 15:04:05.000000" for DATETIME(6).
var CursorTimeLayout = "2006-01-02 15:04:05"

// TimeSortKeyMatcher optionally marks sort keys whose cursor values are unix
// seconds (integers) that must be compared against a wall-clock SQL column.
// When it returns true for a key, ToCursorValue encodes integer values as
// UTC wall-clock time using CursorTimeLayout instead of the bare number —
// a bare integer would be string-compared against a DATETIME column and
// match every row, so pagination would repeat the first page forever.
//
// It is nil by default, which disables the reinterpretation entirely:
// integers are only ever treated as unix seconds when you opt in. Set it
// once at application startup, e.g. when all your time columns follow the
// created_at/updated_at naming convention:
//
//	boilergql.TimeSortKeyMatcher = boilergql.SuffixTimeSortKeyMatcher("_AT")
//
// The UTC formatting assumes the database column stores UTC wall time (the
// go-sql-driver/mysql default without a loc parameter).
var TimeSortKeyMatcher func(key string) bool

// SuffixTimeSortKeyMatcher returns a TimeSortKeyMatcher that reports true
// when the sort key ends in any of the given suffixes. Sort keys are the
// uppercased enum values gqlgen generates (e.g. "CREATED_AT"), so the
// conventional configuration is SuffixTimeSortKeyMatcher("_AT").
func SuffixTimeSortKeyMatcher(suffixes ...string) func(key string) bool {
	return func(key string) bool {
		for _, suffix := range suffixes {
			if strings.HasSuffix(key, suffix) {
				return true
			}
		}
		return false
	}
}

// ToCursorValue encodes one sort key/value pair for a pagination cursor.
//
// Pointers are dereferenced first (gqlgen models nullable sort fields as
// *int/*string; %v on a pointer would print the memory address, so the next
// page would filter on "0xc000..." and repeat the same rows forever). A nil
// pointer encodes as an empty value.
//
// time.Time values are formatted with CursorTimeLayout in their own
// location (preserving the wall-clock reading the SQL driver produced)
// instead of %v, whose output — monotonic clock and zone included — never
// round-trips through FromCursorValue into a value comparable against a
// DATETIME/TIMESTAMP column.
//
// Integer values are additionally interpreted as unix seconds when
// TimeSortKeyMatcher is set and matches the key; see its documentation.
func ToCursorValue(k string, v interface{}) string {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return k + cursorValueSeparator
		}
		rv = rv.Elem()
	}
	if rv.IsValid() {
		v = rv.Interface()
	}
	if t, ok := v.(time.Time); ok {
		return k + cursorValueSeparator + t.Format(CursorTimeLayout)
	}
	if TimeSortKeyMatcher != nil && TimeSortKeyMatcher(k) {
		if unix, ok := asUnixSeconds(v); ok {
			return k + cursorValueSeparator + time.Unix(unix, 0).UTC().Format(CursorTimeLayout)
		}
	}
	return k + cursorValueSeparator + fmt.Sprintf("%v", v)
}

func asUnixSeconds(v interface{}) (int64, bool) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int(), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(rv.Uint()), true
	}
	return 0, false
}

func FromCursorValue(cursor string) (string, string) {
	keyValue := strings.SplitN(cursor, cursorValueSeparator, 2)
	if len(keyValue) != 2 {
		return "", ""
	}
	return keyValue[0], keyValue[1]
}

func StringToInterface(v string) interface{} {
	if v == "" {
		return nil
	}
	return v
}

func CursorValuesToString(v []string) string {
	return base64.StdEncoding.EncodeToString([]byte(strings.Join(v, cursorSliceSeparator)))
}

func CursorStringToValues(v string) []string {
	s, _ := base64.StdEncoding.DecodeString(v) //nolint:errcheck
	return strings.Split(string(s), cursorSliceSeparator)
}

func ToOffsetCursor(index int) string {
	return strconv.Itoa(index + 1)
}

func parenthese(v string) string {
	return "(" + v + ")"
}

func GetCursorWhere(comparisonSign ComparisonSign, columns []string, values []interface{}) string {
	return parenthese(strings.Join(columns, ", ")) + " " +
		string(comparisonSign) + " " +
		parenthese(strings.TrimSuffix(strings.Repeat("?,", len(values)), ","))
}

func EdgeLength(pagination ConnectionPagination, length int) int {
	limit := GetLimit(pagination.Forward, pagination.Backward)
	maxLength := limit - 1
	return int(math.Min(float64(length), float64(maxLength)))
}

func BaseConnection(
	pagination ConnectionPagination,
	length int,
	appendEdge func(i int),
) bool {
	limit := GetLimit(pagination.Forward, pagination.Backward)
	maxLength := limit - 1

	switch {
	case pagination.Backward != nil:
		// If the last argument is provided, reverse the order of the results
		for i := length - 1; i >= 0; i-- {
			if i == maxLength {
				continue
			}
			appendEdge(i)
		}

	case pagination.Forward != nil:
		for i := 0; i < length; i++ {
			if i == maxLength {
				break
			}
			appendEdge(i)
		}
	}

	return length == limit
}

func HasNextAndPreviousPage(pagination ConnectionPagination, hasMore bool, hasMoreReversed bool) (bool, bool) {
	switch {
	case pagination.Backward != nil:
		return hasMoreReversed, hasMore
	case pagination.Forward != nil:
		return hasMore, hasMoreReversed
	default:
		return false, false
	}
}

func NewBackwardPagination(last int, before *string) ConnectionPagination {
	return ConnectionPagination{
		Backward: &ConnectionBackwardPagination{
			Last:   last,
			Before: before,
		},
	}
}

func NewForwardPagination(first int, after *string) ConnectionPagination {
	return ConnectionPagination{
		Forward: &ConnectionForwardPagination{
			First: first,
			After: after,
		},
	}
}
