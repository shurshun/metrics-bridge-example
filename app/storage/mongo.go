package storage

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"metrics/bridge/app/types"
	"time"
)

type MongoStorage struct {
	client 	*mongo.Client
	db 		*mongo.Database
}

func (m *MongoStorage) Convert(e types.MetricEntity) bson.D {
	var doc bson.D

	for k, v := range e.Data {
		doc = append(doc, bson.E{k, v})
	}

	return doc
}

func (m *MongoStorage) Connect(dsn, db string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(dsn))
	_ = mongoClient.Connect(ctx)

	if err != nil {
		return err
	}

	m.client = mongoClient
	m.db = mongoClient.Database(db)

	return nil
}

func (m *MongoStorage) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return m.client.Disconnect(ctx)
}

func (m *MongoStorage) Commit(batch []types.MetricEntity) error {
	buckets := make(map[string][]interface{})

	for _, e := range batch {
		buckets[e.Bucket] = append(buckets[e.Bucket], m.Convert(e))
	}

	for bucket, data := range buckets {
		collection := m.db.Collection(bucket)
		_, err := collection.InsertMany(context.Background(), data)
		if err != nil {
			return err
		}
	}

	return nil
}
