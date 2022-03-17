Golang GeoJson Dijkstra
=================

Golang GeoJson Dijkstra utilizes GeoJson Feature Collections to find the best the shortest route between two points.

Installation
------------
Use go get.

	go get github.com/pitchinnate/golangGeojsonDijkstra

Then import the validator package into your own code.

	import "github.com/pitchinnate/golangGeojsonDijkstra"

Requirements 
------
Must pass in valid GeoJson `FeatureCollection` that has a collection of GeoJson `Feature`'s. Everything but
Features that have a `geometry.type` equal to `LineString` are ignored. Also, currently it is required that all 
`LineString` features **must only have two coordinates**. This may be improved later.

Usage
------
Pass into coordinate arrays and a precision. Coordinate arrays should be in the following format `[longitude, latitude]`
as that corresponds to `[x,y]`. **Precision** is how accurate it rounds the latitude and longitude coordinates to. So if
you pass in a coordinate of `-84.396535122` with a precision of `0.00001` it will round that off to `-84.39654`. This is
important when connections between lines are made.
```
featureCollection.FindPath(Position, Position, precision)
```

Example
------
```go
featureCollection := FeatureCollection{
    Type: "FeatureCollection",
    Features: []Feature{
        Feature{Type: "Feature", Geometry: Geometry{Type: "LineString", Coordinates: []Position{Position{-84.396535, 33.792487}, Position{-84.396535, 33.792578}}}},
        Feature{Type: "Feature", Geometry: Geometry{Type: "LineString", Coordinates: []Position{Position{-84.396536, 33.79281}, Position{-84.396535, 33.792578}}}},
        Feature{Type: "Feature", Geometry: Geometry{Type: "LineString", Coordinates: []Position{Position{-84.396536, 33.79281}, Position{-84.396536, 33.792908}}}},
        Feature{Type: "Feature", Geometry: Geometry{Type: "LineString", Coordinates: []Position{Position{-84.396863, 33.792908}, Position{-84.396536, 33.792908}}}},
        Feature{Type: "Feature", Geometry: Geometry{Type: "LineString", Coordinates: []Position{Position{-84.396537, 33.792996}, Position{-84.396536, 33.792908}}}},
        Feature{Type: "Feature", Geometry: Geometry{Type: "LineString", Coordinates: []Position{Position{-84.396213, 33.792908}, Position{-84.396536, 33.792908}}}},
    },
}
path, distance, err := featureCollection.FindPath(Position{-84.396863, 33.792908}, Position{-84.396535, 33.792578}, 0.00001)

// path should come back as an array of coordinates
// [[-84.39654,33.79258],[-84.39654,33.79281],[-84.39654,33.79291],[-84.39686,33.79291]]
// distance should come back as a float64 measure of how far you traveled (in meters)
```

License
-------
Distributed under GNU GENERAL PUBLIC LICENSE, please see license file within the code for more details.
