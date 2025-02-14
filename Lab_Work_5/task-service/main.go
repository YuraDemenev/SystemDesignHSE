package main

import (
	"log"
	"net/http"
	taskservice "task-service/handlers"
)

func main() {
	http.HandleFunc("/tasks", taskservice.CreateTaskHandler)
	http.HandleFunc("/tasks/list", taskservice.GetTasksHandler)
	log.Println("Task service running on port 8001")
	log.Println(http.ListenAndServe(":8001", nil))
}

//
