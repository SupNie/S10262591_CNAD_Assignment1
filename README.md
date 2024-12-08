a) Design consideration of my microservices
1. Modular Architecture
Each service has a distinct responsibility:
User Management Service: Manages user profiles and membership tiers.
Vehicle Reservation Service: Handles vehicle availability, booking, and cancellation.
Billing Service: Calculates costs, applies membership discounts, and processes payments.
2. Database Separation
Each service has its own database to ensure data independence and scalability:
user_management_db: Stores user-related data.
vehicle_reservation_db: Stores vehicle and reservation records.
billing_db: Maintains billing and payment details.
3. Communication
Services communicate via RESTful APIs using JSON as the standard data format.
Services are loosely coupled to reduce dependencies.
4. Scalability
Services can be deployed and scaled independently based on the load.
5. Error Handling and Fault Tolerance
Graceful error handling is implemented in API endpoints.
Each service is designed to handle failures without impacting the overall system
