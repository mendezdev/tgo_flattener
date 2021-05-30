package flattener

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFlat_ValidRequest(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockGtw := NewMockGateway(mockCtrl)
	h := NewHandler(mockGtw)

	mockedRequest := mockFlatRequest()
	mockedResponse := mockFlatResponse()
	jsonMockedResponse, err := json.Marshal(mockedRequest)
	assert.Nil(t, err)
	assert.NotNil(t, jsonMockedResponse)

	mockGtw.
		EXPECT().
		FlatResponse(gomock.Any()).Return(mockedResponse, nil).
		Times(1)

	nr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(nr)
	c.Request, _ = http.NewRequest("POST", "/flat", strings.NewReader(string(jsonMockedResponse)))
	h.Post(c)

	var fr FlatResponse
	jsonErr := json.Unmarshal(nr.Body.Bytes(), &fr)
	assert.Nil(t, jsonErr)

	assert.Equal(t, http.StatusOK, c.Writer.Status())
	assert.Equal(t, mockedResponse, fr)
}

func TestFlat_BadRequest(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockGtw := NewMockGateway(mockCtrl)
	h := NewHandler(mockGtw)

	mockedRequest := make(map[string]string)
	mockedRequest["superkey"] = "supervalue"

	jsonMockedRequest, err := json.Marshal(mockedRequest)
	assert.Nil(t, err)
	assert.NotNil(t, jsonMockedRequest)

	mockGtw.
		EXPECT().
		FlatResponse(gomock.Any()).
		Times(0)

	nr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(nr)
	c.Request, _ = http.NewRequest("POST", "/flat", strings.NewReader(string(jsonMockedRequest)))
	h.Post(c)

	assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
	assert.Contains(t, nr.Body.String(), "error parsing body")
}

func mockFlatRequest() []interface{} {
	return []interface{}{"test1", "test2", "test3"}
}

func mockFlatResponse() FlatResponse {
	return FlatResponse{
		MaxDepth: 0,
		Data:     []interface{}{"test1", "test2", "test3"},
	}
}
