	package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/Trentham3269/cattledog/models"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

// Define database for global access
var db *gorm.DB

func main() {
	var err error

	// Load .env file
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Assign environment variables
	pg_host := os.Getenv("HOST")
	pg_port := os.Getenv("PORT")
	pg_user := os.Getenv("USERNAME")
	pg_dbname := os.Getenv("DATABASE")

	// Postgres connection string 
	pgInfo := fmt.Sprintf("host=%s port=%s user=%s "+
    "dbname=%s sslmode=disable",
    pg_host, pg_port, pg_user, pg_dbname)

    // Connect to database
	db, err = gorm.Open("postgres", pgInfo)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()			

	// Create url routes
	r := mux.NewRouter()

	// Public routes
	r.HandleFunc("/signup", createUser).Methods("POST")
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/categories", getCategories).Methods("GET")
	r.HandleFunc("/categories/{id}", getCategory).Methods("GET")
	
	// Auth routes
	s := r.PathPrefix("/auth").Subrouter()
	s.HandleFunc("/categories", addCategory).Methods("POST")
	s.HandleFunc("/categories/{id}", updateCategory).Methods("PUT")
	s.HandleFunc("/categories/{id}", deleteCategory).Methods("DELETE")

	// Start http server
	http_port := 8888
	http_host := fmt.Sprintf("localhost:%d", http_port)
	log.Println(fmt.Sprintf("Server listening on %d...", http_port))
	log.Fatal(http.ListenAndServe(http_host, r))
}

func createUser(w http.ResponseWriter, r *http.Request) {
	// Decode request to retrieve password
	user := models.User{}
	json.NewDecoder(r.Body).Decode(&user)

	// Encrypt password
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Password encryption failed")
	}

	// Assign new password
	user.Password = string(pass)
	
	// Create user in db
	db.Create(&user)

	// Log user created
	log.Println(fmt.Sprintf("Create user: %s", user.Email))
}

var (
    // key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
    key = []byte("super-secret-key")
    store = sessions.NewCookieStore(key)
)

func login(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "cookie-name")

    // Decode request
	user := models.User{}
	json.NewDecoder(r.Body).Decode(&user)

	// Grab plain text password
	password := user.Password
	
	// Check database for user
	err := db.
		Where("Email = ?", user.Email).
		First(&user).
		Error
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Email address not found"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Compare hashed password to plain text password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) 
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Invalid credentials...please try again"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Set header and encode as json
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var resp = map[string]interface{}{"status": false, "message": "User is now logged in"}
	json.NewEncoder(w).Encode(resp)

    // Set user as authenticated
    session.Values["authenticated"] = true
    session.Save(r, w)
}

func getCategories(w http.ResponseWriter, r *http.Request) {
	// Query db and return all categories
	categories := []models.Category{}
	db.Find(&categories)

	// Log endpoint
	log.Println("Return all categories")

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

	// Log category returned
	log.Println(fmt.Sprintf("Return category of %s", category.Name))

	// Set header and encode as json
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(category)
}

func addCategory(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	session, _ := store.Get(r, "cookie-name")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Decode request and create record in db
	category := models.Category{}
	json.NewDecoder(r.Body).Decode(&category)
	db.Create(&category)

	// Log new category
	log.Println(fmt.Sprintf("Create %s category", category.Name))
}

func updateCategory(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	session, _ := store.Get(r, "cookie-name")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Access url parameter
	vars := mux.Vars(r)

	// Find correct record by id
	category := models.Category{}
	db.First(&category, vars["id"])
	
	// Parse request body
	json.NewDecoder(r.Body).Decode(&category)
	db.Save(&category.Name)

	// Log category updated
	log.Println(fmt.Sprintf("Update category to %s", category.Name))
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	session, _ := store.Get(r, "cookie-name")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Access url parameter
	vars := mux.Vars(r)

	// Find correct record and delete from db
	category := models.Category{}
	db.First(&category, vars["id"])
	db.Delete(&category)

	// Log category deleted
	log.Println(fmt.Sprintf("Delete %s category", category.Name))
}