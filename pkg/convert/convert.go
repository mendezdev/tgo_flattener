package convert

import (
	"errors"
	"fmt"
)

func GetTypeAndValueStringFromInterface(val interface{}) (dt string, dv string, err error) {
	if val == nil {
		err = errors.New("cannot get type and value from nil interface")
		return
	}
	dt = fmt.Sprintf("%T", val)
	dv = fmt.Sprintf("%v", val)
	return
}
