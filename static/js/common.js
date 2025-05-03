import { fetchWithCache } from "./cached_fetch.js";

// i know its kind of pointless but whatever
document.addEventListener("DOMContentLoaded", () => {
  const gotoTopButton = document.createElement("button");
  gotoTopButton.innerText = "Top";
  gotoTopButton.style.position = "fixed";
  gotoTopButton.style.bottom = "20px";
  gotoTopButton.style.right = "20px";
  gotoTopButton.style.display = "none";
  gotoTopButton.style.zIndex = 1100;
  gotoTopButton.className = "goto-top-button";

  document.body.appendChild(gotoTopButton);

  window.onscroll = () => {
    gotoTopButton.style.display = window.scrollY > 200 ? "block" : "none";
  };

  gotoTopButton.onclick = () => {
    window.scrollTo({ top: 0, behavior: "smooth" });
  };
});

document.addEventListener("DOMContentLoaded", async function () {
  try {
    const user = await fetchWithCache("/api/users/me", {
      credentials: "same-origin",
    });

    console.log("Authenticated as:", user.username);
  } catch (error) {
    console.error("Error checking auth status:", error);
  }
});

document.querySelectorAll(".get-now-button, .sale-button").forEach((button) => {
  button.addEventListener("click", async function () {
    try {
      const response = await fetch("/api/users/me", {
        credentials: "same-origin",
      });

      if (!response.ok) {
        window.location.href = "/login";
        return;
      }
    } catch (error) {
      console.error("Error:", error);
    }
  });
});
