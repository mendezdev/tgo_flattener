package app

import (
	"github.com/mendezdev/tgo_flattener/flattener"
	"github.com/mendezdev/tgo_flattener/internal/storage"
)

type handlers struct {
	Flat flattener.Handler
}

func StartApplication() {
	db := storage.Connect("mongodb://localhost:27017")
	h := handlers{
		Flat: flattener.NewHandler(flattener.NewGateway(flattener.NewStorage(db))),
	}

	router := routes(h)
	router.Run(":8080")
}
