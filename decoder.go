package decoder

import (
    "errors"
    "fmt"
    "github.com/IvanZagoskin/wkt/geometry"
    "github.com/IvanZagoskin/wkt/parser"
    "github.com/golang/geo/s2"
    "io"
)

type Decoder struct {
    parser *parser.Parser
}

func GeomString(g geometry.Type) string {
    switch g {
    case geometry.UndefinedGT: return "Undefined"
    case geometry.PointGT: return "Point"
    case geometry.MultyPointGT: return "Multipoint"
    case geometry.LineStringGT: return "Linestring"
    case geometry.CircularStringGT: return "CircularLinestring"
    case geometry.MultiLineStringGT: return "MultiLineString"
    case geometry.PolygonGT: return "Polygon"
    case geometry.MultiPolygonGT: return "MultiPolygon"
    default:
        return fmt.Sprintf("%d", int(g))
    }
}

func asS2LatLngs(points []*geometry.Point) []s2.LatLng {
    xs := make([]s2.LatLng, len(points))
    for i, p := range points {
        xs[i] = s2.LatLngFromDegrees(p.X, p.Y)
    }
    return xs
}

func asS2Points(points []*geometry.Point) []s2.Point {
    xs := make([]s2.Point, len(points))
    for i, p := range points {
        xs[i] = s2.PointFromLatLng(s2.LatLngFromDegrees(p.X, p.Y))
    }
    return xs
}

func asS2Polygon(lines []*geometry.LineString) *s2.Polygon {
    loops := make([]*s2.Loop, len(lines))
    for i, ls := range lines {
        loops[i] = s2.LoopFromPoints(asS2Points(ls.Points))
    }
    return s2.PolygonFromLoops(loops)
}

func (d *Decoder) ParsePoint(r io.Reader) (s2.Point, error) {
    g, err := d.ParseWKT(r)
    if err != nil {
        return s2.Point{}, err
    }
    if x, ok := g.(s2.Point); ok {
        return x, nil
    }
    return s2.Point{}, errors.New("cast error")
}

func (d *Decoder) ParseLinestring(r io.Reader) (*s2.Polyline, error) {
    g, err := d.ParseWKT(r)
    if err != nil {
        return nil, err
    }
    if x, ok := g.(*s2.Polyline); ok {
        return x, nil
    }
    return nil, errors.New("cast error")
}

func (d *Decoder) ParsePolygon(r io.Reader) (*s2.Polygon, error) {
    g, err := d.ParseWKT(r)
    if err != nil {
        return nil, err
    }
    if x, ok := g.(*s2.Polygon); ok {
        return x, nil
    }
    return nil, errors.New("cast error")
}

func (d *Decoder) ParseMultiPoint(r io.Reader) ([]s2.Point, error) {
    g, err := d.ParseWKT(r)
    if err != nil {
        return nil, err
    }
    if x, ok := g.([]s2.Point); ok {
        return x, nil
    }
    return nil, errors.New("cast error")
}

func (d *Decoder) ParseMultiLinestring(r io.Reader) ([]s2.Polyline, error) {
    g, err := d.ParseWKT(r)
    if err != nil {
        return nil, err
    }
    if x, ok := g.([]s2.Polyline); ok {
        return x, nil
    }
    return nil, errors.New("cast error")
}

func (d *Decoder) ParseMultiPolygon(r io.Reader) ([]s2.Polygon, error) {
    g, err := d.ParseWKT(r)
    if err != nil {
        return nil, err
    }
    if x, ok := g.([]s2.Polygon); ok {
        return x, nil
    }
    return nil, errors.New("cast error")
}

func (d *Decoder) ParseWKT(r io.Reader) (interface{}, error) {
    g, err := d.parser.ParseWKT(r)
    if err != nil {
        return nil, err
    }
    switch geom := g.(type) {
    case *geometry.Point:
        return s2.PointFromLatLng(s2.LatLngFromDegrees(geom.X, geom.Y)), nil
    case *geometry.LineString:
        return s2.PolylineFromLatLngs(asS2LatLngs(geom.Points)), nil
    case *geometry.Polygon:
        return asS2Polygon(geom.LineStrings), nil
    case *geometry.MultiPoint:
        return asS2Points(geom.Points), nil
    case *geometry.MultiLineString:
        lines := make([]s2.Polyline, len(geom.Lines))
        for i, ls := range geom.Lines {
            lines[i] = *s2.PolylineFromLatLngs(asS2LatLngs(ls.Points))
        }
        return lines, nil
    case *geometry.MultiPolygon:
        polygons := make([]s2.Polygon, len(geom.Polygons))
        for i, p := range geom.Polygons {
            polygons[i] = *asS2Polygon(p.LineStrings)
        }
        return polygons, nil
    default:
        return nil, fmt.Errorf("unimplemented geometry type: %s", GeomString(g.GetGeometryType()))
    }
}

func New() Decoder {
    return Decoder{
        parser: parser.New(),
    }
}
