package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type JobStatus string

const (
	StatusQueued     JobStatus = "queued"
	StatusInProgress JobStatus = "in_progress"
	StatusCompleted  JobStatus = "completed"
	StatusFailed     JobStatus = "failed"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

func main() {
	godotenv.Load()

	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Test Redis connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/jobs", createJobHandler).Methods("POST")
	r.HandleFunc("/jobs/{id}", getJobStatusHandler).Methods("GET")

	port := os.Getenv("PORT")
	log.Println("Starting server on PORT:", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}

func createJobHandler(w http.ResponseWriter, r *http.Request) {
	jobID := uuid.New().String()
	err := setJobStatus(ctx, jobID, StatusQueued)
	if err != nil {
		http.Error(w, "Could not store job", http.StatusInternalServerError)
		return
	}

	go processJob(ctx, jobID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	log.Println("[INFO]: Created job", jobID)
	json.NewEncoder(w).Encode(map[string]string{
		"job_id": jobID,
		"status": string(StatusQueued),
	})
}

func getJobStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["id"]

	status, err := getJobStatus(ctx, jobID)
	if err == redis.Nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"job_id": jobID,
		"status": status,
	})
}

func processJob(ctx context.Context, jobID string) {
	v := rand.Intn(10)
	if v < 4 {
		msg := setJobStatus(ctx, jobID, StatusFailed)
		if msg != nil {
			log.Println("[ERROR]:", msg)
		}
		return
	}

	time.Sleep(time.Duration(2*rand.Intn(5)) * time.Second)
	msg := setJobStatus(ctx, jobID, StatusInProgress)
	if msg != nil {
		log.Println("[ERROR]:", msg)
	}

	time.Sleep(time.Duration(5*v) * time.Second)
	msg = setJobStatus(ctx, jobID, StatusCompleted)
	if msg != nil {
		log.Println("[ERROR]:", msg)
	}
}

func jobKey(id string) string {
	return fmt.Sprintf("job:%s", id)
}

func setJobStatus(ctx context.Context, id string, status JobStatus) error {
	return rdb.Set(ctx, jobKey(id), string(status), 0).Err()
}

func getJobStatus(ctx context.Context, id string) (string, error) {
	return rdb.Get(ctx, jobKey(id)).Result()
}
