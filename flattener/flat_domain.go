package flattener

import (
	"errors"
	"fmt"
)

type FlatResponse struct {
	MaxDepth int           `json:"max_depth"`
	Data     []interface{} `json:"flatted_data"`
}

type FlatInfoResponse struct {
	ID          string        `json:"id"`
	DateCreated string        `json:"date_created"`
	Unflatted   []interface{} `json:"unflatted"`
	Flatted     []interface{} `json:"flatted"`
}

type FlatInfo struct {
	ID               string              `json:"id" bson:"_id,omitempty"`
	StructureInfo    []FlatStructureInfo `bson:"structure"`
	StructureFlatted []FlatData          `bson:"structure_flatted"`
	DateCreated      string              `bson:"date_created"`
}

type FlatStructureInfo struct {
	Level int        `bson:"level"`
	Data  []FlatData `bson:"data"`
}

type FlatData struct {
	DataType  string `bson:"type"`
	DataValue string `bson:"value"`
}

// Graph
type Graph struct {
	Vertices map[int]*Vertex
	directed bool
}

type Vertex struct {
	Key      int
	Value    interface{}
	Vertices map[int]*Vertex
}

type GraphSecuence struct {
	VertexSecuence []VertexSecuence `bson:"vertex_secuence"`
	EdgeSecuence   []EdgeSecuence   `bson:"edge_secuence"`
}

type VertexSecuence struct {
	Key      int      `bson:"key"`
	DataInfo DataInfo `bson:"data"`
}

type EdgeSecuence struct {
	From int `bson:"from"`
	To   int `bson:"to"`
}

type DataInfo struct {
	DataType  string `bson:"type"`
	DataValue string `bson:"value"`
}

// CONSTRUCTORS
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

func NewGraphSecuence() GraphSecuence {
	return GraphSecuence{
		VertexSecuence: make([]VertexSecuence, 0),
		EdgeSecuence:   make([]EdgeSecuence, 0),
	}
}

// METHODS
func (g *Graph) ToArray() []interface{} {
	res := make([]interface{}, 0)
	for _, v := range g.Vertices[0].Vertices {
		vtxRes := v.ToArray()
		res = append(res, vtxRes)
	}
	return res
}

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

// AddVertex creates a new Vertex and added to the Graph
func (g *Graph) AddVertex(key int, val interface{}) {
	v := NewVertex(key, val)
	g.Vertices[key] = v
}

// AddEdge connect to Vertex
func (g *Graph) AddEdge(k1, k2 int) {
	v1 := g.Vertices[k1]
	v2 := g.Vertices[k2]

	// TODO: return apierrors ?
	if v1 == nil || v2 == nil {
		panic("not all vertices exists")
	}

	// is already connected
	if _, ok := v1.Vertices[v2.Key]; ok {
		return
	}

	// check if is undirected
	v1.Vertices[v2.Key] = v2
	if !g.directed && v1.Key != v2.Key {
		v2.Vertices[v1.Key] = v1
	}

	// add the vertices to the graph vertex map
	g.Vertices[v1.Key] = v1
	g.Vertices[v2.Key] = v2
}

// FUNCTIONS
func BuildGraphFromArray(input []interface{}) *Graph {
	gs := NewGraphSecuence()
	g := NewDirectedGraph()

	var node int
	g.AddVertex(node, nil)

	cb := func(father int, val interface{}) int {
		node++
		var data interface{}

		if _, ok := val.([]interface{}); !ok {
			data = val
		}
		g.AddVertex(node, data)
		dt, dv, err := getTypeAndValueStringFromInterface(data)

		// TODO: return apierrors
		if err != nil {
			panic(err)
		}
		gs.VertexSecuence = append(gs.VertexSecuence,
			VertexSecuence{node, DataInfo{DataType: dt, DataValue: dv}})
		g.AddEdge(father, node)
		gs.EdgeSecuence = append(gs.EdgeSecuence, EdgeSecuence{father, node})
		return node
	}

	// start from zero node by default
	buildGraphRecursive(input, 0, cb)
	return g
}

func buildGraphRecursive(data []interface{}, father int, cb func(int, interface{}) int) {
	for _, v := range data {
		current := cb(father, v)
		parsed, ok := v.([]interface{})
		if ok {
			buildGraphRecursive(parsed, current, cb)
		}
	}
}

func getTypeAndValueStringFromInterface(val interface{}) (string, string, error) {
	var dt, dv string
	if val == nil {
		return dt, dv, errors.New("cannot get type and value from nil interface")
	}
	dt = fmt.Sprintf("%T", val)
	dv = fmt.Sprintf("%v", val)

	return dt, dv, nil
}
