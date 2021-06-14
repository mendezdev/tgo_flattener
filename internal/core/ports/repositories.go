package ports

import "github.com/mendezdev/tgo_flattener/apierrors"

type FlatRepository interface {
	Create() apierrors.RestErr
	GetAll() apierrors.RestErr
}
