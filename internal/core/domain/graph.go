package domain

import (
	"errors"

	"github.com/mendezdev/tgo_flattener/pkg/convert"
)

// Graph contains all the information about the array
// knows how to flat and unflat an array
type Graph struct {
	Vertices map[int]*Vertex
	directed bool
}

// Vertex represents the nodes in the Graph and the connection between each other
type Vertex struct {
	Key      int
	Value    interface{}
	Vertices map[int]*Vertex
}

// VertexSecuence contains all the nodes information to restore the flatted array
// with the EdgeSecuence
type VertexSecuence struct {
	Key      int      `bson:"key"`
	DataInfo DataInfo `bson:"data"`
	Edges    []int    `json:"edges"`
}

// EdgeSecuence represents the connections between the nodes
type EdgeSecuence struct {
	From int `bson:"from"`
	To   int `bson:"to"`
}

/* CONSTRUCTORS */

func NewVertex(key int, value interface{}) *Vertex {
	return &Vertex{
		Key:      key,
		Value:    value,
		Vertices: map[int]*Vertex{},
	}
}

func NewDirectedGraph() *Graph {
	return &Graph{
		Vertices: map[int]*Vertex{},
		directed: true,
	}
}

/* METHODS */

// ToArray will build the array with the information in the Graph
func (g *Graph) ToArray() []interface{} {
	res := make([]interface{}, 0)
	for _, v := range g.Vertices[0].Vertices {
		vtxRes := v.ToArray()
		res = append(res, vtxRes)
	}
	return res
}

// ToArray is called by Graph to build the array
func (v *Vertex) ToArray() interface{} {
	if len(v.Vertices) <= 0 {
		return v.Value
	}

	res := make([]interface{}, 0)
	for _, neighbor := range v.Vertices {
		val := neighbor.ToArray()
		res = append(res, val)
	}

	return res
}

// ToFlat will return the flatted array with the Graph information
func (g *Graph) ToFlat() []interface{} {
	res := make([]interface{}, 0)
	for _, v := range g.Vertices[0].Vertices {
		vtxRes := v.ToFlat()
		d, ok := vtxRes.([]interface{})
		if ok {
			res = append(res, d...)
		} else {
			res = append(res, vtxRes)
		}
	}
	return res
}

// ToFlat is called by Graph to build the flaated array
func (v *Vertex) ToFlat() interface{} {
	if len(v.Vertices) <= 0 {
		return v.Value
	}

	res := make([]interface{}, 0)
	for _, neighbor := range v.Vertices {
		val := neighbor.ToFlat()
		d, ok := val.([]interface{})
		if ok {
			res = append(res, d...)
		} else {
			res = append(res, val)
		}
	}
	return res
}

// GetVertexSecuence build the secuence necesary to be saved in db to be use
// to rebuild the Graph and the array
func (g *Graph) GetVertexSecuence() []VertexSecuence {
	vtxSecuence := make([]VertexSecuence, 0)
	for _, v := range g.Vertices {
		vtxSecuence = append(vtxSecuence, v.GetVertexSecuence())
	}
	return vtxSecuence
}

// GetVertexSecuence is called by Graph
func (v *Vertex) GetVertexSecuence() VertexSecuence {
	var dt, dv string
	var err error
	if v.Value != nil {
		dt, dv, err = convert.GetTypeAndValueStringFromInterface(v.Value)
	}
	// TODO: return apierrors?
	if err != nil {
		panic(err)
	}
	vs := VertexSecuence{
		Key:      v.Key,
		DataInfo: DataInfo{DataType: dt, DataValue: dv},
		Edges:    make([]int, 0),
	}
	for _, neighbor := range v.Vertices {
		vs.Edges = append(vs.Edges, neighbor.Key)
	}
	return vs
}

// AddVertex creates a new Vertex and added to the Graph
func (g *Graph) AddVertex(key int, val interface{}) {
	v := NewVertex(key, val)
	g.Vertices[key] = v
}

// AddEdge connect to Vertex
func (g *Graph) AddEdge(k1, k2 int) error {
	v1 := g.Vertices[k1]
	v2 := g.Vertices[k2]

	// TODO: return apierrors ?
	if v1 == nil || v2 == nil {
		return errors.New("not all vertices exists")
	}

	// is already connected
	if _, ok := v1.Vertices[v2.Key]; ok {
		return nil
	}

	// check if is undirected
	v1.Vertices[v2.Key] = v2
	if !g.directed && v1.Key != v2.Key {
		v2.Vertices[v1.Key] = v1
	}

	// add the vertices to the graph vertex map
	g.Vertices[v1.Key] = v1
	g.Vertices[v2.Key] = v2
	return nil
}
