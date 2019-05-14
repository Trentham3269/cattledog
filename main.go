package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/subosito/gotenv"
)

type Category struct {
	ID		uint	`gorm:"primary_key" json:"id"`
	Name	string	`gorm:"size:50" json:"name"` 
}

// Define database for global access
var db *gorm.DB

// Load environment variables
func init() {
	gotenv.Load()
}

func main() {
	var err error
	pgUrl, err := pq.ParseURL(os.Getenv("ELEPHANTSQL_URL"))
	if err != nil {
		fmt.Println(err)
	}

	db, err = gorm.Open("postgres", pgUrl)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// Create url routes
	router := mux.NewRouter()
	router.HandleFunc("/categories", getCategories).Methods("GET")
	router.HandleFunc("/categories/{id}", getCategory).Methods("GET")
	router.HandleFunc("/categories", addCategory).Methods("POST")
	router.HandleFunc("/categories", updateCategory).Methods("PUT")
	router.HandleFunc("/categories/{id}", deleteCategory).Methods("DELETE")

	// Start http server
	port := 8888
	host := fmt.Sprintf("localhost:%d", port)
	log.Println(fmt.Sprintf("Server listening on %d...", port))
	log.Fatal(http.ListenAndServe(host, router))
}

func getCategories(w http.ResponseWriter, r *http.Request) {
	// Query db and return all categories
	categories := []Category{}
	db.Find(&categories)

	// Log endpoint
	log.Println("Return all categories")
	
	// Set header and encode as json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func getCategory(w http.ResponseWriter, r *http.Request) {
	// Access url parameter
	vars := mux.Vars(r)

	// Query db and return category by id
	category := Category{}
	db.Where("ID = ?", vars["id"]).Find(&category)

	// Log category returned
	log.Println(fmt.Sprintf("Return category of %s", category.Name))

	// Set header and encode as json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func addCategory(w http.ResponseWriter, r *http.Request) {

}

func updateCategory(w http.ResponseWriter, r *http.Request) {

}

func deleteCategory(w http.ResponseWriter, r *http.Request) {

}
