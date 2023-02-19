package data

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

// client is our mongo client that allows us to perform operations on the Mongo database.
var client *mongo.Client

// LogEntry is the type for all data stored in the logs collection. Note that we specify
// specific bson values, and we *must* include omitempty on ID, or newly inserted records will
// have an empty id! We also specify JSON struct tags, even though we don't use them yet. We
// might in the future.
type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// Models is the type for this package. Note that any model that is included as a member
// in this type is available to us throughout the application, anywhere that the
// app variable is used, provided that the model is also added in the New function.
type Models struct {
	LogEntry LogEntry
}

// New is the function used to create an instance of the data package. It returns the type
// Model, which embeds all the types we want to be available to our application.

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{LogEntry: LogEntry{}}
}

// Insert puts a document in the logs collection.

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("logs").Collection("logs")

	el := LogEntry{
		ID:        entry.ID,
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := collection.InsertOne(context.TODO(), el)

	if err != nil {
		log.Println("Error inserting into logs:", err)
		return err
	}
	return nil
}
