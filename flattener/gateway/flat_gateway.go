package flattener

import (
	models "github.com/mendezdev/tgo_flattener/flattener/models"
)

type Gateway interface {
	FlatResponse([]interface{}) models.FlatResponse
}

type gateway struct{}

func NewGateway() Gateway {
	return &gateway{}
}

func (s *gateway) FlatResponse(req []interface{}) models.FlatResponse {
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
