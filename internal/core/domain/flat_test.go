package domain_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/mendezdev/tgo_flattener/internal/core/domain"
	"github.com/mendezdev/tgo_flattener/pkg/convert"
	"github.com/stretchr/testify/assert"
)

func TestNewGraph(t *testing.T) {
	f := domain.NewFlat()
	assert.NotNil(t, f)
	assert.NotNil(t, f.Graph)
	assert.NotNil(t, f.Graph.Vertices)
}

func TestNewVertex(t *testing.T) {
	vtx := domain.NewVertex(1, "some_value")

	assert.NotNil(t, vtx)
	assert.Equal(t, vtx.Key, 1)
	assert.Equal(t, vtx.Value, "some_value")
	assert.NotNil(t, vtx.Vertices)
}

func TestVertexToArraySingleValue(t *testing.T) {
	mockValue := "some_value"
	vtx := domain.NewVertex(1, mockValue)
	result := vtx.ToArray()

	assert.NotNil(t, result)
	assert.IsType(t, mockValue, result)
}

func TestVertexToArrayInvolvedInArrayValue(t *testing.T) {
	mockValue := "first_lv2"
	vtxLvl1 := domain.NewVertex(1, "first_lvl")
	vtxLvl2 := domain.NewVertex(2, mockValue)
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
	f := domain.NewFlat()

	// re-creating the following array: [["value2","value3"]]
	f.Graph.AddVertex(0, nil)
	f.Graph.AddVertex(1, nil)
	f.Graph.AddVertex(2, "value2")
	f.Graph.AddVertex(3, "value3")

	err := f.Graph.AddEdge(0, 1)
	assert.Nil(t, err)

	err = f.Graph.AddEdge(1, 2)
	assert.Nil(t, err)

	err = f.Graph.AddEdge(1, 3)
	assert.Nil(t, err)

	result := f.Graph.ToArray()
	assert.NotNil(t, result)

	jsonResult, jsonErr := json.Marshal(result)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, jsonResult)

	assert.Contains(t, string(jsonResult), "value2")
	assert.Contains(t, string(jsonResult), "value3")
}

func TestAddEdgeOK(t *testing.T) {
	f := domain.NewFlat()

	f.Graph.AddVertex(0, nil)
	f.Graph.AddVertex(1, nil)
	f.Graph.AddVertex(2, "value2")
	f.Graph.AddVertex(3, "value3")

	err := f.Graph.AddEdge(0, 1)
	assert.Nil(t, err)

	err = f.Graph.AddEdge(1, 2)
	assert.Nil(t, err)

	err = f.Graph.AddEdge(1, 3)
	assert.Nil(t, err)

	// testing an already connecting nodes
	err = f.Graph.AddEdge(1, 3)
	assert.Nil(t, err)

	assert.Len(t, f.Graph.Vertices[0].Vertices, 1)
	assert.Len(t, f.Graph.Vertices[1].Vertices, 2)
	assert.Equal(t, f.Graph.Vertices[0].Vertices[1].Key, 1)
	assert.Equal(t, f.Graph.Vertices[1].Vertices[2].Key, 2)
	assert.Equal(t, f.Graph.Vertices[1].Vertices[3].Key, 3)
}

func TestAddEdgeNotExistVerticesError(t *testing.T) {
	f := domain.NewFlat()

	f.Graph.AddVertex(0, nil)
	f.Graph.AddVertex(1, nil)
	f.Graph.AddVertex(2, "value2")
	f.Graph.AddVertex(3, "value3")

	err := f.Graph.AddEdge(0, 1)
	assert.Nil(t, err)

	err = f.Graph.AddEdge(1, 2)
	assert.Nil(t, err)

	err = f.Graph.AddEdge(1, 3)
	assert.Nil(t, err)

	// not exist the first vertice key
	err = f.Graph.AddEdge(5, 1)
	assert.NotNil(t, err)
	assert.Equal(t, "not all vertices exists", err.Error())

	// not exist the second vertice key
	err = f.Graph.AddEdge(1, 4)
	assert.NotNil(t, err)
	assert.Equal(t, "not all vertices exists", err.Error())

	// not exist any of the vertices
	err = f.Graph.AddEdge(4, 5)
	assert.NotNil(t, err)
	assert.Equal(t, "not all vertices exists", err.Error())
}

func TestToInterfaceBool(t *testing.T) {
	testCases := []struct {
		Name     string
		DataInfo domain.DataInfo
		Value    interface{}
		Err      error
	}{
		{"parsed_bool", domain.DataInfo{"bool", "false"}, false, nil},
		{"parsed_float", domain.DataInfo{"float64", "22"}, float64(22), nil},
		{"parsed_float_with_decimal", domain.DataInfo{"float64", "1.99"}, 1.99, nil},
		{"parsed_error", domain.DataInfo{"float64", "false"}, nil, errors.New("error parsing flat_data")},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := tc.DataInfo.ToInterface()
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
	f := domain.NewFlat()
	vtxSecuences := make([]domain.VertexSecuence, 0)
	vtx0 := domain.VertexSecuence{0, domain.DataInfo{}, []int{1}}
	vtx1 := domain.VertexSecuence{1, domain.DataInfo{}, []int{2, 3}}
	vtx2 := domain.VertexSecuence{2, domain.DataInfo{"string", "value2"}, []int{}}
	vtx3 := domain.VertexSecuence{3, domain.DataInfo{"string", "value2"}, []int{}}
	vtxSecuences = append(vtxSecuences, vtx0, vtx1, vtx2, vtx3)

	g, err := f.BuildGraphFromVertexSecuence(vtxSecuences)
	assert.Nil(t, err)
	assert.NotNil(t, g)

	assert.Len(t, g.Vertices, 4)
	assert.Len(t, g.Vertices[0].Vertices, 1)
	assert.Len(t, g.Vertices[1].Vertices, 2)
	assert.Len(t, g.Vertices[2].Vertices, 0)
	assert.Len(t, g.Vertices[3].Vertices, 0)
}

func TestBuildGraphFromVertexSecuenceErrorParsing(t *testing.T) {
	f := domain.NewFlat()
	vtxSecuences := make([]domain.VertexSecuence, 0)
	vtx0 := domain.VertexSecuence{0, domain.DataInfo{}, []int{1}}
	vtx1 := domain.VertexSecuence{1, domain.DataInfo{}, []int{2, 3}}
	vtx2 := domain.VertexSecuence{2, domain.DataInfo{"string", "value2"}, []int{}}
	// this contains the error type
	vtx3 := domain.VertexSecuence{3, domain.DataInfo{"float64", "value2"}, []int{}}
	vtxSecuences = append(vtxSecuences, vtx0, vtx1, vtx2, vtx3)

	g, err := f.BuildGraphFromVertexSecuence(vtxSecuences)
	assert.NotNil(t, err)
	assert.Nil(t, g)
	assert.Equal(t, "error parsing data_info", err.Message())
}

func TestBuildGraphFromVertexSecuenceErrorAddingEdge(t *testing.T) {
	f := domain.NewFlat()
	vtxSecuences := make([]domain.VertexSecuence, 0)
	vtx0 := domain.VertexSecuence{0, domain.DataInfo{}, []int{1}}
	// this is the bad edge, contains a node that not exists
	vtx1 := domain.VertexSecuence{1, domain.DataInfo{}, []int{2, 4}}
	vtx2 := domain.VertexSecuence{2, domain.DataInfo{"string", "value2"}, []int{}}
	vtx3 := domain.VertexSecuence{3, domain.DataInfo{"string", "value2"}, []int{}}
	vtxSecuences = append(vtxSecuences, vtx0, vtx1, vtx2, vtx3)

	g, err := f.BuildGraphFromVertexSecuence(vtxSecuences)
	assert.NotNil(t, err)
	assert.Nil(t, g)
	assert.Equal(t, "not all vertices exists", err.Message())
}

func TestGetTypeAndValueStringFromInterface(t *testing.T) {
	_, _, err := convert.GetTypeAndValueStringFromInterface(nil)
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
			dt, dv, err := convert.GetTypeAndValueStringFromInterface(tc.Value)
			assert.Nil(t, err)
			assert.Equal(t, tc.DataType, dt)
			assert.Equal(t, tc.DataValue, dv)
		})
	}
}
