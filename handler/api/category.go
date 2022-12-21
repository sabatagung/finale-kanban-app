package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type CategoryAPI interface {
	GetCategory(w http.ResponseWriter, r *http.Request)
	CreateNewCategory(w http.ResponseWriter, r *http.Request)
	DeleteCategory(w http.ResponseWriter, r *http.Request)
	GetCategoryWithTasks(w http.ResponseWriter, r *http.Request)
}

type categoryAPI struct {
	categoryService service.CategoryService
}

func NewCategoryAPI(categoryService service.CategoryService) *categoryAPI {
	return &categoryAPI{categoryService}
}

func (c *categoryAPI) GetCategory(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id")
	if userId == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	id, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	categories, err := c.categoryService.GetCategories(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)

}

func (c *categoryAPI) CreateNewCategory(w http.ResponseWriter, r *http.Request) {
	var category entity.CategoryRequest

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid category request"))
		return
	}

	userId := r.Context().Value("id")
	if userId == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	id, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}
	if category.Type == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid category request"))
		return
	}

	categoryResponse, err := c.categoryService.StoreCategory(r.Context(), &entity.Category{
		Type:   category.Type,
		UserID: id,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(categoryResponse)

}

func (c *categoryAPI) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id")
	if userId == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	categoryID := r.URL.Query().Get("category_id")
	if categoryID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid category id"))
		return
	}

	id, err := strconv.Atoi(categoryID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid category id"))
		return
	}

	response := c.categoryService.DeleteCategory(r.Context(), id)
	if response != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func (c *categoryAPI) GetCategoryWithTasks(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("get category task", err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	categories, err := c.categoryService.GetCategoriesWithTasks(r.Context(), int(idLogin))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("internal server error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)

}
