package geojsonDikstra

import "sync"

type VertexQueue struct {
	Vertices []Vertex
	Lock     sync.RWMutex
}

func (vq *VertexQueue) Size() int {
	vq.Lock.RLock()
	count := len(vq.Vertices)
	vq.Lock.RUnlock()
	return count
}

func (vq *VertexQueue) Empty() bool {
	vq.Lock.RLock()
	empty := len(vq.Vertices) == 0
	vq.Lock.RUnlock()
	return empty
}

func (vq *VertexQueue) Take() Vertex {
	vq.Lock.Lock()

	// remove first one from remaining data
	first := vq.Vertices[0]
	vq.Vertices = vq.Vertices[1:len(vq.Vertices)]

	vq.Lock.Unlock()
	return first
}

func (vq *VertexQueue) Add(vertex Vertex) {
	vq.Lock.Lock()

	// no existing data we don't need to prioritize
	if len(vq.Vertices) == 0 {
		vq.Vertices = append(vq.Vertices, vertex)
		vq.Lock.Unlock()
		return
	}

	// weight queue to take shortest distance first
	added := false
	for key, vert := range vq.Vertices {
		if vertex.Distance < vert.Distance {
			if key == 0 {
				// add to the beginning
				vq.Vertices = append([]Vertex{vertex}, vq.Vertices...)
				added = true
				break
			} else {
				// slice into array
				vq.Vertices = append(vq.Vertices[:key+1], vq.Vertices[key:]...)
				vq.Vertices[key] = vertex
				added = true
				break
			}
		}
	}

	// add to the end if it wasn't put anywhere else
	if !added {
		vq.Vertices = append(vq.Vertices, vertex)
	}

	vq.Lock.Unlock()
}
