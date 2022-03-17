package geojsonDikstra

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func ReadInputFile(path string) []byte {
	dir, _ := os.Getwd()
	filePath := fmt.Sprintf("%s/sample/%s", dir, path)
	contents, _ := ioutil.ReadFile(filePath)
	return contents
}

func TestFindPath(t *testing.T) {
	t.Run("Test FindPath Small Dataset", func(t *testing.T) {
		var fc FeatureCollection
		text := ReadInputFile("sample.json")
		json.Unmarshal(text, &fc)
		path, distance, err := fc.FindPath(Position{-84.396863, 33.792908}, Position{-84.396535, 33.792578}, 0.00001)
		if err != nil {
			t.Errorf("returned an error %s", err)
		}
		if len(path) != 4 {
			t.Errorf("Did not find best route, route had %d coords, should be 3", len(path))
		}
		if distance < 66 || distance > 67 {
			t.Errorf("Distance wasn't correct, had %f meters, should be ~66.265", distance)
		}
	})
	t.Run("Test FindPath Large Dataset", func(t *testing.T) {
		var fc FeatureCollection
		text := ReadInputFile("sample_large.json")
		json.Unmarshal(text, &fc)
		path, distance, err := fc.FindPath(Position{-84.397252, 33.792997}, Position{-84.395111, 33.791666}, 0.00001)
		if err != nil {
			t.Errorf("returned an error %s", err)
		}
		if len(path) != 27 {
			t.Errorf("Did not find best route, route had %d coords, should be 27", len(path))
		}
		if distance < 365 || distance > 366 {
			t.Errorf("Distance wasn't correct, had %f meters, should be ~66.265", distance)
		}
	})
	t.Run("Test FindPath Large Dataset", func(t *testing.T) {
		fc := FeatureCollection{
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
		path, distance, err := fc.FindPath(Position{-84.396863, 33.792908}, Position{-84.396535, 33.792578}, 0.00001)
		if err != nil {
			t.Errorf("returned an error %s", err)
		}
		if len(path) != 4 {
			t.Errorf("Did not find best route, route had %d coords, should be 4", len(path))
		}
		if distance < 66 || distance > 67 {
			t.Errorf("Distance wasn't correct, had %f meters, should be ~66.265", distance)
		}
	})
}
