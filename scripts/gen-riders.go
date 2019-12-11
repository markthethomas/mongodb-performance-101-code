package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/brianvoe/gofakeit"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type rider struct {
	ID        string `json:"_id" bson:"_id"`
	Name      string `json:"name" bson:"name"`
	Email     string `json:"email" bson:"email"`
	AvatarURL string `json:"avatarURL" bson:"avatarURL"`
	Password  string `json:"password" bson:"password"`
	Rides     []struct {
		ID       string `json:"_id" bson:"_id"`
		Date     string `json:"date" bson:"date"`
		DriverID string `json:"driverID" bson:"driverID,omitempty"`
		From     struct {
			Lat float64 `json:"lat" bson:"lat"`
			Lng float64 `json:"lng" bson:"lng"`
		} `json:"from" bson:"from"`
		To struct {
			Lat float64 `json:"lat" bson:"lat"`
			Lng float64 `json:"lng" bson:"lng"`
		} `json:"to" bson:"to"`
	} `json:"rides" bson:"rides"`
	Location struct {
		Current struct {
			Lat float64 `json:"lat" bson:"lat"`
			Lng float64 `json:"lng" bson:"lng"`
		} `json:"current" bson:"current"`
	} `json:"location" bson:"location"`
	RideNotes string `json:"rideNotes" bson:"rideNotes"`
}

func generateRandomRider() *bson.M {
	rides := bson.A{}
	randRides := gofakeit.Float32Range(0, 10)
	for index := 0; index < int(randRides); index++ {
		fakeRide := bson.M{
			"_id":      primitive.NewObjectID(),
			"date":     gofakeit.Date(),
			"driverID": primitive.NewObjectID(),
			"from": bson.M{
				"lat": gofakeit.Latitude(),
				"lng": gofakeit.Longitude(),
			},
			"to": bson.M{
				"lat": gofakeit.Latitude(),
				"lng": gofakeit.Longitude(),
			},
		}
		rides = append(rides, fakeRide)
	}
	return &bson.M{
		"_id":       primitive.NewObjectID(),
		"name":      gofakeit.Name() + " " + gofakeit.Name(),
		"email":     gofakeit.Name() + gofakeit.Email(),
		"avatarURL": gofakeit.ImageURL(60, 60),
		"password":  gofakeit.Password(true, true, true, true, true, 64),
		"rides":     rides,
		"location": bson.M{
			"current": bson.M{
				"lat": gofakeit.Latitude(),
				"lng": gofakeit.Longitude(),
			},
		},
		"rideNotes": gofakeit.Paragraph(1, int(gofakeit.Float32Range(1, 5)), int(gofakeit.Float32Range(3, 15)), " "),
	}
}

func main() {
	gofakeit.Seed(0)

	nDocs := flag.Int("docs", 50000, "Set a number of documents to generate into the riders collection")
	batchSize := flag.Int("batchsize", 5000, "Batch size of docs to insert at a time")
	connectionURL := flag.String("connectionurl", "mongodb://localhost:27017", "Pass a valid MongoDB URL")
	flag.Parse()

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(*connectionURL).SetMinPoolSize(50))
	if err != nil {
		log.Fatal(err)
	}
	coll := client.Database("floqars").Collection("riders")
	batch := []interface{}{}
	for index := 0; index < *nDocs; index++ {
		r := generateRandomRider()
		batch = append(batch, r)
		batch = append(batch)
		if index%*batchSize == 0 {
			if _, err := coll.InsertMany(context.Background(), batch); err != nil {
				log.Fatal(err)
			}
			batch = nil
		}

	}
	// flush remaining
	if _, err := coll.InsertMany(context.Background(), batch); err != nil {
		log.Fatal(err)
	}
	log.Printf("Generated %v docs!", *nDocs)
}
