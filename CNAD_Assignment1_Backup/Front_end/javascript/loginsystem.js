document.getElementById("loginForm").addEventListener("submit", (e) => {
    e.preventDefault();

    const email = document.getElementById("email").value;
    const password = document.getElementById("password").value;

    console.log("Attempting login with email:", email); // Debug log

    fetch("http://localhost:5001/api/v1/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
    })
        .then((response) => {
            if (!response.ok) {
                throw new Error("Login failed");
            }
            return response.json();
        })
        .then((data) => {
            console.log("Login successful, user ID:", data.id); // Debug log
            if (!data.id || typeof data.id !== "number") {
                throw new Error("Invalid user ID");
            }
            sessionStorage.setItem("userId", data.id);
            window.location.href = "profile.html";
        })
        .catch((error) => {
            console.error("Error logging in:", error);
            document.getElementById("errorMessage").textContent = "Invalid email or password.";
            document.getElementById("errorMessage").style.display = "block";
        });
});
