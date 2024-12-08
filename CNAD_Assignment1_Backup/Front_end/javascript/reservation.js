const baseURL = "http://localhost:5000/reservations";

// Fetch Available Vehicles
document.getElementById("availabilityForm").addEventListener("submit", (e) => {
    e.preventDefault();

    const startTime = document.getElementById("startTime").value;
    const endTime = document.getElementById("endTime").value;

    fetch(`${baseURL}/vehicles/available?start_time=${startTime}&end_time=${endTime}`)
        .then((response) => response.json())
        .then((data) => {
            const vehicleList = document.getElementById("vehicleList");
            vehicleList.innerHTML = "";

            if (data.length === 0) {
                vehicleList.innerHTML = "<li>No vehicles available for the selected time range.</li>";
            } else {
                data.forEach((vehicle) => {
                    const listItem = document.createElement("li");
                    listItem.textContent = `ID: ${vehicle.id}, Make: ${vehicle.make}, Model: ${vehicle.model}`;
                    vehicleList.appendChild(listItem);
                });
            }
        })
        .catch((error) => console.error("Error fetching available vehicles:", error));
});

// Create a Reservation
document.getElementById("reservationForm").addEventListener("submit", (e) => {
    e.preventDefault();

    const vehicleId = document.getElementById("vehicleId").value;
    const userId = document.getElementById("userId").value;
    const startTime = document.getElementById("reservationStartTime").value;
    const endTime = document.getElementById("reservationEndTime").value;

    fetch(`${baseURL}/reservations`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ vehicle_id: vehicleId, user_id: userId, start_time: startTime, end_time: endTime }),
    })
        .then((response) => response.json())
        .then((data) => {
            document.getElementById("reservationResponse").textContent = "Reservation created successfully!";
        })
        .catch((error) => {
            console.error("Error creating reservation:", error);
            document.getElementById("reservationResponse").textContent = "Failed to create reservation.";
        });
});

// Modify a Reservation
document.getElementById("modifyReservationForm").addEventListener("submit", (e) => {
    e.preventDefault();

    const reservationId = document.getElementById("reservationId").value;
    const newStartTime = document.getElementById("newStartTime").value;
    const newEndTime = document.getElementById("newEndTime").value;

    fetch(`${baseURL}/reservations/${reservationId}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ start_time: newStartTime, end_time: newEndTime }),
    })
        .then((response) => response.json())
        .then((data) => {
            document.getElementById("modifyResponse").textContent = "Reservation modified successfully!";
        })
        .catch((error) => {
            console.error("Error modifying reservation:", error);
            document.getElementById("modifyResponse").textContent = "Failed to modify reservation.";
        });
});

// Cancel a Reservation
document.getElementById("cancelReservationForm").addEventListener("submit", (e) => {
    e.preventDefault();

    const reservationId = document.getElementById("cancelReservationId").value;

    fetch(`${baseURL}/reservations/${reservationId}`, {
        method: "DELETE",
    })
        .then((response) => {
            if (response.ok) {
                document.getElementById("cancelResponse").textContent = "Reservation canceled successfully!";
            } else {
                document.getElementById("cancelResponse").textContent = "Failed to cancel reservation.";
            }
        })
        .catch((error) => {
            console.error("Error canceling reservation:", error);
            document.getElementById("cancelResponse").textContent = "Failed to cancel reservation.";
        });
});
