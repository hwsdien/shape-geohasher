package geohasher

import (
	"github.com/hwsdien/polyclip-go"
	"github.com/mmcloughlin/geohash"
	geom "github.com/twpayne/go-geom"
	xy "github.com/twpayne/go-geom/xy"
)

type Point struct {
	X float64
	Y float64
}

type GeoHasher struct {}

func NewGeoHasher() *GeoHasher {
	return &GeoHasher{}
}


func (g *GeoHasher) NewGeomPolygon(coords [][]float64) *geom.Polygon {
	polygonCoords := make([][]geom.Coord, 0)
	polygonTempCoords := make([]geom.Coord, 0)
	for _, location := range coords {
		polygonTempCoords = append(polygonTempCoords, location)
	}
	polygonCoords = append(polygonCoords, polygonTempCoords)
	return geom.NewPolygon(geom.XY).MustSetCoords(polygonCoords)
}

func (g *GeoHasher) getPolygon(coords [][]float64) polyclip.Polygon {
	contour := make(polyclip.Contour, 0)

	for _, location := range coords{
		contour = append(contour, polyclip.Point{location[0], location[1]})
	}

	return polyclip.Polygon{contour}
}

func (g *GeoHasher) checkIntersection(p1, p2 polyclip.Polygon) bool {
	result := p1.Construct(polyclip.INTERSECTION, p2)
	if result.NumVertices() > 0 {
		return true
	}
	return false
}


func (g *GeoHasher) getCentroid(p *geom.Polygon) (Point, error) {
	point := Point{}
	centroid, err := xy.Centroid(p)
	if err != nil {
		return point, err
	}

	point.X = centroid.X()
	point.Y = centroid.Y()
	return point, nil
}

func (g *GeoHasher) getPolygonByGeohash(hash string) polyclip.Polygon {
	box := geohash.BoundingBox(hash)
	northWest := []float64{box.MinLng, box.MaxLat}
	southWest := []float64{box.MinLng, box.MinLat}
	southEast := []float64{box.MaxLng, box.MinLat}
	northEast := []float64{box.MaxLng, box.MaxLat}

	coords := [][]float64{northWest, southWest, southEast, northEast, northWest}

	return g.getPolygon(coords)
}


func (g *GeoHasher) GetGeohashesOfPolygon(p *geom.Polygon, precision uint) []string{
	coordList := make([][]float64, 0)
	for _, coord := range p.Coords()[0] {
		coordList = append(coordList, []float64{coord.X(), coord.Y()})
	}

	subject := g.getPolygon(coordList)
	centroid, _ := g.getCentroid(p)
	geohashCode := geohash.EncodeWithPrecision(centroid.Y, centroid.X, precision)

	geohashList := make([]string, 0)
	geohashList = append(geohashList, geohashCode)
	geohashHandledMap := make(map[string]bool, 0)
	geohashMap := make(map[string]bool, 0)


	for len(geohashList) > 0 {
		currentGeohash := geohashList[0]

		geohashList = geohashList[1:]
		currentPolygon := g.getPolygonByGeohash(currentGeohash)
		geohashMap[currentGeohash] = true

		_, ok := geohashHandledMap[currentGeohash]
		if !ok {
			if g.checkIntersection(subject, currentPolygon) {
				geohashHandledMap[currentGeohash] = true

				for _, neighbor := range geohash.Neighbors(currentGeohash) {
					_, ok := geohashMap[neighbor]
					if !ok {
						geohashMap[neighbor] = true
						geohashList = append(geohashList, neighbor)
					}
				}
			}
		}
	}

	geohashes := make([]string , 0)
	for k := range geohashHandledMap {
		geohashes = append(geohashes, k)
	}

	return geohashes
}



