package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var MongoClient *mongo.Client
var DB *mongo.Database

const (
	uri    = "mongodb+srv://satircemonat_db_user:2PKxc2hByMp6j05e@ilgazmuhdb.inirtn0.mongodb.net/?appName=IlgazMuhDB"
	dbName = "IlgazDB"
)

func ConnectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal("MongoDB Bağlantı Hatası:", err)
	}

	// Ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB Atlas'a erişilemiyor:", err)
	}

	MongoClient = client
	DB = client.Database(dbName)
	log.Println("MongoDB Atlas'a başarıyla bağlanıldı!")
}

func GetCollection(collectionName string) *mongo.Collection {
	return DB.Collection(collectionName)
}

func SeedAdmin() {
	collection := GetCollection("admins")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if admin exists
	var admin map[string]interface{}
	err := collection.FindOne(ctx, bson.M{"username": "admin"}).Decode(&admin)
	if err != nil {
		// Not found, create it
		_, err = collection.InsertOne(ctx, bson.M{
			"username": "admin",
			"password": "ilgaz2026",
			"email":    "admin@ilgazmuhendislik.com",
		})
		if err != nil {
			log.Println("Admin seed hatası:", err)
		} else {
			log.Println("Varsayılan admin hesabı oluşturuldu (admin / ilgaz2026)")
		}
	}
}
