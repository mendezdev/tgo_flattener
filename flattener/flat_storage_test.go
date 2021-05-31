package flattener

import (
	"context"
	"testing"
	"time"

	"github.com/mendezdev/tgo_flattener/config"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestGetAllFlats(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	assert.Nil(t, err)
	assert.NotNil(t, client)
	defer client.Disconnect(ctx)

	storage := NewTestStorage(client)
	qtyNewDocuments, createErr := createFlatsInfo(storage)
	assert.Nil(t, createErr)

	flats, getErr := storage.getAll()
	assert.Nil(t, getErr)
	assert.NotNil(t, flats)
	assert.Equal(t, config.FlatsLimit, int64(len(flats)))

	now := time.Now().UTC()
	var counter int
	for _, f := range flats {
		if f.ProcessedAt.After(now) {
			counter++
		}
	}
	assert.Equal(t, qtyNewDocuments, counter)

	dropErr := client.Database(DbNameTest).Collection(FlatCollection).Drop(context.Background())
	assert.Nil(t, dropErr)
}

// this creates a 140 records:
// 90 of them are 1 day after now to simulate a recent and old records
// with this, the getAll can check if it is getting the last ones
func createFlatsInfo(s Storage) (qtyNewDocuments int, err error) {
	qtyOldDocuments := 50
	qtyNewDocuments = 90
	newProcessedTime := time.Now().UTC().Add(time.Hour * 24)
	oldProcessedTime := time.Now().UTC()

	for i := 0; i < qtyOldDocuments; i++ {
		fi := buildFlatInfo(oldProcessedTime)
		err := s.create(fi)
		if err != nil {
			return qtyNewDocuments, err
		}
	}

	for i := 0; i < qtyNewDocuments; i++ {
		fi := buildFlatInfo(newProcessedTime)
		err := s.create(fi)
		if err != nil {
			return qtyNewDocuments, err
		}
	}
	return qtyNewDocuments, nil
}

// is the same flat_info only for test purposes
func buildFlatInfo(processedAt time.Time) FlatInfo {
	vtxSecuences := make([]VertexSecuence, 0)
	vtx0 := VertexSecuence{0, DataInfo{}, []int{1}}
	vtx1 := VertexSecuence{1, DataInfo{}, []int{2, 3}}
	vtx2 := VertexSecuence{2, DataInfo{"string", "value2"}, []int{}}
	vtx3 := VertexSecuence{3, DataInfo{"string", "value2"}, []int{}}
	vtxSecuences = append(vtxSecuences, vtx0, vtx1, vtx2, vtx3)
	return FlatInfo{
		MaxDepth:       0,
		VertexSecuence: vtxSecuences,
		ProcessedAt:    processedAt,
	}
}
