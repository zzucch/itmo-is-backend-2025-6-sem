document.addEventListener("DOMContentLoaded", function () {
  checkAuthStatus();

  document.addEventListener("click", function (e) {
    if (e.target && e.target.id === "logout-link") {
      e.preventDefault();
      handleLogout();
    }
  });
});

async function checkAuthStatus() {
  const authLinks = document.getElementById("auth-links");
  if (!authLinks) return;

  try {
    const response = await fetch("/api/users/me", {
      method: "GET",
      credentials: "same-origin",
    });

    if (response.ok) {
      authLinks.innerHTML = `
                <a href="/cart">Cart</a>
                <a href="/sell">Sell</a>
                <a href="#" id="logout-link">Sign out</a>
            `;
    } else {
      authLinks.innerHTML = `
                <a href="/login">Sign in</a>
                <a href="/signup">Sign up</a>
            `;
    }
  } catch (error) {
    console.error("Error checking auth status:", error);
    authLinks.innerHTML = `
            <a href="/login">Sign in</a>
            <a href="/signup">Sign up</a>
        `;
  }
}

async function handleLogout() {
  try {
    const response = await fetch("/api/users/logout", {
      method: "POST",
      credentials: "same-origin",
    });

    if (response.ok) {
      window.location.href = "/";
    } else {
      console.error("Logout failed");
    }
  } catch (error) {
    console.error("Error during logout:", error);
  }
}
