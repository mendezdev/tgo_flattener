package flattener

import (
	"context"
	"fmt"
	"time"

	"github.com/mendezdev/tgo_flattener/apierrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DbName         = "flattenerdb"
	FlatCollection = "flats"
)

type Storage interface {
	create(FlatInfo) apierrors.RestErr
	getAll() ([]FlatInfo, apierrors.RestErr)
}

type storage struct {
	db *mongo.Client
}

func NewStorage(db *mongo.Client) Storage {
	return &storage{
		db,
	}
}

func (s *storage) create(fi FlatInfo) apierrors.RestErr {
	fi.DateCreated = time.Now().UTC().String()

	collection := s.db.Database(DbName).Collection(FlatCollection)
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
	collection := s.db.Database(DbName).Collection(FlatCollection)
	ctx := context.TODO()
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, apierrors.NewInternalServerError(fmt.Sprintf("database error getting all flat_info: %s", err.Error()))
	}
	var res []FlatInfo
	if cursorErr := cursor.All(ctx, &res); cursorErr != nil {
		return nil, apierrors.NewInternalServerError(fmt.Sprintf("database error iterating cursor of all flat_info: %s", err.Error()))
	}
	return res, nil
}
