const vehicleForm = document.getElementById("vehicleForm");
const vehicleTable = document.getElementById("vehicleTable").getElementsByTagName("tbody")[0];

// Corrected API URL variable name
const apiUrl = "http://localhost:5000/vehicles"; // Go server API URL

// Event listener for form submission
vehicleForm.addEventListener("submit", function (event) {
    event.preventDefault();
    const vehicleId = document.getElementById("vehicleId").value;
    const make = document.getElementById("make").value;
    const model = document.getElementById("model").value;
    const availability = document.getElementById("availability").checked;

    if (vehicleId) {
        // Update vehicle
        updateVehicle(vehicleId, make, model, availability);
    } else {
        // Create new vehicle
        createVehicle(make, model, availability);
    }
});

// Fetch all vehicles and display them in the table
function fetchVehicles() {
    fetch(apiUrl)
        .then(response => response.json())
        .then(vehicles => {
            vehicleTable.innerHTML = ""; // Clear the table body
            vehicles.forEach(vehicle => {
                const row = vehicleTable.insertRow();
                row.innerHTML = `
                    <td>${vehicle.id}</td>
                    <td>${vehicle.make}</td>
                    <td>${vehicle.model}</td>
                    <td>${vehicle.availability ? "Available" : "Not Available"}</td>
                    <td>
                        <button class="edit" onclick="editVehicle(${vehicle.id})">Edit</button>
                        <button class="delete" onclick="deleteVehicle(${vehicle.id})">Delete</button>
                    </td>
                `;
            });
        })
        .catch(err => console.error("Error fetching vehicles:", err));
}

// Create a new vehicle
function createVehicle(make, model, availability) {
    const vehicle = { make, model, availability };
    fetch(apiUrl, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(vehicle)
    })
    .then(response => response.json())
    .then(() => {
        fetchVehicles(); // Refresh the vehicle list
        vehicleForm.reset(); // Reset the form
    })
    .catch(err => console.error("Error creating vehicle:", err));
}

// Edit a vehicle (populate form for updating)
function editVehicle(id) {
    fetch(`${apiUrl}/${id}`)
        .then(response => response.json())
        .then(vehicle => {
            document.getElementById("vehicleId").value = vehicle.id;
            document.getElementById("make").value = vehicle.make;
            document.getElementById("model").value = vehicle.model;
            document.getElementById("availability").checked = vehicle.availability;
        })
        .catch(err => console.error("Error fetching vehicle for editing:", err));
}

// Update a vehicle
function updateVehicle(id, make, model, availability) {
    const vehicle = { make, model, availability };
    fetch(`${apiUrl}/${id}`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(vehicle)
    })
    .then(() => {
        fetchVehicles(); // Refresh the vehicle list
        vehicleForm.reset(); // Reset the form
    })
    .catch(err => console.error("Error updating vehicle:", err));
}

// Delete a vehicle
function deleteVehicle(id) {
    fetch(`${apiUrl}/${id}`, {
        method: "DELETE"
    })
    .then(() => fetchVehicles()) // Refresh the vehicle list
    .catch(err => console.error("Error deleting vehicle:", err));
}

// Initial call to fetch all vehicles when the page loads
document.addEventListener("DOMContentLoaded", fetchVehicles);