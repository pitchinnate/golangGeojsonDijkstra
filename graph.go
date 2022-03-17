package geojsonDikstra

import (
	"encoding/json"
	"fmt"
	"github.com/tomchavakis/turf-go/measurement"
	"log"
	"math"
	"strconv"
	"strings"
)

type Position [2]float64

type Geometry struct {
	Coordinates []Position `json:"coordinates"`
	Type        string     `json:"type"`
}

type Feature struct {
	Type       string                 `json:"type"`
	Geometry   Geometry               `json:"geometry"`
	ID         *string                `json:"id"`
	Properties map[string]interface{} `json:"properties"`
}

type FeatureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Edge struct {
	Position1 Position
	Position2 Position
}

type Topology struct {
	Vertices []Position
	Edges    []Edge
}

type Path []Position

type Graph struct {
	Vertices           map[string]map[string]float64
	SimplifiedVertices map[string]map[string]float64
	Topo               *Topology
}

type Vertex struct {
	Key      string
	Distance float64
}

var UnitOfMeasure = "meters"

func (graph *Graph) PrintDebug() {
	log.Println("Topo Vertices: ", graph.Topo.Vertices)

	for index, edge := range graph.Topo.Edges {
		log.Printf("Edge %d: %s -> %s \n", index, edge.Position1.Key(), edge.Position2.Key())
	}

	for key, vert := range graph.Vertices {
		log.Printf("%s Connections\n", key)
		for key2, vert2 := range vert {
			log.Printf("    %s -> %s = %f \n", key, key2, vert2)
		}
	}
}

func (path *Path) Json() {
	jsonString, _ := json.Marshal(path)
	log.Print(string(jsonString))
}

func (graph *Graph) ClosestVertex(position Position) (Position, float64) {
	minimum := math.MaxFloat64
	var closestPosition Position
	for _, vertex := range graph.Topo.Vertices {
		weight, err := measurement.Distance(vertex[0], vertex[1], position[0], position[1], UnitOfMeasure)
		if err == nil && weight < minimum {
			minimum = weight
			closestPosition = vertex
		}
	}
	return closestPosition, minimum
}

func (graph *Graph) ShortestPath(start Position, end Position) (Path, float64, error) {
	dist := make(map[string]float64)
	prev := make(map[string]string)
	visited := make(map[string]bool)
	queue := VertexQueue{}

	// set max distance for each
	for _, pos := range graph.Topo.Vertices {
		dist[pos.Key()] = math.MaxFloat64
	}
	dist[start.Key()] = 0
	queue.Add(Vertex{
		Key:      start.Key(),
		Distance: 0,
	})
	for !queue.Empty() {
		nextVertex := queue.Take()

		// if we have already visited this node we can ignore there was a shorter path
		if visited[nextVertex.Key] {
			continue
		}

		// mark this vertex as visited so we don't visit it again
		visited[nextVertex.Key] = true

		// find connections to this vertex
		connections := graph.Vertices[nextVertex.Key]

		for key, distance := range connections {
			// if we already have been to this connection ignore it
			if visited[key] {
				continue
			}
			// if the distance to the new node is shorter than we have already found use it instead
			newDistance := dist[nextVertex.Key] + distance
			if newDistance < dist[key] {
				dist[key] = newDistance
				prev[key] = nextVertex.Key
				queue.Add(Vertex{
					Key:      key,
					Distance: newDistance,
				})
			}
		}
	}
	finalPath := []string{end.Key()}
	nextVal := prev[end.Key()]
	for nextVal != start.Key() && nextVal != "" {
		finalPath = append(finalPath, nextVal)
		nextVal = prev[nextVal]
	}
	finalPath = append(finalPath, nextVal)

	var path Path

	if len(finalPath) == 0 {
		return path, 0, fmt.Errorf("no route was found")
	}

	for _, stringCords := range finalPath {
		pieces := strings.Split(stringCords, ",")
		var newPos Position
		var err error
		newPos[0], err = strconv.ParseFloat(pieces[0], 64)
		if err != nil {
			return path, 0, fmt.Errorf("error parsing to float %s", pieces[0])
		}
		newPos[1], err = strconv.ParseFloat(pieces[1], 64)
		if err != nil {
			return path, 0, fmt.Errorf("error parsing to float %s", pieces[0])
		}
		path = append(path, newPos)
	}

	return path, dist[end.Key()], nil
}

func (fc *FeatureCollection) ToTopology(precision float64) Topology {
	var topo Topology
	allPositions := make(map[string]Position)

	for _, feature := range fc.Features {
		if feature.Geometry.Type != "LineString" {
			continue
		}
		for _, cord := range feature.Geometry.Coordinates {
			roundedCord := cord.Roundoff(precision)
			allPositions[roundedCord.Key()] = roundedCord
		}
		topo.Edges = append(topo.Edges, Edge{
			Position1: feature.Geometry.Coordinates[0].Roundoff(precision),
			Position2: feature.Geometry.Coordinates[1].Roundoff(precision),
		})
	}

	for _, pos := range allPositions {
		topo.Vertices = append(topo.Vertices, pos)
	}

	return topo
}

func (topo *Topology) Preprocess() Graph {
	var graph Graph
	graph.Vertices = make(map[string]map[string]float64)
	graph.Topo = topo
	for _, edge := range topo.Edges {
		weight, err := measurement.Distance(edge.Position1[0], edge.Position1[1], edge.Position2[0], edge.Position2[1], UnitOfMeasure)
		if err == nil && weight > 0 {
			_, ok := graph.Vertices[edge.Position1.Key()]
			if !ok {
				graph.Vertices[edge.Position1.Key()] = make(map[string]float64)
			}
			_, ok = graph.Vertices[edge.Position2.Key()]
			if !ok {
				graph.Vertices[edge.Position2.Key()] = map[string]float64{}
			}
			graph.Vertices[edge.Position1.Key()][edge.Position2.Key()] = weight
			graph.Vertices[edge.Position2.Key()][edge.Position1.Key()] = weight
		}
	}
	return graph
}

func (position *Position) Roundoff(precision float64) Position {
	var newPosition Position
	newPosition[0] = math.Round(position[0]/precision) * precision
	newPosition[1] = math.Round(position[1]/precision) * precision
	return newPosition
}

func (position *Position) Key() string {
	return fmt.Sprintf("%f,%f", position[0], position[1])
}
