package domain

import (
	"fmt"
	"strconv"
)

// DataInfo represents the information of the value of the array node
type DataInfo struct {
	DataType  string `bson:"type"`
	DataValue string `bson:"value"`
}

// toInterface rebuild the original value of in the array
func (di DataInfo) ToInterface() (interface{}, error) {
	var convertedValue interface{}
	var err error

	switch di.DataType {
	case "float64":
		convertedValue, err = strconv.ParseFloat(di.DataValue, 64)
	case "bool":
		convertedValue, err = strconv.ParseBool(di.DataValue)
	case "": // in a v2, this should be improved by checking 'nil' or 'array' like an special data type and value
		convertedValue = nil
	default:
		convertedValue = di.DataValue
	}
	if err != nil {
		return nil, fmt.Errorf("error parsing flat_data: %s", err.Error())
	}
	return convertedValue, nil
}
