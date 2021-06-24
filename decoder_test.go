package decoder

import (
	"bytes"
	"github.com/golang/geo/s2"
	"math"
	"testing"
)
func floatEq(a, b float64) bool {
	return math.Abs(a - b) <= 1e-9
}

func TestLinestring(t *testing.T) {
	d := New()
	poly, err := d.ParseLinestring(bytes.NewReader([]byte("LINESTRING (30 10, 10 30, 40 40)")))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(*poly) != 3 {
		t.Errorf("Polyine should have 3 points, found %d", len(*poly))
	}
	first := s2.LatLngFromPoint((*poly)[0])
	if !floatEq(first.Lng.Degrees(), 30.0) || !floatEq(first.Lat.Degrees(), 10) {
		t.Errorf("Incorrect first point, expected (30, 10), but found (%f, %f)", first.Lng.Degrees(), first.Lat.Degrees())
	}
}
