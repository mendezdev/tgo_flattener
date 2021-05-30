package app

import (
	"fmt"

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
	if db != nil {
		fmt.Println("testing db not nil")
	}
	router := routes(h)

	router.Run(":8080")
}
