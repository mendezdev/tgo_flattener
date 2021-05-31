package flattener

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/mendezdev/tgo_flattener/apierrors"
)

// FlatResponse represents the client response for POST /flats
type FlatResponse struct {
	MaxDepth int           `json:"max_depth"`
	Data     []interface{} `json:"flatted_data"`
}

// FlatResponse represents the client response for GET /flats
type FlatInfoResponse struct {
	ID          string        `json:"id"`
	ProcessedAt time.Time     `json:"processed_at"`
	Unflatted   []interface{} `json:"unflatted"`
	Flatted     []interface{} `json:"flatted"`
}

// FlatInfo represents the structure to be saved in the db
type FlatInfo struct {
	ID             string           `json:"id" bson:"_id,omitempty"`
	Graph          *Graph           `bson:"-"`
	VertexSecuence []VertexSecuence `bson:"vertex_secuence"`
	MaxDepth       int              `bson:"max_depth"`
	ProcessedAt    time.Time        `bson:"processed_at"`
}

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

// DataInfo represents the information of the value of the array node
type DataInfo struct {
	DataType  string `bson:"type"`
	DataValue string `bson:"value"`
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
		dt, dv, err = getTypeAndValueStringFromInterface(v.Value)
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

// toInterface rebuild the original value of in the array
func (di DataInfo) toInterface() (interface{}, error) {
	var convertedValue interface{}
	var err error

	switch di.DataType {
	case "float64":
		convertedValue, err = strconv.ParseFloat(di.DataValue, 64)
	case "bool":
		convertedValue, err = strconv.ParseBool(di.DataValue)
	default:
		convertedValue = di.DataValue
	}
	if err != nil {
		return nil, fmt.Errorf("error parsing flat_data: %s", err.Error())
	}
	return convertedValue, nil
}

/* FUNCTIONS */

// FlatArray it receive an input array an recursive will find
// the max depth of the array and will build a Graph. This info is wrapped
// in a FlatInfo
func FlatArray(input []interface{}) (FlatInfo, apierrors.RestErr) {
	g := NewDirectedGraph()

	var node int
	var maxDepth int
	g.AddVertex(node, nil)

	// this callback func  will create the nodes and added the connections
	// to build the Graph. Also will track the max depth
	cb := func(father int, depth int, val interface{}) (int, error) {
		if depth > maxDepth {
			maxDepth = depth
		}

		var data interface{}
		if _, ok := val.([]interface{}); !ok {
			data = val
		}

		// every this cb is execute, it means that it is in a node value inside the array
		// so add a vertex (node) to the Graph and the connection with father-son relation
		// e.g: after added 1 to node, this is the father for the next iteration and the "father"
		// is the node in the before iteration
		node++
		g.AddVertex(node, data)
		if err := g.AddEdge(father, node); err != nil {
			return 0, err
		}

		return node, nil
	}

	// start from zero node by default
	if err := buildGraphRecursive(input, 0, 0, cb); err != nil {
		return FlatInfo{}, apierrors.NewInternalServerError(err.Error())
	}

	return FlatInfo{
		Graph:          g,
		VertexSecuence: g.GetVertexSecuence(),
		MaxDepth:       maxDepth,
		ProcessedAt:    time.Now().UTC(),
	}, nil
}

func buildGraphRecursive(data []interface{}, father int, depth int, cb func(int, int, interface{}) (int, error)) error {
	for _, v := range data {
		var d int

		// if it is an array, add one to depth
		// call the function again to go more in depth and pass the info
		// to the next iteration
		parsed, ok := v.([]interface{})
		if ok {
			d = depth + 1
		}

		// current will be the father for the next iteration and actual father is for the current
		current, err := cb(father, d, v)
		if err != nil {
			return err
		}
		if ok {
			if err := buildGraphRecursive(parsed, current, d, cb); err != nil {
				return err
			}
		}
	}
	return nil
}

func BuildGraphFromVertexSecuence(vertexSecuence []VertexSecuence) (*Graph, apierrors.RestErr) {
	g := NewDirectedGraph()

	// creating all the vertex's
	for _, vs := range vertexSecuence {
		parsedValue, err := vs.DataInfo.toInterface()
		if err != nil {
			return nil, apierrors.NewInternalServerError("error parsing data_info")
		}
		g.AddVertex(vs.Key, parsedValue)
	}

	// creating all the edge connections
	for _, vs := range vertexSecuence {
		for _, e := range vs.Edges {
			if err := g.AddEdge(vs.Key, e); err != nil {
				return nil, apierrors.NewInternalServerError(err.Error())
			}
		}
	}

	return g, nil
}

func getTypeAndValueStringFromInterface(val interface{}) (dt string, dv string, err error) {
	if val == nil {
		err = errors.New("cannot get type and value from nil interface")
		return
	}
	dt = fmt.Sprintf("%T", val)
	dv = fmt.Sprintf("%v", val)
	return
}
