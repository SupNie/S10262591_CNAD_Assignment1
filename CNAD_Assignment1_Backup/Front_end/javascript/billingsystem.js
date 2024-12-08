const baseURL = "http://localhost:5002"; // Replace with your backend URL if different

// Fetch billing info when the user enters a billing ID
document.getElementById("fetchBillingBtn").addEventListener("click", () => {
    const billingId = document.getElementById("billingId").value;

    if (billingId === "") {
        alert("Please enter a billing ID.");
        return;
    }

    fetch(`${baseURL}/billings/${billingId}`)
        .then(response => response.json())
        .then(data => {
            document.getElementById("billingIdDetails").textContent = data.id;
            document.getElementById("amountDetails").textContent = data.amount.toFixed(2);
            document.getElementById("paymentStatusDetails").textContent = data.payment_status;
        })
        .catch(error => {
            console.error("Error fetching billing details:", error);
            alert("Failed to fetch billing details.");
        });
});

// Generate invoice when the user clicks the "Generate Invoice" button
// Generate invoice with cost details
document.getElementById("generateInvoiceBtn").addEventListener("click", () => {
    const billingId = document.getElementById("billingId").value;

    if (billingId === "") {
        alert("Please enter a billing ID.");
        return;
    }

    fetch(`${baseURL}/invoices/${billingId}`)
        .then(response => response.json())
        .then(data => {
            document.getElementById("vehicleTypeDetails").textContent = data.vehicle_type;
            document.getElementById("membershipLevelDetails").textContent = data.membership_level;
            document.getElementById("totalCostDetails").textContent = data.amount.toFixed(2);
            document.getElementById("invoiceMessage").textContent = `Invoice generated. Total Cost: $${data.amount.toFixed(2)}`;
        })
        .catch(error => {
            console.error("Error generating invoice:", error);
            document.getElementById("invoiceMessage").textContent = "Failed to generate invoice.";
        });
});
// Generate receipt when the user clicks the "Generate Receipt" button
document.getElementById("generateReceiptBtn").addEventListener("click", () => {
    const billingId = document.getElementById("billingId").value;

    if (billingId === "") {
        alert("Please enter a billing ID.");
        return;
    }

    fetch(`${baseURL}/receipts/${billingId}`)
        .then(response => response.json())
        .then(data => {
            document.getElementById("receiptMessage").textContent = `Receipt generated. Amount: $${data.amount.toFixed(2)}. Date: ${new Date(data.payment_date).toLocaleDateString()}`;
        })
        .catch(error => {
            console.error("Error generating receipt:", error);
            document.getElementById("receiptMessage").textContent = "Failed to generate receipt.";
        });
});