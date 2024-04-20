package mongodb

import (
	"context"
	"franky/go-api-gin/envs"
	"sync"
	"time"

	"github.com/fatih/color"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoOnce      sync.Once
	DB             *mongo.Database
	UserCollection *mongo.Collection
)

func InitMongoDB() *mongo.Database {
	mongoOnce.Do(func() {
		var err error
		clientOptions := options.Client().ApplyURI(envs.MONGO_DATABASE_URL).
			SetConnectTimeout(10 * time.Second).
			SetMaxPoolSize(100)
		client, err := mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			panic(err)
		}

		DB = client.Database("go-api-gin")
		UserCollection = DB.Collection("users")

		color.New(color.Bold, color.BgGreen).Println("Connected to MongoDB")

		// User Collection Index
		name, err := DB.Collection("users").Indexes().CreateOne(context.Background(), mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		})
		if err != nil {
			panic(err)
		}
		color.New(color.FgGreen).Println("User collection index created: " + name)
	})

	return DB
}
