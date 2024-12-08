a) Design consideration of my microservices

1. Modular Architecture
Each service has a distinct responsibility:
User Management Service: Manages user profiles and membership tiers.
Vehicle Reservation Service: Handles vehicle availability, booking, and cancellation.
Billing Service: Calculates costs, applies membership discounts, and processes payments.

3. Database Separation
Each service has its own database to ensure data independence and scalability:
user_management_db: Stores user-related data.
vehicle_reservation_db: Stores vehicle and reservation records.
billing_db: Maintains billing and payment details.

5. Communication
Services communicate via RESTful APIs using JSON as the standard data format.
Services are loosely coupled to reduce dependencies.

7. Scalability
Services can be deployed and scaled independently based on the load.

9. Error Handling and Fault Tolerance
Graceful error handling is implemented in API endpoints.
Each service is designed to handle failures without impacting the overall system

b) Instructions for Setting Up and Running Microservices
Pre-requirements
Install the following dependencies:

Go (version 1.19 or later)
MySQL (version 8.0 or later)



Step 1: Configure Databases
Create three databases:
user_management_db
vehicle_reservation_db
billing_db



Step 2: Run Microservices Locally
Start each service in its respective directory:

To access User Management Service:

cd User_Management
go run main.go

To access Vehicle Reservation Service:

Copy code
cd Vehicle_Management
go run vehicles.go

To access Billing Service:

Copy code
cd Billing_Management
go run billing.go
The services will be available at:

User Management Service: http://localhost:5001
Vehicle Reservation Service: http://localhost:5000
Billing Service: http://localhost:5002

Step 3: Test the Setup
Use a tool like Postman or curl to interact with the APIs.
Example: Fetch user data from User Management Service:

Copy code:
curl -X GET http://localhost:5001/users/1

Step 4: Frontend Integration
Place the frontend/ directory in a web server.

Configure the API endpoints in frontend/js/config.js:

javascript
Copy code

const baseURL = "http://localhost:5001/api/v1/users"  For the users
const apiUrl = "http://localhost:5000/vehicles";   For the vehicles
const baseURL = "http://localhost:5002";     For the billing

Launch the frontend and test the integrated application.
