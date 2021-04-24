package boilergql

import "strings"

func GetQuestionMarksForColumns(columns []string) string {
	b := new(strings.Builder)
	b.WriteString("(")
	for i := range columns {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString("?")
	}
	b.WriteString(")")
	return b.String()
}
