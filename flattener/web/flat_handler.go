package flattener

import (
	"net/http"

	"github.com/gin-gonic/gin"
	models "github.com/mendezdev/tgo_flattener/flattener/models"
)

type FlatService interface {
	FlatHandler(c *gin.Context)
}

type flatServiceImpl struct{}

func NewFlatService() FlatService {
	return &flatServiceImpl{}
}

func (s *flatServiceImpl) FlatHandler(c *gin.Context) {
	var unflatted []interface{}
	if err := c.ShouldBindJSON(&unflatted); err != nil {
		c.JSON(http.StatusBadRequest, "error parsing body")
		return
	}

	flatResponse := getFlattenedResponse(unflatted)
	c.JSON(http.StatusOK, flatResponse)
}

func getFlattenedResponse(req []interface{}) models.FlatResponse {
	parsedFlat := flatRecursive(req, 0)

	fr := models.FlatResponse{
		Data: make([]interface{}, 0),
	}

	for k, v := range parsedFlat {
		fr.Data = append(fr.Data, v...)
		if k > fr.MaxDepth {
			fr.MaxDepth = k
		}
	}

	return fr
}

func flatRecursive(arr []interface{}, depth int) map[int][]interface{} {
	flatResult := make(map[int][]interface{})
	flatResult[depth] = make([]interface{}, 0)

	for _, v := range arr {
		nextDepthArr, ok := v.([]interface{})
		if ok {
			newFlatResult := flatRecursive(nextDepthArr, depth+1)
			for k, nfr := range newFlatResult {
				flatResult[k] = append(flatResult[k], nfr...)
			}
			continue
		}
		flatResult[depth] = append(flatResult[depth], v)
	}

	return flatResult
}
