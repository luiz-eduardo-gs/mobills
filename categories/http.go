package categories

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type CategoryType uint8

const (
	Income CategoryType = iota
	Expense
)

type Category struct {
	ID   uint64       `json:"ID"`
	Name string       `json:"name"`
	Type CategoryType `json:"type"`
}

type CategoryService struct {
	Repository CategoryRepository
}

func logError(w http.ResponseWriter, err error) {
	log.Println("Error: ", err)
	w.WriteHeader(http.StatusInternalServerError)
}

func (s CategoryService) CreateCategory(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to read body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	c := Category{}

	if err := json.Unmarshal(body, &c); err != nil {
		log.Println("Error unmarshal: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := s.Repository.Add(c); err != nil {
		log.Println("Error trying to save category: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type ListCategoriesResp struct {
	Categories []Category `json:"categories"`
}

func (s CategoryService) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := s.Repository.GetAll()
	if err != nil {
		logError(w, err)
		return
	}

	resp := ListCategoriesResp{
		Categories: categories,
	}

	respJSON, err := json.Marshal(resp)
	if err != nil {
		logError(w, err)
		return
	}

	w.Write(respJSON)
}
