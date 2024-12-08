package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// User represents a user in the car-sharing system
type User struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	MembershipTier string `json:"membership_tier"`
}

type MembershipBenefits struct {
	DiscountRate   float64 `json:"discount_rate"`   // Discount rate for hourly rentals
	PriorityAccess bool    `json:"priority_access"` // Priority vehicle access
	IncreasedLimit int     `json:"increased_limit"` // Increased booking limits
}

var db *sql.DB

// Initialize MySQL connection
func initDB() {
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/user_management_db")
	if err != nil {
		log.Fatalf("Failed to connect to Mysql database: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping Mysql database: %v", err)
	}
	fmt.Println("Connected to the Mysql database.")
}

func getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, name, email, membership_tier FROM users")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.MembershipTier); err != nil {
			http.Error(w, "Error reading database", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}
	json.NewEncoder(w).Encode(users)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]

	var user User
	err := db.QueryRow("SELECT id, name, email, membership_tier FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Name, &user.Email, &user.MembershipTier)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO users (name, email, password, membership_tier) VALUES (?, ?, ?, ?)",
		user.Name, user.Email, user.Password, user.MembershipTier)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]

	var user struct {
		Name           string `json:"name"`
		Email          string `json:"email"`
		MembershipTier string `json:"membership_tier"`
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(
		"UPDATE users SET name = ?, email = ?, membership_tier = ? WHERE id = ?",
		user.Name, user.Email, user.MembershipTier, userId,
	)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]

	res, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var user struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	err := db.QueryRow("SELECT id, name, email FROM users WHERE email = ? AND password = ?", credentials.Email, credentials.Password).Scan(&user.ID, &user.Name, &user.Email)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      user.ID,
		"name":    user.Name,
		"email":   user.Email,
		"message": "Login successful",
	})
}

// View user profile by ID
func getUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]

	var user User
	err := db.QueryRow("SELECT id, name, email, membership_tier FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Name, &user.Email, &user.MembershipTier)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "User not found"})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Database error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func updateUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate input fields
	if user.Name == "" || user.Email == "" || user.Password == "" || user.MembershipTier == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	validTiers := map[string]bool{"Basic": true, "Premium": true, "VIP": true}
	if !validTiers[user.MembershipTier] {
		http.Error(w, "Invalid membership tier", http.StatusBadRequest)
		return
	}

	// Update user information in the database
	res, err := db.Exec(
		"UPDATE users SET name = ?, email = ?, membership_tier = ? WHERE id = ?",
		user.Name, user.Email, user.MembershipTier, id,
	)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			http.Error(w, "Email or name already in use", http.StatusConflict)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Return the updated user details without the password
	updatedUser := struct {
		Name           string `json:"name"`
		Email          string `json:"email"`
		MembershipTier string `json:"membership_tier"`
	}{
		Name:           user.Name,
		Email:          user.Email,
		MembershipTier: user.MembershipTier,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedUser)
}

func main() {
	initDB()
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/users", getAllUsersHandler).Methods("GET")
	router.HandleFunc("/api/v1/users/{id}", getUserHandler).Methods("GET")
	router.HandleFunc("/api/v1/users", createUserHandler).Methods("POST")
	router.HandleFunc("/api/v1/users/{id}", updateUserHandler).Methods("PUT")
	router.HandleFunc("/api/v1/users/{id}", deleteUserHandler).Methods("DELETE")
	router.HandleFunc("/api/v1/login", loginHandler).Methods("POST")
	router.HandleFunc("/api/v1/users/{id}", getUserProfileHandler).Methods("GET")
	router.HandleFunc("/api/v1/users/{id}", updateUserProfileHandler).Methods("PUT")

	corsHandler := handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedOrigins([]string{"*"}), // Replace "*" with the frontend origin if known
	)

	fmt.Println("Listening on port 5001")
	log.Fatal(http.ListenAndServe(":5001", corsHandler(router)))
}
