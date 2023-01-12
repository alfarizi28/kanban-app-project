package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type TaskAPI interface {
	GetTask(w http.ResponseWriter, r *http.Request)
	CreateNewTask(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
	UpdateTaskCategory(w http.ResponseWriter, r *http.Request)
}

type taskAPI struct {
	taskService service.TaskService
}

func NewTaskAPI(taskService service.TaskService) *taskAPI {
	return &taskAPI{taskService}
}

func (t *taskAPI) GetTask(w http.ResponseWriter, r *http.Request) {
	getID := fmt.Sprintf("%s", r.Context().Value("id"))
	getIDInt, _ := strconv.Atoi(getID)

	if getID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	taskID := r.URL.Query().Get("task_id")
	taskIDInt, _ := strconv.Atoi(taskID)

	if taskID == "" {
		gTask, err := t.taskService.GetTasks(r.Context(), getIDInt)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
			return
		} else {
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(gTask)
			return
		}
	} else {
		geTask, err := t.taskService.GetTaskByID(r.Context(), taskIDInt)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
			return
		} else {
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(geTask)
			return
		}
	}

}

func (t *taskAPI) CreateNewTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest
	var tsk entity.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}

	catIDStr := fmt.Sprintf("%d", task.CategoryID)

	if task.Title == "" || task.Description == "" || catIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}

	getID := fmt.Sprintf("%s", r.Context().Value("id"))
	getIDInt, _ := strconv.Atoi(getID)

	if getID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	tsk = entity.Task{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		CategoryID:  task.CategoryID,
		UserID:      getIDInt,
	}
	getTsk, err := t.taskService.StoreTask(r.Context(), &tsk)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(201)
	success := map[string]interface{}{
		"user_id": getTsk.UserID,
		"task_id": getTsk.ID,
		"message": "success create new task",
	}
	json.NewEncoder(w).Encode(success)
	return

}

func (t *taskAPI) DeleteTask(w http.ResponseWriter, r *http.Request) {
	getID := fmt.Sprintf("%s", r.Context().Value("id"))
	getIDInt, _ := strconv.Atoi(getID)

	taskID := r.URL.Query().Get("task_id")
	taskIDInt, _ := strconv.Atoi(taskID)

	err := t.taskService.DeleteTask(r.Context(), taskIDInt)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(200)
	success := map[string]interface{}{
		"user_id": getIDInt,
		"task_id": taskIDInt,
		"message": "success delete task",
	}
	json.NewEncoder(w).Encode(success)
	return

}

func (t *taskAPI) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	getID := fmt.Sprintf("%s", r.Context().Value("id"))
	getIDInt, _ := strconv.Atoi(getID)

	if getID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}
	tsk := entity.Task{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		CategoryID:  task.CategoryID,
		UserID:      getIDInt,
	}
	getTsk, err := t.taskService.UpdateTask(r.Context(), &tsk)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(200)
	success := map[string]interface{}{
		"user_id": getTsk.UserID,
		"task_id": getTsk.ID,
		"message": "success update task",
	}
	json.NewEncoder(w).Encode(success)
	return

}

func (t *taskAPI) UpdateTaskCategory(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskCategoryRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	var updateTask = entity.Task{
		ID:         task.ID,
		CategoryID: task.CategoryID,
		UserID:     int(idLogin),
	}

	_, err = t.taskService.UpdateTask(r.Context(), &updateTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userId,
		"task_id": task.ID,
		"message": "success update task category",
	})
}
