document.addEventListener("DOMContentLoaded", () => {
    const userId = sessionStorage.getItem("userId");

    console.log("Retrieved user ID from sessionStorage:", userId); // Debug log

    if (!userId || isNaN(userId)) {
        console.error("Invalid or missing user ID. Redirecting to login.");
        alert("No user logged in. Redirecting to login page.");
        window.location.href = "login.html";
        return;
    }

    fetch(`http://localhost:5001/api/v1/users/${userId}`)
        .then((response) => {
            if (!response.ok) {
                throw new Error("User not found");
            }
            return response.json();
        })
        .then((data) => {
            console.log("User profile data fetched:", data); // Debug log
            document.getElementById("userName").textContent = data.name;
            document.getElementById("userEmail").textContent = data.email;
            document.getElementById("membershipTier").textContent = data.membership_tier;
        })
        .catch((error) => {
            console.error("Error fetching user profile:", error);
            alert("Failed to load user profile.");
        });
});

document.getElementById("updateProfileForm").addEventListener("submit", (e) => {
    e.preventDefault();

    const userId = sessionStorage.getItem("userId");
    const updatedDetails = {
        name: document.getElementById("updateName").value,
        email: document.getElementById("updateEmail").value,
        password: document.getElementById("updatePassword").value, // Include password if necessary
        membership_tier: document.getElementById("updateTier").value,
    };

    fetch(`http://localhost:5001/api/v1/users/${userId}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(updatedDetails),
    })
        .then((response) => {
            if (!response.ok) {
                throw new Error("Failed to update profile");
            }
            return response.json();
        })
        .then((data) => {
            alert("Profile updated successfully!");
            console.log("Updated profile:", data);
        })
        .catch((error) => console.error("Error updating profile:", error));
});