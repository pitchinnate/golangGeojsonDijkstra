package geojsonDikstra

func (fc *FeatureCollection) FindPath(start Position, end Position, precision float64) (Path, float64, error) {
	topology := fc.ToTopology(precision)
	graph := topology.Preprocess()

	// Find the closest points on our graph to the requested start and end
	closestStart, _ := graph.ClosestVertex(start)
	closestEnd, _ := graph.ClosestVertex(end)

	// Find the shortest path from start to end
	path, distance, err := graph.ShortestPath(closestStart, closestEnd)

	return path, distance, err
}
