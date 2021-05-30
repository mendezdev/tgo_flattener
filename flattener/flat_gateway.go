package flattener

import (
	"github.com/mendezdev/tgo_flattener/apierrors"
)

//go:generate mockgen -destination=mock_gateway.go -package=flattener -source=flat_gateway.go Gateway

type Gateway interface {
	FlatResponse([]interface{}) (FlatResponse, apierrors.RestErr)
	GetFlats() ([]FlatInfoResponse, apierrors.RestErr)
}

type gateway struct {
	storage Storage
}

func NewGateway(s Storage) Gateway {
	return &gateway{storage: s}
}

func (s *gateway) FlatResponse(input []interface{}) (FlatResponse, apierrors.RestErr) {
	var fr FlatResponse

	flatInfo, err := FlatArray(input)
	if err != nil {
		return fr, apierrors.NewInternalServerError("error flatting the array")
	}

	if dbErr := s.storage.create(flatInfo); dbErr != nil {
		return fr, apierrors.NewInternalServerError("error saving the flat_info")
	}

	fr.MaxDepth = flatInfo.MaxDepth
	fr.Data = flatInfo.Graph.ToFlat()

	return fr, nil
}

func (s *gateway) GetFlats() ([]FlatInfoResponse, apierrors.RestErr) {
	flats, err := s.storage.getAll()
	if err != nil {
		return nil, apierrors.NewInternalServerError("error getting flat_info from db")
	}

	res := make([]FlatInfoResponse, 0)
	for _, f := range flats {
		g, buildErr := BuildGraphFromVertexSecuence(f.VertexSecuence)
		if buildErr != nil {
			return nil, buildErr
		}
		unflatted := g.ToArray()
		flatted := g.ToFlat()
		fir := FlatInfoResponse{
			ID:          f.ID,
			DateCreated: f.DateCreated,
			Unflatted:   unflatted,
			Flatted:     flatted,
		}
		res = append(res, fir)
	}

	return res, nil
}
