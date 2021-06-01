package flattener

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/mendezdev/tgo_flattener/apierrors"
	"github.com/stretchr/testify/assert"
)

func TestPostFlatsOK(t *testing.T) {
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

func TestPostFlatBadRequest(t *testing.T) {
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
	c.Request, _ = http.NewRequest("POST", "/flats", strings.NewReader(string(jsonMockedRequest)))
	h.Post(c)

	assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
	assert.Contains(t, nr.Body.String(), "error parsing body")
}

func TestGetFlatsOK(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockGtw := NewMockGateway(mockCtrl)
	h := NewHandler(mockGtw)

	mockedResponse := mockFlatInfoResponse()
	jsonResponse, err := json.Marshal(mockedResponse)
	assert.Nil(t, err)
	assert.NotNil(t, jsonResponse)

	mockGtw.
		EXPECT().
		GetFlats().
		Return(mockedResponse, nil).
		Times(1)

	nr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(nr)
	c.Request, _ = http.NewRequest(http.MethodGet, "/flats", strings.NewReader(string(jsonResponse)))
	h.GetAll(c)

	var response []FlatInfoResponse
	jsonErr := json.Unmarshal(nr.Body.Bytes(), &response)
	assert.Nil(t, jsonErr)
	assert.Equal(t, http.StatusOK, c.Writer.Status())
	assert.NotNil(t, response)
	assert.Len(t, response, 1)
}

func TestGetFlatsError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockGtw := NewMockGateway(mockCtrl)
	h := NewHandler(mockGtw)

	msgErr := "error getting flats from database"
	mockGtw.
		EXPECT().
		GetFlats().
		Return(nil, apierrors.NewInternalServerError(msgErr)).
		Times(1)

	nr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(nr)
	c.Request, _ = http.NewRequest(http.MethodGet, "/flats", nil)
	h.GetAll(c)

	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
	assert.Contains(t, nr.Body.String(), msgErr)
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

func mockFlatInfoResponse() []FlatInfoResponse {
	unflatted := make([]interface{}, 0)
	item := make([]interface{}, 0)

	item = append(item, "lvl1_item0")
	item = append(item, "lvl1_item1")
	unflatted = append(unflatted, "string")
	unflatted = append(unflatted, item)

	flatted := []interface{}{"string", "lvl1_item0", "lvl1_item1"}
	firItem := FlatInfoResponse{
		ID:          "1234qwerty",
		ProcessedAt: time.Now().UTC(),
		Unflatted:   unflatted,
		Flatted:     flatted,
	}

	return []FlatInfoResponse{firItem}
}
