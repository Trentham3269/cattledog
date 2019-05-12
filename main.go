package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	//"strconv"

	"github.com/lib/pq"
	"github.com/subosito/gotenv"
	"github.com/gorilla/mux"
)

type Category struct {
	ID	int	`json:id`
	Name	string	`json:name`
}

var categories []Category

var db *sql.DB

func init() {
	gotenv.Load()
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	pgUrl, err := pq.ParseURL(os.Getenv("ELEPHANTSQL_URL"))
	logFatal(err)

	db, err = sql.Open("postgres", pgUrl)
	logFatal(err)

	err = db.Ping()
	if err != nil {
		log.Println("Could not connect to database")
	} else {
		log.Println("Connected to database successfully")
	}
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/categories", getCategories).Methods("GET")
	router.HandleFunc("/categories/{id}", getCategory).Methods("GET")
	router.HandleFunc("/categories", addCategory).Methods("POST")
	router.HandleFunc("/categories", updateCategory).Methods("PUT")
	router.HandleFunc("/categories/{id}", deleteCategory).Methods("DELETE")

	port := 8888
	host := fmt.Sprintf("localhost:%d", port)
	log.Println(fmt.Sprintf("Server listening on %d...", port))
	log.Fatal(http.ListenAndServe(host, router))
}

func getCategories(w http.ResponseWriter, r *http.Request) {
	var category Category
	categories = []Category{}

	rows, err := db.Query("select * from categories")
	logFatal(err)
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&category.ID, &category.Name)
		logFatal(err)

		categories = append(categories, category)
	}

	json.NewEncoder(w).Encode(categories)
}

func getCategory(w http.ResponseWriter, r *http.Request) {
	var category Category
	params := mux.Vars(r)

	row := db.QueryRow("select * from categories where id=$1", params["id"])

	err := row.Scan(&category.ID, &category.Name)
	logFatal(err)

	json.NewEncoder(w).Encode(category)
}


func addCategory(w http.ResponseWriter, r *http.Request) {

}

func updateCategory(w http.ResponseWriter, r *http.Request) {

}

func deleteCategory(w http.ResponseWriter, r *http.Request) {

}
