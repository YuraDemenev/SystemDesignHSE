package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	taskservice "task-service/handlers"
	"testing"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestCreateTaskHandler(t *testing.T) {
	task := taskservice.Task{ID: -1, Text: "Test task"}
	body, _ := json.Marshal(task)
	req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(taskservice.CreateTaskHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestGetTasksHandler(t *testing.T) {
	taskservice.Tasks = append(taskservice.Tasks, taskservice.Task{ID: -1, Text: "Test task"})
	req, err := http.NewRequest("GET", "/tasks/list", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(taskservice.GetTasksHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var got []taskservice.Task
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}

	expected := []taskservice.Task{{ID: -1, Text: "Test task"}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
