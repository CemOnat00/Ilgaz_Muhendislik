package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var MongoClient *mongo.Client
var DB *mongo.Database

const (
	dbName = "IlgazDB"
)

func ConnectDB() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI environment variable is not set")
	}

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

	// Default credentials (In production, move to .env)
	defaultUser := "admin"
	defaultPass := "ilgaz2026"

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(defaultPass), bcrypt.DefaultCost)

	// Update or Create admin
	opts := options.UpdateOne().SetUpsert(true)
	filter := bson.M{"username": defaultUser}
	update := bson.M{
		"$set": bson.M{
			"username": defaultUser,
			"password": string(hashedPassword),
			"email":    "admin@ilgazmuhendislik.com",
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Println("Admin seed hatası:", err)
	} else {
		log.Println("Admin hesabı güncellendi/oluşturuldu (Güvenli Hashlendi)")
	}
}
