package boilergql

import (
	"fmt"
	"math"
	"strconv"
	"strings"
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
	return dbColumn + " " + direction.String()
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

func HasReversePage(cursor *string, pagination ConnectionPagination, cursorType CursorType, countFunc func() (int64, error)) (bool, error) {
	if cursor != nil {
		if cursorType == CursorTypeCursor {
			reverseCount, err := countFunc()
			if err != nil {
				return false, err
			}
			return reverseCount > 0, nil
		} else {
			return true, nil
		}
	}
	return false, nil
}

const (
	cursorValueSeparator = ":"
	cursorSliceSeparator = "--##%$(_)$%##--"
)

func ToCursorValue(k string, v interface{}) string {
	return k + cursorValueSeparator + fmt.Sprintf("%v", v)
}

func FromCursorValue(cursor string) (string, string) {
	keyValue := strings.SplitN(cursor, cursorValueSeparator, 2)
	if len(keyValue) != 2 {
		return "", ""
	}
	return keyValue[0], keyValue[1]
}

func CursorValuesToString(v []string) string {
	return strings.Join(v, cursorSliceSeparator)
}

func CursorStringToValues(v string) []string {
	return strings.Split(v, cursorSliceSeparator)
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
