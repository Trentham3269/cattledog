package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/Trentham3269/cattledog/config"
	"github.com/Trentham3269/cattledog/middleware"
	"github.com/Trentham3269/cattledog/models"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

// Define database for global access
var db *gorm.DB

func main() {
	var err error

	cf := GetConfig()

	// Postgres connection string 
	pgInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"dbname=%s sslmode=disable",
		cf.Host, cf.Port, cf.Username, cf.Database)

	// Connect to database
	db, err = gorm.Open("postgres", pgInfo)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()			

	// Create url routes
	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)

	// Public routes
	r.HandleFunc("/signup", createUser).Methods("POST")
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/logout", logout).Methods("GET")
	r.HandleFunc("/categories", getCategories).Methods("GET")
	r.HandleFunc("/categories/{id}", getCategory).Methods("GET")

	// Auth routes
	s := r.PathPrefix("/auth").Subrouter()
	s.Use(middleware.SessionMiddleware)
	s.HandleFunc("/categories", addCategory).Methods("POST")
	s.HandleFunc("/categories/{id}", updateCategory).Methods("PUT")
	s.HandleFunc("/categories/{id}", deleteCategory).Methods("DELETE")

	// Start http server
	http_port := 8888
	http_host := fmt.Sprintf("localhost:%d", http_port)
	log.Println(fmt.Sprintf("Server listening on %d...", http_port))
	log.Fatal(http.ListenAndServe(http_host, r))
}

func GetConfig() config.Config {
	var err error
	// Load .env file
	err = godotenv.Load()
	if err != nil {
		log.Fatal(fmt.Println("Could not load .env file"))
	}

	// Populate config struct
	cf := config.Config{
		Host: os.Getenv("HOST"),
		Port: os.Getenv("PORT"),
		Username: os.Getenv("USERNAME"),
		Database: os.Getenv("DATABASE"),
	}

	// Make struct available to db
	return cf
}

func createUser(w http.ResponseWriter, r *http.Request) {
	// Decode request to retrieve password
	user := models.User{}
	json.NewDecoder(r.Body).Decode(&user)

	// Encrypt password
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Password encryption failed")
	}

	// Assign new password
	user.Password = string(pass)
	
	// Create user in db
	db.Create(&user)
}

func login(w http.ResponseWriter, r *http.Request) {
	session, err := middleware.Store.Get(r, "cookie-name")
	if err != nil {
		log.Println(err)
	}

	// Decode request
	user := models.User{}
	json.NewDecoder(r.Body).Decode(&user)

	// Grab plain text password
	password := user.Password
	
	// Check database for user
	err = db.
		Where("Email = ?", user.Email).
		First(&user).
		Error
	if err != nil {
		var resp = map[string]interface{}{"message": "Email address not found"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Compare hashed password to plain text password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) 
	if err != nil {
		var resp = map[string]interface{}{"message": "Invalid credentials...please try again"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Set user as authenticated
	session.Values["authenticated"] = true
	session.Save(r, w)

	// Set header and encode as json
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var resp = map[string]interface{}{"message": "User is now logged in"}
	json.NewEncoder(w).Encode(resp)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, err := middleware.Store.Get(r, "cookie-name")
	if err != nil {
		log.Println(err)
	}

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)

	// Set header and encode as json
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var resp = map[string]interface{}{"message": "User is now logged out"}
	json.NewEncoder(w).Encode(resp)
}

func getCategories(w http.ResponseWriter, r *http.Request) {
	// Query db and return all categories
	categories := []models.Category{}
	db.Find(&categories)

	// Set header and encode as json
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(categories)
}

func getCategory(w http.ResponseWriter, r *http.Request) {
	// Access url parameter
	vars := mux.Vars(r)

	// Query db and return category by id
	category := models.Category{}
	db.
		Preload("Items").
		Preload("Items.User").
		Where("ID = ?", vars["id"]).
		Find(&category)

	// Set header and encode as json
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(category)
}

func addCategory(w http.ResponseWriter, r *http.Request) {
	// Decode request and create record in db
	category := models.Category{}
	json.NewDecoder(r.Body).Decode(&category)
	db.Create(&category)
}

func updateCategory(w http.ResponseWriter, r *http.Request) {
	// Access url parameter
	vars := mux.Vars(r)

	// Find correct record by id
	category := models.Category{}
	db.First(&category, vars["id"])
	
	// Parse request body
	json.NewDecoder(r.Body).Decode(&category)
	db.Save(&category.Name)
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	// Access url parameter
	vars := mux.Vars(r)

	// Find correct record and delete from db
	category := models.Category{}
	db.First(&category, vars["id"])
	db.Delete(&category)
}