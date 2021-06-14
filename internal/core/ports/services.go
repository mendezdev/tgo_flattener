package ports

import "github.com/mendezdev/tgo_flattener/apierrors"

type FlatService interface {
	Flat() apierrors.RestErr
	GetAll() apierrors.RestErr
}
