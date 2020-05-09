package boilergql

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// GeoPoint is serialized as a simple array, eg [1, 2]
type GeoPoint struct {
	X float64
	Y float64
}

func (p *GeoPoint) UnmarshalGQL(v interface{}) error {
	pointStr, ok := v.(string)
	if !ok {
		return fmt.Errorf("points must be strings")
	}

	parts := strings.Split(pointStr, ",")

	if len(parts) != 2 {
		return fmt.Errorf("points must have 2 parts")
	}

	var err error
	if p.X, err = strconv.ParseFloat(parts[0], 64); err != nil {
		return err
	}
	if p.Y, err = strconv.ParseFloat(parts[1], 64); err != nil {
		return err
	}
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (p GeoPoint) MarshalGQL(w io.Writer) {
	fmt.Fprintf(w, `"%g,%g"`, p.X, p.Y)
}
