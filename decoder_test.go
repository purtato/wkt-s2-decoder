package decoder

import (
	"bytes"
	"github.com/golang/geo/s2"
	"math"
	"testing"
)

func floatEq(a, b float64) bool {
	return math.Abs(a-b) <= 1e-9
}

func TestPoint(t *testing.T) {
	d := New()
	p, err := d.ParsePoint(bytes.NewReader([]byte("POINT(138.5199401149407 -34.97629388791155)")))
	if err != nil {
		t.Fatalf(err.Error())
	}
	lng := s2.LatLngFromPoint(p).Lng.Degrees()
	lat := s2.LatLngFromPoint(p).Lat.Degrees()
	if !floatEq(lng, 138.5199401149407) {
		t.Errorf("Incorrect lng, expected 138.5199401149407, but found %f", lng)
	}
	if !floatEq(lat, -34.97629388791155) {
		t.Errorf("Incorrect lat, expected -34.97629388791155, but found %f", lat)
	}
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

func TestSanity(t *testing.T) {
	testCases := []string{
		"LINESTRING (151.23941 -33.82436, 151.239761 -33.824387, 151.239913 -33.824402, 151.240036 -33.824418, 151.240097 -33.824425, 151.240158 -33.824433, 151.240417 -33.824463, 151.240692 -33.824482, 151.240798 -33.824471, 151.240921 -33.824437, 151.241012 -33.824387, 151.24118 -33.824276, 151.24147 -33.824055, 151.241958 -33.823746)",
		"LINESTRING (151.058837 -33.95028, 151.059005 -33.950257, 151.059005 -33.95023, 151.05902 -33.950207, 151.059036 -33.950192, 151.059036 -33.950192, 151.059051 -33.950184, 151.059066 -33.950173, 151.059097 -33.950165, 151.059127 -33.950165, 151.059158 -33.950177, 151.059188 -33.950196, 151.059204 -33.950223, 151.059219 -33.950249, 151.059326 -33.95023, 151.059326 -33.95023, 151.059631 -33.950177)",
		"LINESTRING (151.238677 -33.824341, 151.238922 -33.824341, 151.239242 -33.824345, 151.239761 -33.824387, 151.239913 -33.824402, 151.240036 -33.824418, 151.240097 -33.824425, 151.240158 -33.824433, 151.240417 -33.824463, 151.240692 -33.824482, 151.240798 -33.824471, 151.240921 -33.824437, 151.241012 -33.824387, 151.24118 -33.824276, 151.24147 -33.824055, 151.241836 -33.823826)",
		"LINESTRING (151.238769 -33.824341, 151.238922 -33.824341, 151.239242 -33.824345, 151.239761 -33.824387, 151.239913 -33.824402, 151.240036 -33.824418, 151.240097 -33.824425, 151.240158 -33.824433, 151.240417 -33.824463, 151.240692 -33.824482, 151.240798 -33.824471, 151.240921 -33.824437, 151.240997 -33.824566, 151.241058 -33.824616, 151.241134 -33.824643, 151.241226 -33.824681, 151.241561 -33.824841, 151.241897 -33.82499)",
	}
	d := New()
	for _, tc := range testCases {
		_, err := d.ParseLinestring(bytes.NewReader([]byte(tc)))
		if err != nil {
			t.Errorf(err.Error())
		}
	}
}

func TestNil(t *testing.T) {
	d := New()
	_, err := d.ParseLinestring(nil)
	if err == nil {
		t.Errorf("Missing err for nil input")
	}
}

func TestEmpty(t *testing.T) {
	d := New()
	_, err := d.ParseLinestring(bytes.NewReader([]byte("")))
	if err == nil {
		t.Errorf("Missing err for empty input")
	}
}
