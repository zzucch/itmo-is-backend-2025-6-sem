document.addEventListener("DOMContentLoaded", function () {
  const form = document.getElementById("signup-form");

  form.addEventListener("submit", async function (e) {
    e.preventDefault();

    const username = document.getElementById("username").value;
    const email = document.getElementById("email").value;
    const password = document.getElementById("password").value;
    const confirmPassword = document.getElementById("confirm-password").value;

    if (password !== confirmPassword) {
      alert("Passwords do not match");
      return;
    }

    try {
      const response = await fetch("/api/users", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, email, password }),
        credentials: "same-origin",
      });

      if (response.ok) {
        alert("Signed up. Please log in.");
        window.location.href = "/login";
      } else {
        const errorData = await response.json();
        alert(errorData.message || "Signing up failed");
      }
    } catch (error) {
      alert(error);
    }
  });
});
