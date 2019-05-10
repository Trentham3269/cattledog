package main

import (
	"encoding/json"
	//"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Category struct {
	ID	int	`json:id`
	Name	string	`json:name`
}

var categories []Category

func main() {
	router := mux.NewRouter()

	categories = append(categories,
		Category{ID: 1, Name: "Soccer"},
		Category{ID: 2, Name: "Basketball"},
		Category{ID: 3, Name: "Baseball"},
		Category{ID: 4, Name: "Frisbee"},
		Category{ID: 5, Name: "Snowboarding"},
		Category{ID: 6, Name: "Rock Climbing"},
		Category{ID: 7, Name: "Football"},
		Category{ID: 8, Name: "Surfing"},
		Category{ID: 9, Name: "Hockey"})

	router.HandleFunc("/categories", getCategories).Methods("GET")
	router.HandleFunc("/categories/{id}", getCategory).Methods("GET")

	log.Fatal(http.ListenAndServe(":8888", router))
}

func getCategories(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(categories)
}

func getCategory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	i, _ := strconv.Atoi(params["id"])

	for _, category := range categories {
		if category.ID == i {
			json.NewEncoder(w).Encode(&category)
		}
	}
}
