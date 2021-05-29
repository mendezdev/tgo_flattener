package flattener

type FlatResponse struct {
	MaxDepth int           `json:"max_depth"`
	Data     []interface{} `json:"flatted_data"`
}

type FlatInfoResponse struct {
	ID          string        `json:"id"`
	DateCreated string        `json:"date_created"`
	Unflatted   []interface{} `json:"unflatted"`
	Flatted     []interface{} `json:"flatted"`
}

type FlatInfo struct {
	ID               string              `json:"id" bson:"_id,omitempty"`
	StructureInfo    []FlatStructureInfo `bson:"structure"`
	StructureFlatted []FlatData          `bson:"structure_flatted"`
	DateCreated      string              `bson:"date_created"`
}

type FlatStructureInfo struct {
	Level int        `bson:"level"`
	Data  []FlatData `bson:"data"`
}

type FlatData struct {
	DataType  string `bson:"type"`
	DataValue string `bson:"value"`
}
