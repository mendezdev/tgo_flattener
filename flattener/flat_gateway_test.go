package flattener

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/mendezdev/tgo_flattener/apierrors"
	"github.com/stretchr/testify/assert"
)

func TestFlatResponse_OK(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStorage := NewMockStorage(mockCtrl)
	gwt := NewGateway(mockStorage)

	mockStorage.
		EXPECT().
		create(gomock.Any()).Return(nil).
		Times(7)

	testCases := []struct {
		Name      string
		UseCaseFn func() ([]interface{}, error)
		MaxDepth  int
		Len       int
	}{
		{"level_0", buildDepthLevel0, 0, 3},
		{"level_1", buildDepthLevel1, 1, 3},
		{"level_2", buildDepthLevel2, 2, 5},
		{"level_3", buildDepthLevel3, 3, 8},
		{"level_4", buildDepthLevel4, 4, 1},
		{"level_5", buildDepthLevel5, 5, 2},
		{"level_6", buildDepthLevel6, 6, 14},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			useCase, err := tc.UseCaseFn()
			assert.Nil(t, err)
			assert.NotNil(t, useCase)

			fr, apiErr := gwt.FlatResponse(useCase)
			assert.Nil(t, apiErr)

			assert.NotNil(t, fr)
			assert.Equal(t, tc.MaxDepth, fr.MaxDepth)
			assert.Equal(t, tc.Len, len(fr.Data))
		})
	}
}

func TestFlatResponseDatabaseError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStorage := NewMockStorage(mockCtrl)
	gwt := NewGateway(mockStorage)

	dbErr := apierrors.NewInternalServerError("db error")
	mockStorage.
		EXPECT().
		create(gomock.Any()).Return(dbErr).
		Times(1)

	input, buildErr := buildDepthLevel0()
	assert.Nil(t, buildErr)
	assert.NotNil(t, input)

	_, apiErr := gwt.FlatResponse(input)
	assert.NotNil(t, apiErr)
	assert.Equal(t, "error saving the flat_info", apiErr.Message())
}

func TestGetFlatsOk(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStorage := NewMockStorage(mockCtrl)
	gwt := NewGateway(mockStorage)

	mockFlatInfo := getMockFlatInfo()
	mockStorage.
		EXPECT().
		getAll().
		Return(mockFlatInfo, nil).
		Times(1)

	flats, apiErr := gwt.GetFlats()
	assert.Nil(t, apiErr)
	assert.NotNil(t, flats)
	assert.Len(t, flats, 1)
}

func TestGetFlatsDbError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStorage := NewMockStorage(mockCtrl)
	gwt := NewGateway(mockStorage)

	dbErr := apierrors.NewInternalServerError("database error")
	mockStorage.
		EXPECT().
		getAll().
		Return(nil, dbErr).
		Times(1)

	flats, apiErr := gwt.GetFlats()
	assert.NotNil(t, apiErr)
	assert.Nil(t, flats)
	assert.Equal(t, "error getting flat_info from db", apiErr.Message())
}

func getMockFlatInfo() []FlatInfo {
	vs := []VertexSecuence{
		{
			Key:      0,
			DataInfo: DataInfo{DataType: "float64", DataValue: "1"},
			Edges:    []int{},
		},
		{
			Key:      1,
			DataInfo: DataInfo{DataType: "float64", DataValue: "2"},
			Edges:    []int{},
		},
		{
			Key:      2,
			DataInfo: DataInfo{DataType: "float64", DataValue: "3"},
			Edges:    []int{},
		},
	}

	fi := []FlatInfo{
		{
			ID:             "qwery12345",
			MaxDepth:       0,
			VertexSecuence: vs,
			DateCreated:    time.Now().UTC(),
		},
	}
	return fi
}

// returns depth: 0, total_items: 3
func buildDepthLevel0() ([]interface{}, error) {
	b := []byte(`[1,2,3]`)
	return unmarshalDepth(b)
}

// returns depth: 1, total_items: 3
func buildDepthLevel1() ([]interface{}, error) {
	b := []byte(`[1,2,[3]]`)
	return unmarshalDepth(b)
}

// returns depth: 2, total_items: 5
func buildDepthLevel2() ([]interface{}, error) {
	b := []byte(`[1,2,[[false],3],["some"]]`)
	return unmarshalDepth(b)
}

// returns depth: 3, total_items: 8
func buildDepthLevel3() ([]interface{}, error) {
	b := []byte(`[1,2,[[false,"test",[8]],3,7],["some"]]`)
	return unmarshalDepth(b)
}

// returns depth: 4, total_items: 1
func buildDepthLevel4() ([]interface{}, error) {
	b := []byte(`[[[[[1]]]]]`)
	return unmarshalDepth(b)
}

// returns depth: 5, total_items: 2
func buildDepthLevel5() ([]interface{}, error) {
	b := []byte(`[[[[[1,["test"]]]]]]`)
	return unmarshalDepth(b)
}

// returns depth: 6, total_items: 14
func buildDepthLevel6() ([]interface{}, error) {
	b := []byte(`[[10,20,[["some",20,[30,2]]],[[20,[30,30.99,false,[100],2,[[101]]]]],[40]]]`)
	return unmarshalDepth(b)
}

func unmarshalDepth(b []byte) ([]interface{}, error) {
	var res []interface{}
	err := json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
