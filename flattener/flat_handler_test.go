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
		FlatResponse(gomock.Any()).Return(mockedResponse).
		Times(1)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("POST", "/flat", strings.NewReader(string(jsonMockedResponse)))
	h.Post(c)

	assert.Equal(t, http.StatusOK, c.Writer.Status())
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

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("POST", "/flat", strings.NewReader(string(jsonMockedRequest)))
	h.Post(c)

	assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
}

func mockFlatRequest() []interface{} {
	return []interface{}{1, 2, 3, 4}
}

func mockFlatResponse() FlatResponse {
	return FlatResponse{
		MaxDepth: 0,
		Data:     []interface{}{1, 2, 3, 4},
	}
}
