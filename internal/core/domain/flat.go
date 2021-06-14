package domain

import (
	"time"

	"github.com/mendezdev/tgo_flattener/apierrors"
)

type Flat struct {
	Graph *Graph
}

func NewFlat() Flat {
	return Flat{
		Graph: NewDirectedGraph(),
	}
}

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

/* FUNCTIONS */

// FlatArray it receive an input array an recursive will find
// the max depth of the array and will build a Graph. This info is wrapped
// in a FlatInfo
func (f Flat) FlatArray(input []interface{}) (FlatInfo, apierrors.RestErr) {
	var node int
	var maxDepth int
	f.Graph.AddVertex(node, nil)

	// this callback func  will create the nodes and added the connections
	// to build the Graph. Also will track the max depth
	cb := func(father int, depth int, val interface{}) (int, apierrors.RestErr) {
		if depth > maxDepth {
			maxDepth = depth
		}

		var data interface{}
		if _, ok := val.([]interface{}); !ok {
			switch val.(type) {
			case map[string]interface{}:
				return 0, apierrors.NewBadRequestError("object is not a valid value inside an array")
			}
			data = val
		}

		// every this cb is execute, it means that it is in a node value inside the array
		// so add a vertex (node) to the Graph and the connection with father-son relation
		// e.g: after added 1 to node, this is the father for the next iteration and the "father"
		// is the node in the before iteration
		node++
		f.Graph.AddVertex(node, data)
		if err := f.Graph.AddEdge(father, node); err != nil {
			return 0, apierrors.NewInternalServerError(err.Error())
		}

		return node, nil
	}

	// start from zero node by default
	if err := buildGraphRecursive(input, 0, 0, cb); err != nil {
		return FlatInfo{}, err
	}

	return FlatInfo{
		Graph:          f.Graph,
		VertexSecuence: f.Graph.GetVertexSecuence(),
		MaxDepth:       maxDepth,
		ProcessedAt:    time.Now().UTC(),
	}, nil
}

func (f Flat) BuildGraphFromVertexSecuence(vertexSecuence []VertexSecuence) (*Graph, apierrors.RestErr) {
	g := NewDirectedGraph()

	// creating all the vertex's
	for _, vs := range vertexSecuence {
		parsedValue, err := vs.DataInfo.ToInterface()
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

/******** PRIVATE FUNC ********/

func buildGraphRecursive(data []interface{}, father int, depth int, cb func(int, int, interface{}) (int, apierrors.RestErr)) apierrors.RestErr {
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
