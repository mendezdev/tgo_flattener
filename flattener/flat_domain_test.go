package flattener

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGraph(t *testing.T) {
	g := NewDirectedGraph()
	assert.NotNil(t, g)
	assert.NotNil(t, g.Vertices)
	assert.True(t, g.directed)
}

func TestNewVertex(t *testing.T) {
	vtx := NewVertex(1, "some_value")

	assert.NotNil(t, vtx)
	assert.Equal(t, vtx.Key, 1)
	assert.Equal(t, vtx.Value, "some_value")
	assert.NotNil(t, vtx.Vertices)
}

func TestVertexToArraySingleValue(t *testing.T) {
	mockValue := "some_value"
	vtx := NewVertex(1, mockValue)
	result := vtx.ToArray()

	assert.NotNil(t, result)
	assert.IsType(t, mockValue, result)
}

func TestVertexToArrayInvolvedInArrayValue(t *testing.T) {
	mockValue := "first_lv2"
	vtxLvl1 := NewVertex(1, "first_lvl")
	vtxLvl2 := NewVertex(2, mockValue)
	vtxLvl1.Vertices[vtxLvl1.Key] = vtxLvl2
	result := vtxLvl1.ToArray()

	assert.NotNil(t, result)

	resultType, ok := result.([]interface{})
	assert.True(t, ok)
	assert.NotNil(t, resultType)
	assert.Len(t, resultType, 1)
	for _, v := range resultType {
		assert.Equal(t, mockValue, v)
	}
}

func TestGraphToArray(t *testing.T) {
	g := NewDirectedGraph()

	// re-creating the following array: [["value2","value3"]]
	g.AddVertex(0, nil)
	g.AddVertex(1, nil)
	g.AddVertex(2, "value2")
	g.AddVertex(3, "value3")

	err := g.AddEdge(0, 1)
	assert.Nil(t, err)

	err = g.AddEdge(1, 2)
	assert.Nil(t, err)

	err = g.AddEdge(1, 3)
	assert.Nil(t, err)

	result := g.ToArray()
	assert.NotNil(t, result)

	jsonResult, jsonErr := json.Marshal(result)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, jsonResult)

	assert.Contains(t, string(jsonResult), "value2")
	assert.Contains(t, string(jsonResult), "value3")
}

func TestAddEdgeOK(t *testing.T) {
	g := NewDirectedGraph()

	g.AddVertex(0, nil)
	g.AddVertex(1, nil)
	g.AddVertex(2, "value2")
	g.AddVertex(3, "value3")

	err := g.AddEdge(0, 1)
	assert.Nil(t, err)

	err = g.AddEdge(1, 2)
	assert.Nil(t, err)

	err = g.AddEdge(1, 3)
	assert.Nil(t, err)

	// testing an already connecting nodes
	err = g.AddEdge(1, 3)
	assert.Nil(t, err)

	assert.Len(t, g.Vertices[0].Vertices, 1)
	assert.Len(t, g.Vertices[1].Vertices, 2)
	assert.Equal(t, g.Vertices[0].Vertices[1].Key, 1)
	assert.Equal(t, g.Vertices[1].Vertices[2].Key, 2)
	assert.Equal(t, g.Vertices[1].Vertices[3].Key, 3)
}

func TestAddEdgeNotExistVerticesError(t *testing.T) {
	g := NewDirectedGraph()

	g.AddVertex(0, nil)
	g.AddVertex(1, nil)
	g.AddVertex(2, "value2")
	g.AddVertex(3, "value3")

	err := g.AddEdge(0, 1)
	assert.Nil(t, err)

	err = g.AddEdge(1, 2)
	assert.Nil(t, err)

	err = g.AddEdge(1, 3)
	assert.Nil(t, err)

	// not exist the first vertice key
	err = g.AddEdge(5, 1)
	assert.NotNil(t, err)
	assert.Equal(t, "not all vertices exists", err.Error())

	// not exist the second vertice key
	err = g.AddEdge(1, 4)
	assert.NotNil(t, err)
	assert.Equal(t, "not all vertices exists", err.Error())

	// not exist any of the vertices
	err = g.AddEdge(4, 5)
	assert.NotNil(t, err)
	assert.Equal(t, "not all vertices exists", err.Error())
}

func TestToInterfaceBool(t *testing.T) {
	testCases := []struct {
		Name     string
		DataInfo DataInfo
		Value    interface{}
		Err      error
	}{
		{"parsed_bool", DataInfo{"bool", "false"}, false, nil},
		{"parsed_float", DataInfo{"float64", "22"}, float64(22), nil},
		{"parsed_float_with_decimal", DataInfo{"float64", "1.99"}, 1.99, nil},
		{"parsed_error", DataInfo{"float64", "false"}, nil, errors.New("error parsing flat_data")},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := tc.DataInfo.toInterface()
			if tc.Err == nil {
				assert.Nil(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, result, tc.Value)
			} else {
				assert.NotNil(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), tc.Err.Error())
			}
		})
	}
}

func TestBuildGraphFromVertexSecuenceOK(t *testing.T) {
	vtxSecuences := make([]VertexSecuence, 0)
	vtx0 := VertexSecuence{0, DataInfo{}, []int{1}}
	vtx1 := VertexSecuence{1, DataInfo{}, []int{2, 3}}
	vtx2 := VertexSecuence{2, DataInfo{"string", "value2"}, []int{}}
	vtx3 := VertexSecuence{3, DataInfo{"string", "value2"}, []int{}}
	vtxSecuences = append(vtxSecuences, vtx0, vtx1, vtx2, vtx3)

	g, err := BuildGraphFromVertexSecuence(vtxSecuences)
	assert.Nil(t, err)
	assert.NotNil(t, g)

	assert.Len(t, g.Vertices, 4)
	assert.Len(t, g.Vertices[0].Vertices, 1)
	assert.Len(t, g.Vertices[1].Vertices, 2)
	assert.Len(t, g.Vertices[2].Vertices, 0)
	assert.Len(t, g.Vertices[3].Vertices, 0)
}

func TestBuildGraphFromVertexSecuenceErrorParsing(t *testing.T) {
	vtxSecuences := make([]VertexSecuence, 0)
	vtx0 := VertexSecuence{0, DataInfo{}, []int{1}}
	vtx1 := VertexSecuence{1, DataInfo{}, []int{2, 3}}
	vtx2 := VertexSecuence{2, DataInfo{"string", "value2"}, []int{}}
	// this contains the error type
	vtx3 := VertexSecuence{3, DataInfo{"float64", "value2"}, []int{}}
	vtxSecuences = append(vtxSecuences, vtx0, vtx1, vtx2, vtx3)

	g, err := BuildGraphFromVertexSecuence(vtxSecuences)
	assert.NotNil(t, err)
	assert.Nil(t, g)
	assert.Equal(t, "error parsing data_info", err.Message())
}

func TestBuildGraphFromVertexSecuenceErrorAddingEdge(t *testing.T) {
	vtxSecuences := make([]VertexSecuence, 0)
	vtx0 := VertexSecuence{0, DataInfo{}, []int{1}}
	// this is the bad edge, contains a node that not exists
	vtx1 := VertexSecuence{1, DataInfo{}, []int{2, 4}}
	vtx2 := VertexSecuence{2, DataInfo{"string", "value2"}, []int{}}
	vtx3 := VertexSecuence{3, DataInfo{"string", "value2"}, []int{}}
	vtxSecuences = append(vtxSecuences, vtx0, vtx1, vtx2, vtx3)

	g, err := BuildGraphFromVertexSecuence(vtxSecuences)
	assert.NotNil(t, err)
	assert.Nil(t, g)
	assert.Equal(t, "not all vertices exists", err.Message())
}

func TestGetTypeAndValueStringFromInterface(t *testing.T) {
	_, _, err := getTypeAndValueStringFromInterface(nil)
	assert.NotNil(t, err)
	assert.Equal(t, "cannot get type and value from nil interface", err.Error())

	testCases := []struct {
		Name      string
		DataType  string
		DataValue string
		Value     interface{}
	}{
		{"parsing_string", "string", "test", "test"},
		{"parsing_float64", "float64", "25", float64(25)},
		{"parsing_float64_with_decimal", "float64", "1.99", 1.99},
		{"parsing_bool", "bool", "false", false},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			dt, dv, err := getTypeAndValueStringFromInterface(tc.Value)
			assert.Nil(t, err)
			assert.Equal(t, tc.DataType, dt)
			assert.Equal(t, tc.DataValue, dv)
		})
	}
}
