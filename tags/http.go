package tags

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

var TagsInMemory []Tag

type TagService struct {
	Repository TagRepository
}

type Tag struct {
	ID   uint64 `json:"ID"`
	Name string `json:"name"`
}

func (s TagService) CreateTag(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error trying to read from body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	tag := Tag{}

	if err := json.Unmarshal(body, &tag); err != nil {
		log.Println("Error trying to unmarshal: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.Repository.Add(tag)
	w.WriteHeader(http.StatusCreated)
}

type GetAllTagsResponse struct {
	Tags []Tag `json:"tags"`
}

func (s TagService) ListTags(w http.ResponseWriter, r *http.Request) {
	tags, err := s.Repository.GetAll()
	if err != nil {
		log.Println("Failed to fetch tags: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := GetAllTagsResponse{
		Tags: tags,
	}

	respJSON, err := json.Marshal(resp)
	if err != nil {
		log.Println("Error trying to marshal JSON: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(respJSON)
}
