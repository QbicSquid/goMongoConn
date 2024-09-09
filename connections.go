package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Monster struct {
	Name     string `bson:"name"`     // Use BSON tags to match MongoDB document fields
	Category string `bson:"category"` // BSON is the format MongoDB uses
}

func connectToDatabases(DB_URI1, DB_URI2 string) (*mongo.Client, *mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions1 := options.Client().ApplyURI(DB_URI1)
	client1, err := mongo.Connect(ctx, clientOptions1)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to first database: %v", err)
	}
	if err := client1.Ping(ctx, nil); err != nil {
		return nil, nil, fmt.Errorf("failed to ping first database: %v", err)
	}

	clientOptions2 := options.Client().ApplyURI(DB_URI2)
	client2, err := mongo.Connect(ctx, clientOptions2)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to second database: %v", err)
	}
	if err := client2.Ping(ctx, nil); err != nil {
		return nil, nil, fmt.Errorf("failed to ping second database: %v", err)
	}

	return client1, client2, nil
}

func insertMonstersData(client *mongo.Client, databaseName, collectionName string, monsters []Monster) error {
	collection := client.Database(databaseName).Collection(collectionName)

	var docs []interface{}
	for _, monster := range monsters {
		docs = append(docs, monster)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("failed to insert monsters data: %v", err)
	}

	fmt.Println("Successfully inserted monsters data")
	return nil
}

func main() {
	DB_URI1 := "mongodb+srv://first:dummypass@first.ulhfa.mongodb.net/elemental?retryWrites=true&w=majority&appName=myAtlasClusterEDU"
	DB_URI2 := "mongodb+srv://second:dummypass@cluster0.pqvjtbp.mongodb.net/elemental?retryWrites=true&w=majority&appName=myAtlasClusterEDU"

	monstersData := []Monster{
		{
			Name:     "Katakan",
			Category: "Vampire",
		},
		{
			Name:     "Drowner",
			Category: "Drowner",
		},
		{
			Name:     "Nekker",
			Category: "Nekker",
		},
		{
			Name:     "Leshen",
			Category: "Relict",
		},
		{
			Name:     "Fiend",
			Category: "Relict",
		},
		{
			Name:     "Griffin",
			Category: "Hybrid",
		},
		{
			Name:     "Ekimma",
			Category: "Vampire",
		},
		{
			Name:     "Werewolf",
			Category: "Cursed One",
		},
		{
			Name:     "Basilisk",
			Category: "Draconid",
		},
		{
			Name:     "Chort",
			Category: "Relict",
		},
		{
			Name:     "Forktail",
			Category: "Draconid",
		},
		{
			Name:     "Harpie",
			Category: "Hybrid",
		},
		{
			Name:     "Succubus",
			Category: "Relict",
		},
	}

	client1, client2, err := connectToDatabases(DB_URI1, DB_URI2)
	if err != nil {
		log.Fatalf("Error connecting to databases: %v", err)
	}

	if err := insertMonstersData(client1, "elemental", "monsters", monstersData); err != nil {
		log.Fatalf("Error inserting monsters data: %v", err)
	}

	defer func() {
		if err := client1.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Error disconnecting from first database: %v", err)
		}
		if err := client2.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Error disconnecting from second database: %v", err)
		}
		fmt.Println("Disconnected from both databases!")
	}()
}
