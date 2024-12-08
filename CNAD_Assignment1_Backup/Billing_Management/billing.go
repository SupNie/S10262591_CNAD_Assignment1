package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Billing struct {
	ID            int     `json:"id"`
	ReservationID int     `json:"reservation_id"`
	Amount        float64 `json:"amount"`
	PaymentStatus string  `json:"payment_status"`
}

type Invoice struct {
	InvoiceID     int       `json:"invoice_id"`
	BillingID     int       `json:"billing_id"`
	ReservationID int       `json:"reservation_id"`
	Amount        float64   `json:"amount"`
	GeneratedDate time.Time `json:"generated_date"`
}

type Receipt struct {
	ReceiptID   int       `json:"receipt_id"`
	BillingID   int       `json:"billing_id"`
	Amount      float64   `json:"amount"`
	PaymentDate time.Time `json:"payment_date"`
}

var vehicleDB *sql.DB
var userDB *sql.DB
var billingDB *sql.DB // Billing database connection

func initDB() {
	var err error
	// Connect to the vehicle reservation database
	vehicleDB, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/vehicle_reservation_db")
	if err != nil {
		log.Fatalf("Failed to connect to vehicle reservation database: %v", err)
	}
	if err := vehicleDB.Ping(); err != nil {
		log.Fatalf("Failed to ping vehicle reservation database: %v", err)
	}

	// Connect to the user management database
	userDB, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/user_management_db")
	if err != nil {
		log.Fatalf("Failed to connect to user database: %v", err)
	}
	if err := userDB.Ping(); err != nil {
		log.Fatalf("Failed to ping user database: %v", err)
	}

	// Connect to the billing database
	billingDB, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/billingpayment_db")
	if err != nil {
		log.Fatalf("Failed to connect to billing database: %v", err)
	}
	if err := billingDB.Ping(); err != nil {
		log.Fatalf("Failed to ping billing database: %v", err)
	}

	fmt.Println("Connected to all databases.")
}

func generateInvoice(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	billingID := params["billing_id"]

	// Get the billing details from the billing database
	var billing Billing
	err := billingDB.QueryRow("SELECT id, reservation_id, amount, payment_status FROM billings WHERE id = ?", billingID).
		Scan(&billing.ID, &billing.ReservationID, &billing.Amount, &billing.PaymentStatus)
	if err == sql.ErrNoRows {
		http.Error(w, "Billing not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to fetch billing data", http.StatusInternalServerError)
		return
	}

	// Get reservation details from vehicle reservation database
	var reservationDetails struct {
		VehicleID int
		UserID    int
	}
	err = vehicleDB.QueryRow("SELECT vehicle_id, user_id FROM reservations WHERE id = ?", billing.ReservationID).
		Scan(&reservationDetails.VehicleID, &reservationDetails.UserID)
	if err != nil {
		http.Error(w, "Failed to fetch reservation details", http.StatusInternalServerError)
		return
	}

	// Get user details from user management database
	var userName string
	err = userDB.QueryRow("SELECT name FROM users WHERE id = ?", reservationDetails.UserID).Scan(&userName)
	if err != nil {
		http.Error(w, "Failed to fetch user details", http.StatusInternalServerError)
		return
	}

	// Generate the invoice
	invoice := Invoice{
		InvoiceID:     billing.ID,
		BillingID:     billing.ID,
		ReservationID: billing.ReservationID,
		Amount:        billing.Amount,
		GeneratedDate: time.Now(),
	}

	// Return the invoice as a response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoice)
}

func generateReceipt(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	billingID := params["billing_id"]

	// Get the billing details
	var billing Billing
	err := billingDB.QueryRow("SELECT id, amount, payment_status FROM billings WHERE id = ?", billingID).
		Scan(&billing.ID, &billing.Amount, &billing.PaymentStatus)
	if err == sql.ErrNoRows {
		http.Error(w, "Billing not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to fetch billing data", http.StatusInternalServerError)
		return
	}

	if billing.PaymentStatus != "Paid" {
		http.Error(w, "Payment not completed for this billing", http.StatusBadRequest)
		return
	}

	// Generate receipt
	receipt := Receipt{
		ReceiptID:   billing.ID,
		BillingID:   billing.ID,
		Amount:      billing.Amount,
		PaymentDate: time.Now(),
	}

	// Return receipt
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(receipt)
}

func getVehiclePricing(vehicleType string) (float64, float64, float64, float64, error) {
	var baseRate, discountBasic, discountPremium, discountVIP float64
	err := billingDB.QueryRow("SELECT base_rate_per_hour, discount_basic, discount_premium, discount_vip FROM vehicle_pricing WHERE vehicle_type = ?", vehicleType).
		Scan(&baseRate, &discountBasic, &discountPremium, &discountVIP)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	return baseRate, discountBasic, discountPremium, discountVIP, nil
}

func calculateCost(vehicleType string, membershipLevel string, startTime, endTime time.Time) (float64, error) {
	baseRate, discountBasic, discountPremium, discountVIP, err := getVehiclePricing(vehicleType)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch vehicle pricing: %v", err)
	}

	// Calculate rental duration in hours
	rentalDuration := endTime.Sub(startTime).Hours()

	// Determine the discount based on membership level
	discount := 0.0
	switch membershipLevel {
	case "Basic":
		discount = discountBasic
	case "Premium":
		discount = discountPremium
	case "VIP":
		discount = discountVIP
	}

	// Calculate the total cost
	totalCost := rentalDuration * baseRate * (1 - discount/100)
	return totalCost, nil
}

func main() {
	initDB()
	defer billingDB.Close()
	defer vehicleDB.Close()
	defer userDB.Close()

	router := mux.NewRouter()
	router.HandleFunc("/invoices/{billing_id}", generateInvoice).Methods("GET")
	router.HandleFunc("/receipts/{billing_id}", generateReceipt).Methods("GET")

	// Configure CORS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Replace "*" with the frontend origin if needed
		handlers.AllowedMethods([]string{"GET", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	fmt.Println("Billing and payment processing service started on port 5002")
	log.Fatal(http.ListenAndServe(":5002", corsHandler(router)))
}
