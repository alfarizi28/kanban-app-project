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
	getID := fmt.Sprintf("%s", r.Context().Value("id"))

	if getID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	ID, _ := strconv.Atoi(getID)
	getCategory, err := c.categoryService.GetCategories(r.Context(), ID)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(getCategory)
	return

}

func (c *categoryAPI) CreateNewCategory(w http.ResponseWriter, r *http.Request) {
	var category entity.CategoryRequest
	var cat entity.Category

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid category request"))
		return
	}

	if category.Type == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid category request"))
		return
	}

	getID := fmt.Sprintf("%s", r.Context().Value("id"))
	// fmt.Println(getID)

	if getID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	getIDInt, _ := strconv.Atoi(getID)
	// fmt.Println(getIDInt)
	cat = entity.Category{
		Type:   category.Type,
		UserID: getIDInt,
	}
	getCat, err := c.categoryService.StoreCategory(r.Context(), &cat)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(201)
	success := map[string]interface{}{
		"user_id":     getCat.UserID,
		"category_id": getCat.ID,
		"message":     "success create new category",
	}
	json.NewEncoder(w).Encode(success)
	return

}

func (c *categoryAPI) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	getID := fmt.Sprintf("%s", r.Context().Value("id"))
	categoryID := r.URL.Query().Get("category_id")

	getIDInt, _ := strconv.Atoi(getID)
	categoryIDInt, _ := strconv.Atoi(categoryID)

	err := c.categoryService.DeleteCategory(r.Context(), categoryIDInt)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(200)
	success := map[string]interface{}{
		"user_id":     getIDInt,
		"category_id": categoryIDInt,
		"message":     "success delete category",
	}
	json.NewEncoder(w).Encode(success)
	return

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
