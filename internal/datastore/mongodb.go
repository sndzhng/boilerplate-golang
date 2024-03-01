package datastore

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sndzhng/gin-template/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	contextTimeout, _ = context.WithTimeout(context.Background(), 10*time.Second)
	mongoClient       *mongo.Client
	MongoDatabase     *mongo.Database
)

func ConnectMongodb() {
	uri := fmt.Sprintf(
		"%s://%s:%s@%s/%s?%s",
		config.Datastore.Mongodb.Format,
		config.Datastore.Mongodb.User,
		config.Datastore.Mongodb.Password,
		config.Datastore.Mongodb.Host,
		config.Datastore.Mongodb.Database,
		config.Datastore.Mongodb.Options,
	)

	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	err = mongoClient.Connect(contextTimeout)
	if err != nil {
		log.Fatal(err)
	}

	err = mongoClient.Ping(contextTimeout, nil)
	if err != nil {
		log.Fatal(err)
	}

	MongoDatabase = mongoClient.Database(config.Datastore.Mongodb.Database)
}

func DisconnectMongodb() {
	err := mongoClient.Disconnect(contextTimeout)
	if err != nil {
		log.Fatal(err)
	}
}
