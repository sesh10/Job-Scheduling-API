package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// JobStatus represents the status of a job
type JobStatus string

const (
	StatusQueued     JobStatus = "queued"
	StatusInProgress JobStatus = "in_progress"
	StatusCompleted  JobStatus = "completed"
	StatusFailed     JobStatus = "failed"
)

// JobStore is a thread-safe map for storing job statuses
var jobStore = sync.Map{}

func main() {
	godotenv.Load()
	r := mux.NewRouter()
	r.HandleFunc("/jobs", createJobHandler).Methods("POST")
	r.HandleFunc("/jobs/{id}", getJobStatusHandler).Methods("GET")

	port := os.Getenv("PORT")
	log.Println("Starting server on PORT:", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}

func createJobHandler(w http.ResponseWriter, r *http.Request) {
	jobID := uuid.New().String()
	jobStore.Store(jobID, StatusQueued)

	go processJob(jobID)

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

	value, ok := jobStore.Load(jobID)
	if !ok {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	status := value.(JobStatus)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"job_id": jobID,
		"status": string(status),
	})
}

func processJob(jobID string) {
	v := rand.Intn(10)
	if v < 4 {
		jobStore.Store(jobID, StatusFailed)
		return
	}
	time.Sleep(time.Duration(2*rand.Intn(5)) * time.Second)

	jobStore.Store(jobID, StatusInProgress)
	time.Sleep(time.Duration(5*v) * time.Second)

	jobStore.Store(jobID, StatusCompleted)

}
