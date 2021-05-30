package config

import "os"

const (
	FlatsLimit = int64(100)

	Production = "production"
	Develop    = "develop"
	Test       = "test"
)

var (
	Scope       = os.Getenv("SCOPE")
	Environment = os.Getenv("TGO_ENVIRONMENT")
)
