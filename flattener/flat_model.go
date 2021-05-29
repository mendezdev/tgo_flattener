package flattener

type FlatResponse struct {
	MaxDepth int           `json:"max_depth"`
	Data     []interface{} `json:"flatted_data"`
}

type Flat struct {
	ID          string        `json:"id" bson:"_id,omitempty"`
	DateCreated string        `json:"date_created" bson:"date_creted"`
	Unflatted   []interface{} `json:"unflatted" bson:"unflatted"`
	Flatted     []interface{} `json:"flatted" bson:"flatted"`
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
