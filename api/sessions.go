package api

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"strings"
)

type sessionModel struct {
	Id       string `json:"_id"`
	Trackers string `json:"trackers"`
}

type sessionsGetResponse struct {
	Trackers string `json:"trackers"`
}

type sessionsPostResponse struct {
	Id string `json:"id"`
}

func Sessions(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		handleGet(w, r)
	} else if r.Method == http.MethodPost {
		handlePost(w, r)
	} else {
		log.Printf("Unsupported HTTP method used: %v", r.Method)
		w.Header().Set("Allow", strings.Join([]string{http.MethodGet, http.MethodPost}, ","))
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Printf("Id missing from GET request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessions := getSessionsCollection()
	if sessions == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Failed to create ObjectID from id '%v': %v", id, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var session sessionModel
	if err := sessions.
		FindOne(context.Background(), bson.M{"_id": objectId}).
		Decode(&session); err != nil {
		log.Printf("Failed to get entry for id '%v': %v", id, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(sessionsGetResponse{Trackers: session.Trackers})
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if _, err := r.Body.Read(body); err != nil {
		log.Printf("Failed to read body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	trackers := string(body)

	sessions := getSessionsCollection()
	if sessions == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id := r.URL.Query().Get("id")

	if id == "" {
		result, err := sessions.InsertOne(context.Background(), bson.D{{"trackers", trackers}})
		if err != nil {
			log.Printf("Failed to insert new session: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		objectId, ok := result.InsertedID.(primitive.ObjectID)
		if !ok {
			log.Printf("Inserted ID is not an ObjectID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(sessionsPostResponse{Id: objectId.Hex()})
	} else {
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Printf("Failed to create ObjectID from id: %v", id)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = sessions.ReplaceOne(
			context.Background(),
			bson.D{{"_id", objectId}},
			bson.D{{"trackers", trackers}})
		if err != nil {
			log.Printf("Failed to replace session: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(sessionsPostResponse{Id: id})
	}
}

func getSessionsCollection() *mongo.Collection {
	db, err := mongo.Connect(
		context.Background(),
		options.Client().
			ApplyURI(fmt.Sprint(
				"mongodb+srv://client:",
				os.Getenv("MONGODB_PASSWORD"),
				"@personal-3effx.azure.mongodb.net/test?retryWrites=true&w=majority")))
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil
	}

	return db.Database("timeTracker").Collection("sessions")
}
