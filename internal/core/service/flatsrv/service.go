package flatsrv

import (
	"github.com/mendezdev/tgo_flattener/apierrors"
	"github.com/mendezdev/tgo_flattener/internal/core/ports"
)

type service struct {
	flatRepository ports.FlatRepository
}

func New(fr ports.FlatRepository) *service {
	return &service{
		flatRepository: fr,
	}
}

func (srv *service) Flat() apierrors.RestErr {
	return nil
}

func (srv *service) GetAll() apierrors.RestErr {
	return nil
}
