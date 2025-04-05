document.addEventListener("DOMContentLoaded", function () {
  const form = document.getElementById("login-form");

  form.addEventListener("submit", async function (e) {
    e.preventDefault();

    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;

    console.log(username, password);

    try {
      const response = await fetch("/api/users/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, password }),
        credentials: "same-origin",
      });

      if (response.ok) {
        window.location.href = "/";
      } else {
        const text = await response.text();
        alert(text || "Login failed");
      }
    } catch (error) {
      console.error("Error during login:", error);
      alert(error);
    }
  });
});
