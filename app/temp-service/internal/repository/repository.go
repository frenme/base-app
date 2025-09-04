package repository

import (
	"context"
	"encoding/json"
	"fmt"
	shareddb "shared/pkg/db"
	"temp/internal/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Repository struct {
	db shareddb.Postgres
}

func NewRepository(db shareddb.Postgres) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetCountDataMongo(ctx context.Context) (int64, error) {
	db := db.MongoDB.Database("tempDB")
	collection := db.Collection("counter")
	return collection.CountDocuments(ctx, bson.M{})
}

func (r *Repository) GetAllDataMongo(ctx context.Context) ([]bson.M, error) {
	db := db.MongoDB.Database("tempDB")
	collection := db.Collection("counter")
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	err = cursor.All(ctx, &results)
	return results, err
}

func (r *Repository) SetDataMongo(ctx context.Context) error {
	database := db.MongoDB.Database("tempDB")
	collection := database.Collection("counter")

	count, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return err
	}
	if count >= 100000 {
		return nil
	}

	var docs []interface{}
	for i := 0; i < 100000; i++ {
		docs = append(docs, bson.D{
			{Key: "index", Value: i},
			{Key: "value", Value: fmt.Sprintf("data_%d", i)},
		})
	}
	_, err = collection.InsertMany(ctx, docs)
	return err
}

func (r *Repository) GetDataRedis(ctx context.Context) (string, error) {
	return db.RedisDB.Get(ctx, "data_key").Result()
}

func (r *Repository) SetDataRedis(ctx context.Context, data interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return db.RedisDB.Set(ctx, "data_key", dataBytes, 10*time.Minute).Err()
}
