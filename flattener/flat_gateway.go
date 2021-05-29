package flattener

import (
	"fmt"

	"github.com/mendezdev/tgo_flattener/apierrors"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate mockgen -destination=mock_gateway.go -package=flattener -source=flat_gateway.go Gateway

type Gateway interface {
	FlatResponse([]interface{}) FlatResponse
	GetFlats() ([]FlatInfo, apierrors.RestErr)
}

type gateway struct {
	Storage
}

func NewGateway(db *mongo.Client) Gateway {
	return &gateway{NewStorage(db)}
}

func (s *gateway) FlatResponse(req []interface{}) FlatResponse {
	parsedFlat := flatRecursive(req, 0)
	fr := FlatResponse{
		Data: make([]interface{}, 0),
	}

	structures := []FlatStructureInfo{}
	flattedStructure := make([]FlatData, 0)

	for k, dArr := range parsedFlat {
		fr.Data = append(fr.Data, dArr...)

		// find max depth
		if k > fr.MaxDepth {
			fr.MaxDepth = k
		}

		structure := FlatStructureInfo{
			Level: k,
			Data:  make([]FlatData, 0),
		}
		for _, d := range dArr {
			t := fmt.Sprintf("%T", d)
			v := fmt.Sprintf("%v", d)
			fd := FlatData{
				DataType:  t,
				DataValue: v,
			}
			structure.Data = append(structure.Data, fd)
			flattedStructure = append(flattedStructure, fd)
		}
		structures = append(structures, structure)
	}
	fi := FlatInfo{
		StructureInfo:    structures,
		StructureFlatted: flattedStructure,
	}

	saveErr := s.Storage.create(fi)
	if saveErr != nil {
		//TODO: return this with an apierror
		panic(saveErr)
	}

	fmt.Println("FLAT INFO SAVED!")
	return fr
}

func (s *gateway) GetFlats() ([]FlatInfo, apierrors.RestErr) {
	flats, err := s.Storage.getAll()
	if err != nil {
		return nil, err
	}
	return flats, nil
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

func toFlatInfo(data map[int][]interface{}) (FlatInfo, apierrors.RestErr) {
	structures := []FlatStructureInfo{}
	for k, dArr := range data {
		structure := FlatStructureInfo{
			Level: k,
			Data:  make([]FlatData, 0),
		}
		for _, d := range dArr {
			t := fmt.Sprintf("%T", d)
			v := fmt.Sprintf("%v", d)
			fd := FlatData{
				DataType:  t,
				DataValue: v,
			}
			structure.Data = append(structure.Data, fd)
		}
		structures = append(structures, structure)
	}
	return FlatInfo{
		StructureInfo: structures,
	}, nil
}
