import { fetchWithCache } from "./cached_fetch.js";

function displayCart(cartItems) {
  const cartTable = document.getElementById("cart-items");
  cartTable.innerHTML = "";

  if (!cartItems || cartItems.length === 0) {
    cartTable.innerHTML = '<tr><td colspan="4">Your cart is empty</td></tr>';
    document.getElementById("cart-total").textContent = "$0.00";
    return;
  }

  let total = 0;

  cartItems.forEach((item) => {
    const itemRow = document.createElement("tr");
    itemRow.className = "cart-item";
    itemRow.dataset.itemId = item.ID;

    const propertiesToShow = {
      Name: item.Name || "Unknown Item",
      Price: `$${(item.Price || 0).toFixed(2)}`,
      Image: item.URL ? `<img src="${item.URL}" width="100">` : "No image",
    };

    for (const key in propertiesToShow) {
      const cell = document.createElement("td");
      cell.innerHTML = propertiesToShow[key];
      itemRow.appendChild(cell);
    }

    const actionCell = document.createElement("td");
    const removeBtn = document.createElement("button");
    removeBtn.textContent = "Remove";
    removeBtn.onclick = () => removeFromCart(item.ID);
    actionCell.appendChild(removeBtn);
    itemRow.appendChild(actionCell);

    cartTable.appendChild(itemRow);
    total += item.Price || 0;
  });

  document.getElementById("cart-total").textContent = `$${total.toFixed(2)}`;
}

function loadCart() {
  fetchWithCache("/api/users/me/cart")
    .then((data) => {
      displayCart(data);
    })
    .catch((error) => {
      alert("Error loading cart: " + error.message);
      document.getElementById("cart-items").innerHTML = `
                    <tr>
                        <td colspan="4" class="error-message">
                            Error loading cart: ${error.message}
                            <button onclick="loadCart()">Retry</button>
                        </td>
                    </tr>
                `;
    });
}

async function removeFromCart(itemId) {
  if (!confirm("Are you sure you want to remove this item from your cart?")) {
    return;
  }

  try {
    const response = await fetch(`/api/users/me/cart/${itemId}`, {
      method: "DELETE",
    });

    const text = await response.text();

    if (!response.ok) {
      throw new Error(`Error removing item: ${text}`);
    }

    alert(text);
    loadCart();
  } catch (error) {
    alert("Failed to remove item: " + error.message);
  }
}

function checkout() {
  fetchWithCache("/api/users/me/cart")
    .then((cartItems) => {
      if (!cartItems || cartItems.length === 0) {
        alert("Your cart is empty!");
        return;
      }

      const orderSummary = document.getElementById("order-summary");
      let total = 0;
      orderSummary.innerHTML = "<ul>";

      cartItems.forEach((item) => {
        orderSummary.innerHTML += `<li>${item.Name || "Unknown Item"} - $${(item.Price || 0).toFixed(2)}</li>`;
        total += item.Price || 0;
      });

      orderSummary.innerHTML += `</ul><p><strong>Total: $${total.toFixed(2)}</strong></p>`;

      document.getElementById("order-modal").style.display = "block";
    })
    .catch((error) => {
      alert("Error during checkout: " + error.message);
    });
}

async function confirmOrder() {
  const shippingAddress = document.getElementById("shipping-address").value;
  if (!shippingAddress.trim()) {
    alert("Please enter a shipping address");
    return;
  }

  try {
    const cartResponse = await fetchWithCache("/api/users/me/cart");
    const cartItems = cartResponse;

    if (!cartItems || cartItems.length === 0) {
      alert("Your cart is empty! Please add items before placing an order.");
      return;
    }

    const formData = new URLSearchParams();
    formData.append("shipping_address", shippingAddress);
    cartItems.forEach((item) => {
      formData.append("phone_ids", item.ID.toString());
    });

    const orderResponse = await fetch("/api/orders", {
      method: "POST",
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
      },
      body: formData.toString(),
    });

    if (!orderResponse.ok) {
      const errorText = await orderResponse.text();
      throw new Error(errorText);
    }

    const result = await orderResponse.json();
    alert(`Order #${result.order_id} placed successfully!`);
    document.getElementById("order-modal").style.display = "none";
    loadCart();
  } catch (error) {
    alert(error.message);
  }
}

document.addEventListener("DOMContentLoaded", function () {
  loadCart();
  loadOrders();
});

function displayOrders(orders) {
  const orderHistoryDiv = document.getElementById("order-history");

  if (!orders || orders.length === 0) {
    orderHistoryDiv.innerHTML = "<p>You haven't placed any orders yet.</p>";
    return;
  }

  let html = `
        <table border="1" width="100%" style="margin-top: 20px;">
            <thead>
                <tr>
                    <th>Order ID</th>
                    <th>Date</th>
                    <th>Items</th>
                    <th>Total</th>
                    <th>Status</th>
                </tr>
            </thead>
            <tbody>`;

  orders.forEach((order) => {
    const orderDate = new Date(order.CreatedAt).toLocaleDateString();
    let total = 0;
    let itemsHtml = '<ul style="margin: 0; padding-left: 20px;">';

    order.Phones.forEach((phone) => {
      itemsHtml += `<li>${phone.Name} - $${phone.Price.toFixed(2)}</li>`;
      total += phone.Price;
    });

    itemsHtml += "</ul>";

    html += `
                <tr>
                    <td>#${order.ID}</td>
                    <td>${orderDate}</td>
                    <td>${itemsHtml}</td>
                    <td>$${total.toFixed(2)}</td>
                    <td>${order.Status || "Processing"}</td>
                </tr>`;
  });

  html += `</tbody></table>`;
  orderHistoryDiv.innerHTML = html;
}

function loadOrders() {
  fetchWithCache("/api/orders/me")
    .then((data) => {
      displayOrders(data);
    })
    .catch((error) => {
      document.getElementById("order-history").innerHTML = `
                    <p class="error-message">
                        ${error.message}
                        <button onclick="loadOrders()">Retry</button>
                    </p>`;
    });
}
