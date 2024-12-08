package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Vehicle struct {
	ID           int    `json:"id"`
	Make         string `json:"make"`
	Model        string `json:"model"`
	Availability bool   `json:"availability"`
}

var vehicleDB *sql.DB // Connection to the vehicles database
var userDB *sql.DB    // Connection to the users database

// Initialize the database connection
func initDB() {
	var err error
	vehicleDB, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/vehicle_reservation_db")
	if err != nil {
		log.Fatalf("Failed to connect to MySQL database: %v", err)
	}
	if err = vehicleDB.Ping(); err != nil {
		log.Fatalf("Failed to ping MySQL database: %v", err)
	}

	userDB, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/user_management_db")
	if err != nil {
		log.Fatalf("Failed to connect to user database: %v", err)
	}
	if err := userDB.Ping(); err != nil {
		log.Fatalf("Failed to ping user database: %v", err)
	}

	fmt.Println("Connected to the user and vehicle databases.")
}

// CRUD Handlers for Vehicles

// Get all vehicles
func getVehicles(w http.ResponseWriter, r *http.Request) {
	vehicles := []Vehicle{}
	rows, err := vehicleDB.Query("SELECT id, make, model, availability FROM vehicles")
	if err != nil {
		http.Error(w, "Failed to fetch vehicles", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var v Vehicle
		if err := rows.Scan(&v.ID, &v.Make, &v.Model, &v.Availability); err != nil {
			http.Error(w, "Failed to parse vehicles", http.StatusInternalServerError)
			return
		}
		vehicles = append(vehicles, v)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(vehicles)
}

// Get a single vehicle by ID
func getVehicle(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var v Vehicle
	err := vehicleDB.QueryRow("SELECT id, make, model, availability FROM vehicles WHERE id = ?", id).Scan(
		&v.ID, &v.Make, &v.Model, &v.Availability,
	)
	if err == sql.ErrNoRows {
		http.Error(w, "Vehicle not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to fetch vehicle", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(v)
}

// Create a new vehicle
func createVehicle(w http.ResponseWriter, r *http.Request) {
	var v Vehicle
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	res, err := vehicleDB.Exec("INSERT INTO vehicles (make, model, availability) VALUES (?, ?, ?)", v.Make, v.Model, v.Availability)
	if err != nil {
		http.Error(w, "Failed to create vehicle", http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()
	v.ID = int(id)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(v)
}

// Update an existing vehicle
func updateVehicle(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var v Vehicle
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	res, err := vehicleDB.Exec("UPDATE vehicles SET make = ?, model = ?, availability = ? WHERE id = ?", v.Make, v.Model, v.Availability, id)
	if err != nil {
		http.Error(w, "Failed to update vehicle", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Vehicle not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Vehicle updated successfully"})
}

// Delete a vehicle
func deleteVehicle(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	res, err := vehicleDB.Exec("DELETE FROM vehicles WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Failed to delete vehicle", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Vehicle not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Vehicle deleted successfully"})
}

func getAvailableVehiclesForUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Assuming userId is passed as a query parameter
	userId := r.URL.Query().Get("userId")
	if userId == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Query to fetch available vehicles
	rows, err := vehicleDB.Query(`
        SELECT id, make, model, availability 
        FROM vehicles 
        WHERE availability = TRUE
    `)
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var vehicles []Vehicle
	for rows.Next() {
		var vehicle Vehicle
		if err := rows.Scan(&vehicle.ID, &vehicle.Make, &vehicle.Model, &vehicle.Availability); err != nil {
			http.Error(w, "Error scanning vehicle data", http.StatusInternalServerError)
			return
		}
		vehicles = append(vehicles, vehicle)
	}

	// Return the list of vehicles as JSON
	json.NewEncoder(w).Encode(vehicles)
}

func getAvailableVehiclesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := vehicleDB.Query("SELECT id, make, model, availability FROM vehicles WHERE availability = TRUE")
	if err != nil {
		http.Error(w, "Failed to fetch available vehicles", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var vehicles []Vehicle
	for rows.Next() {
		var vehicle Vehicle
		if err := rows.Scan(&vehicle.ID, &vehicle.Make, &vehicle.Model, &vehicle.Availability); err != nil {
			http.Error(w, "Failed to parse vehicle data", http.StatusInternalServerError)
			return
		}
		vehicles = append(vehicles, vehicle)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(vehicles)
}

func createReservation(w http.ResponseWriter, r *http.Request) {
	var input struct {
		VehicleID int    `json:"vehicle_id"`
		UserID    int    `json:"user_id"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Ensure the vehicle is available for the requested time
	query := `
        SELECT COUNT(*) FROM reservations
        WHERE vehicle_id = ? AND status = 'active'
        AND ((start_time BETWEEN ? AND ?) OR (end_time BETWEEN ? AND ?))
    `
	var count int
	err := vehicleDB.QueryRow(query, input.VehicleID, input.StartTime, input.EndTime, input.StartTime, input.EndTime).Scan(&count)
	if err != nil || count > 0 {
		http.Error(w, "Vehicle is not available for the requested time", http.StatusConflict)
		return
	}

	res, err := vehicleDB.Exec("INSERT INTO reservations (vehicle_id, user_id, start_time, end_time) VALUES (?, ?, ?, ?)",
		input.VehicleID, input.UserID, input.StartTime, input.EndTime)
	if err != nil {
		http.Error(w, "Failed to create reservation", http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"reservation_id": int(id)})
}

func checkIfUserExists(userID int) (bool, error) {
	var count int
	err := userDB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func createReservations(vehicleID, userID int, startTime, endTime string) error {
	// First, check if the user exists
	userExists, err := checkIfUserExists(userID)
	if err != nil {
		return fmt.Errorf("error checking if user exists: %v", err)
	}
	if !userExists {
		return fmt.Errorf("user with ID %d does not exist", userID)
	}

	// If user exists, proceed with creating the reservation
	_, err = vehicleDB.Exec("INSERT INTO reservations (vehicle_id, user_id, start_time, end_time, status) VALUES (?, ?, ?, ?, 'active')", vehicleID, userID, startTime, endTime)
	if err != nil {
		return fmt.Errorf("failed to create reservation: %v", err)
	}

	fmt.Println("Reservation created successfully.")
	return nil
}

func getReservationsByUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Query to get reservations for the user
	rows, err := vehicleDB.Query(`
        SELECT 
            r.id, v.make, v.model, r.start_time, r.end_time, r.status
        FROM 
            reservations r
        JOIN 
            vehicles v ON r.vehicle_id = v.id
        WHERE 
            r.user_id = ?`, userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var reservations []struct {
		ID        int    `json:"id"`
		Vehicle   string `json:"vehicle"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
		Status    string `json:"status"`
	}

	for rows.Next() {
		var reservation struct {
			ID        int
			Make      string
			Model     string
			StartTime string
			EndTime   string
			Status    string
		}
		if err := rows.Scan(&reservation.ID, &reservation.Make, &reservation.Model, &reservation.StartTime, &reservation.EndTime, &reservation.Status); err != nil {
			http.Error(w, "Error scanning data", http.StatusInternalServerError)
			return
		}
		reservations = append(reservations, struct {
			ID        int    `json:"id"`
			Vehicle   string `json:"vehicle"`
			StartTime string `json:"start_time"`
			EndTime   string `json:"end_time"`
			Status    string `json:"status"`
		}{
			ID:        reservation.ID,
			Vehicle:   reservation.Make + " " + reservation.Model,
			StartTime: reservation.StartTime,
			EndTime:   reservation.EndTime,
			Status:    reservation.Status,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservations)
}

func createReservationHandler(w http.ResponseWriter, r *http.Request) {
	// Parse input data
	var input struct {
		VehicleID int    `json:"vehicle_id"`
		UserID    int    `json:"user_id"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if user exists in the user database
	var userCount int
	err := userDB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", input.UserID).Scan(&userCount)
	if err != nil || userCount == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if the vehicle is available in the vehicle database
	var vehicleAvailable bool
	err = vehicleDB.QueryRow(`
		SELECT availability 
		FROM vehicles 
		WHERE id = ?`, input.VehicleID).Scan(&vehicleAvailable)
	if err != nil {
		http.Error(w, "Error fetching vehicle availability", http.StatusInternalServerError)
		return
	}
	if !vehicleAvailable {
		http.Error(w, "Vehicle is not available", http.StatusConflict)
		return
	}

	// Insert reservation into reservations table
	_, err = vehicleDB.Exec(`
		INSERT INTO reservations (vehicle_id, user_id, start_time, end_time, status)
		VALUES (?, ?, ?, ?, 'active')`,
		input.VehicleID, input.UserID, input.StartTime, input.EndTime)
	if err != nil {
		http.Error(w, "Failed to create reservation", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Reservation created successfully"))
}

func updateReservationHandler(w http.ResponseWriter, r *http.Request) {
	reservationId := mux.Vars(r)["id"]

	var reservation struct {
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	_, err := vehicleDB.Exec("UPDATE reservations SET start_time = ?, end_time = ? WHERE id = ?", reservation.StartTime, reservation.EndTime, reservationId)
	if err != nil {
		http.Error(w, "Failed to update reservation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Reservation updated successfully"})
}

func cancelReservation(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	res, err := vehicleDB.Exec("UPDATE reservations SET status = 'cancelled' WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Failed to cancel reservation", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Reservation not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Reservation cancelled successfully"})
}

func checkUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	userExists, err := checkIfUserExists(userIDInt)
	if err != nil {
		http.Error(w, "Error checking user existence", http.StatusInternalServerError)
		return
	}

	if userExists {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User exists"))
	} else {
		http.Error(w, "User not found", http.StatusNotFound)
	}
}

func checkVehicleAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	vehicleID := r.URL.Query().Get("vehicle_id")
	startTime := r.URL.Query().Get("start_time")
	endTime := r.URL.Query().Get("end_time")

	// Convert vehicleID to int
	vehicleIDInt, err := strconv.Atoi(vehicleID)
	if err != nil {
		http.Error(w, "Invalid vehicle ID", http.StatusBadRequest)
		return
	}

	// Check vehicle availability
	available, err := checkVehicleAvailability(vehicleIDInt, startTime, endTime)
	if err != nil {
		http.Error(w, "Error checking vehicle availability", http.StatusInternalServerError)
		return
	}

	if available {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Vehicle is available"))
	} else {
		http.Error(w, "Vehicle not available", http.StatusConflict)
	}
}

func checkVehicleAvailability(vehicleID int, startTime, endTime string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM reservations 
              WHERE vehicle_id = ? AND status = 'active' 
              AND (start_time < ? AND end_time > ?)`
	err := vehicleDB.QueryRow(query, vehicleID, endTime, startTime).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	initDB()
	defer vehicleDB.Close()
	defer userDB.Close()

	router := mux.NewRouter()

	// Vehicle routes
	router.HandleFunc("/vehicles", getVehicles).Methods("GET")
	router.HandleFunc("/vehicles/{id}", getVehicle).Methods("GET")
	router.HandleFunc("/vehicles", createVehicle).Methods("POST")
	router.HandleFunc("/vehicles/{id}", updateVehicle).Methods("PUT")
	router.HandleFunc("/vehicles/{id}", deleteVehicle).Methods("DELETE")

	router.HandleFunc("/reservations", createReservation).Methods("POST")
	router.HandleFunc("/reservations/{id}", updateReservationHandler).Methods("PUT")
	router.HandleFunc("/reservations/{id}", cancelReservation).Methods("DELETE")

	router.HandleFunc("/api/v1/vehicles/available", getAvailableVehiclesHandler).Methods("GET")
	router.HandleFunc("/api/v1/vehicles/available", getAvailableVehiclesForUserHandler).Methods("GET")

	http.HandleFunc("/check-user", checkUserHandler)
	http.HandleFunc("/check-vehicle-availability", checkVehicleAvailabilityHandler)
	http.HandleFunc("/create-reservation", createReservationHandler)

	router.HandleFunc("/api/reservations", getReservationsByUserHandler).Methods("GET")

	// CORS handling
	corsHandler := handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedOrigins([]string{"*"}), // Replace "*" with the frontend origin if known
	)

	fmt.Println("Listening on port 5000")
	log.Fatal(http.ListenAndServe(":5000", corsHandler(router)))
	log.Fatal(http.ListenAndServe(":5000", enableCORS(router)))
}
