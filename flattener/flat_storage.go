package flattener

import (
	"context"
	"fmt"

	"github.com/mendezdev/tgo_flattener/apierrors"
	"github.com/mendezdev/tgo_flattener/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//go:generate mockgen -destination=mock_storage.go -package=flattener -source=flat_storage.go Storage

const (
	FlatCollection = "flats"
	DbName         = "flattenerdb"
	DbNameTest     = "flattenerdbtest"
)

// Storage will execute all de CRUD operations flat_info related
type Storage interface {
	create(FlatInfo) apierrors.RestErr
	getAll() ([]FlatInfo, apierrors.RestErr)
}

type storage struct {
	db     *mongo.Client
	dbName string
}

func NewStorage(db *mongo.Client) Storage {
	return &storage{
		db,
		DbName,
	}
}

func NewTestStorage(db *mongo.Client) Storage {
	return &storage{
		db,
		DbNameTest,
	}
}

func (s *storage) create(fi FlatInfo) apierrors.RestErr {
	collection := s.db.Database(s.dbName).Collection(FlatCollection)
	insertResult, err := collection.InsertOne(context.TODO(), fi)

	if err != nil {
		return apierrors.NewInternalServerError(fmt.Sprintf("database error creating flat_info: %s", err.Error()))
	}

	if insertedID, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
		fi.ID = insertedID.Hex()
	}

	return nil
}

func (s *storage) getAll() ([]FlatInfo, apierrors.RestErr) {
	collection := s.db.Database(s.dbName).Collection(FlatCollection)
	ctx := context.TODO()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"processed_at", -1}}).SetLimit(config.FlatsLimit)
	cursor, err := collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, apierrors.NewInternalServerError(fmt.Sprintf("database error getting all flat_info: %s", err.Error()))
	}

	var res []FlatInfo
	if cursorErr := cursor.All(ctx, &res); cursorErr != nil {
		return nil, apierrors.NewInternalServerError(fmt.Sprintf("database error iterating cursor of all flat_info: %s", err.Error()))
	}

	return res, nil
}
