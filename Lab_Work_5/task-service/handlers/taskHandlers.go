package taskservice

import (
	"encoding/json"
	"net/http"
)

type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

var Tasks = []Task{}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	Tasks = append(Tasks, task)
	w.WriteHeader(http.StatusCreated)
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Tasks)
}
