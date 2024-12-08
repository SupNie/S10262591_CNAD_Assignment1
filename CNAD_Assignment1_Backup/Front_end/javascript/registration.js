const baseURL = "http://localhost:5001/api/v1/users";

// Handle Registration Form Submission
document.getElementById("registerForm").addEventListener("submit", (e) => {
  e.preventDefault();

  const id = document.getElementById("id").value;
  const name = document.getElementById("name").value;
  const email = document.getElementById("email").value;
  const password = document.getElementById("password").value;
  const membershipTier = document.getElementById("membershipTier").value;

  fetch(baseURL, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ id, name, email, password, membership_tier: membershipTier }),
  })
    .then((response) => response.json())
    .then((data) => {
      document.getElementById("registerResponse").textContent = "User registered successfully!";
      document.getElementById("registerForm").reset();
    })
    .catch((error) => {
      console.error("Error:", error);
      document.getElementById("registerResponse").textContent = "Registration failed.";
    });
});

// Handle Get All Users
document.getElementById("getAllUsers").addEventListener("click", () => {
  fetch(baseURL)
    .then((response) => response.json())
    .then((data) => {
      const usersList = document.getElementById("usersList");
      usersList.innerHTML = "";
      data.forEach((user) => {
        const div = document.createElement("div");
        div.textContent = `ID: ${user.id}, Name: ${user.name}, Email: ${user.email}, Membership Tier: ${user.membership_tier}`;
        usersList.appendChild(div);
      });
    })
    .catch((error) => console.error("Error:", error));
});

// Handle Update Form Submission
document.getElementById("updateForm").addEventListener("submit", (e) => {
    e.preventDefault();
  
    const userId = document.getElementById("userId").value; // The ID of the user to update
    const name = document.getElementById("updateName").value;
    const email = document.getElementById("updateEmail").value;
    const password = document.getElementById("updatePassword").value;
    const membershipTier = document.getElementById("updateMembershipTier").value;
  
    fetch(`${baseURL}/${userId}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        name,
        email,
        password,
        membership_tier: membershipTier,
      }),
    })
      .then((response) => response.json())
      .then((data) => {
        document.getElementById("updateResponse").textContent = "User updated successfully!";
        document.getElementById("updateForm").reset();
      })
      .catch((error) => {
        console.error("Error:", error);
        document.getElementById("updateResponse").textContent = "Update failed.";
      });
  });
  
// Handle Delete User
document.getElementById("deleteForm").addEventListener("submit", (e) => {
  e.preventDefault();

  const userId = document.getElementById("deleteUserId").value;

  fetch(`${baseURL}/${userId}`, {
    method: "DELETE",
  })
    .then((response) => {
      if (response.ok) {
        document.getElementById("deleteResponse").textContent = "User deleted successfully!";
      } else {
        document.getElementById("deleteResponse").textContent = "Failed to delete user.";
      }
    })
    .catch((error) => {
      console.error("Error:", error);
      document.getElementById("deleteResponse").textContent = "Delete operation failed.";
    });
});


