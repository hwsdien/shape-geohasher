package main

import (
	"fmt"
	"sort"
	"encoding/json"

	"github.com/hwsdien/shape-geohasher/geohasher"

)


func main()  {
	coords := [][]float64{
		{ 113.98212432861328, 22.560915397692284 },
		{ 114.02486801147461, 22.560915397692284 },
		{ 114.02486801147461, 22.584850519363435 },
		{ 113.98212432861328, 22.584850519363435 },
		{ 113.98212432861328, 22.560915397692284 },
	}

	g := geohasher.NewGeoHasher()

	polygon := g.NewGeomPolygon(coords)
	geohashList := g.GetGeohashesOfPolygon(polygon, 7)
	sort.Strings(geohashList)
	jsonResult, err := json.Marshal(geohashList)
	if err != nil {
		panic(err)
	}

	fmt.Println("The content of geohash list: ")
	fmt.Println(string(jsonResult))
	fmt.Printf("The length of geohash list : %d", len(geohashList))
}
