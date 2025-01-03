package utils

import (
    "context"
    //"encoding/json"
    "os"
    "github.com/joho/godotenv"
    //"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "errors"
)

func DatabaseConnection() (*mongo.Client, context.Context, error) {
    
    if err := godotenv.Load(); err != nil {
	err1 := errors.New("No .env file found")
	return nil, nil, err1
    }

    url := os.Getenv("MONGO_URL")
    if url == "" {
        //log.Fatal("set your 'MONGO_URL in envirment vairable'" + docs + "usage examples/#environment-variable")
	err := errors.New("No mongourl")
	return nil, nil, err
    }
    ctx := context.TODO()
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
    if err != nil {
        return nil, nil, err
    }

    return client, ctx, nil
}
