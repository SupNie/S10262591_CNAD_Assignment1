const vehicleList = document.getElementById("vehicleItems"); // Reference to the <ul> element for vehicles
const apiUrl = "http://localhost:5000/api/v1/vehicles/available"; // Endpoint for available vehicles

// Fetch available vehicles
function fetchAvailableVehicles() {
    console.log("Fetching available vehicles...");
    fetch(apiUrl)
        .then(response => {
            if (!response.ok) {
                console.error(`Error fetching vehicles: ${response.status}`);
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(vehicles => {
            console.log("Fetched vehicles:", vehicles); // Debug: check the vehicle data
            vehicleList.innerHTML = ""; // Clear the list before updating
            if (vehicles.length === 0) {
                vehicleList.innerHTML = "<li>No vehicles available at the moment.</li>";
                return;
            }
            vehicles.forEach(vehicle => {
                const li = document.createElement("li");
                li.textContent = `ID: ${vehicle.id}, Make: ${vehicle.make}, Model: ${vehicle.model}, Availability: ${
                    vehicle.availability ? "Available" : "Not Available"
                }`;
                vehicleList.appendChild(li);
            });
        })
        .catch(err => {
            console.error("Error fetching available vehicles:", err);
        });
}

// Reservation Form Handler
document.getElementById("reservationForm").addEventListener("submit", function (event) {
    event.preventDefault();

    const vehicleId = document.getElementById("vehicleId").value;
    const startTime = document.getElementById("startTime").value;
    const endTime = document.getElementById("endTime").value;
    const userId = localStorage.getItem("userId");

    console.log("Submitting reservation with the following data:");
    console.log(`User ID: ${userId}, Vehicle ID: ${vehicleId}, Start Time: ${startTime}, End Time: ${endTime}`);

    if (!userId) {
        alert("User not logged in");
        console.log("Error: User not logged in");
        return;
    }

    const reservationData = { vehicle_id: vehicleId, user_id: userId, start_time: startTime, end_time: endTime };

    const reservationApiUrl = "http://localhost:5001/api/v1/reservations";

    fetch(reservationApiUrl, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(reservationData),
    })
        .then(response => {
            if (!response.ok) {
                console.error(`Error creating reservation: ${response.status}`);
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            console.log("Reservation successful:", data); // Debug: check the reservation response
            alert(`Reservation successful! Reservation ID: ${data.reservation_id}`);
            fetchAvailableVehicles(); // Refresh the list of available vehicles
        })
        .catch(err => {
            console.error("Error creating reservation:", err);
        });
});

// Fetch available vehicles on page load
document.addEventListener("DOMContentLoaded", function () {
    console.log("Page loaded, fetching available vehicles...");
    fetchAvailableVehicles();
});