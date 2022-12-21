package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
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
	var task entity.TaskRequest
	_, err := t.taskService.GetTasks(r.Context(), task.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}
	if task.CategoryID == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
	return

}

func (t *taskAPI) CreateNewTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}
	if task.Title == "" || task.Description == "" || task.CategoryID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}
	var userTask = entity.Task{
		Title:       task.Title,
		Description: task.Description,
		CategoryID:  task.CategoryID,
	}
	userTask, err = t.taskService.StoreTask(r.Context(), &userTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userTask.ID,
		"task_id": userTask.CategoryID,
		"message": "success create new task",
	})

}

func (t *taskAPI) DeleteTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest
	taskID := r.URL.Query().Get("task_id")
	tID, _ := strconv.Atoi(taskID)

	err := t.taskService.DeleteTask(r.Context(), tID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": task.ID,
		"task_id": task.CategoryID,
		"message": "success delete task",
	})
	return

}

func (t *taskAPI) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task entity.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}
	if task.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("register data is empty"))
		return
	}

	var updateTask = entity.Task{
		ID:          task.CategoryID,
		Title:       task.Title,
		Description: task.Description,
		CategoryID:  task.CategoryID,
	}
	updateTask, err = t.taskService.UpdateTask(r.Context(), &updateTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": task.UserID,
		"task_id": task.ID,
		"message": "success update task",
	})
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
