const startTime = performance.now();

(function () {
  window.onload = function () {
    const endTime = performance.now();
    const loadTimeMs = endTime - startTime;

    const loadTimeInfo = document.createElement("div");
    loadTimeInfo.innerText = `Page loaded in ${loadTimeMs.toFixed(2)} ms`;
    loadTimeInfo.style.backgroundColor = "#222";
    loadTimeInfo.style.color = "#aaa";
    loadTimeInfo.style.padding = "10px";

    document.body.appendChild(loadTimeInfo);
  };

  fetch(window.location.href, { method: "HEAD" })
    .then((res) => {
      const serverElapsed = res.headers.get("X-Elapsed-Time");

      if (serverElapsed) {
        const serverTimeInfo = document.createElement("div");
        serverTimeInfo.innerText = `backend elapsed time: ${serverElapsed}`;
        serverTimeInfo.style.backgroundColor = "#222";
        serverTimeInfo.style.color = "#aaa";
        serverTimeInfo.style.padding = "10px";
        document.body.appendChild(serverTimeInfo);
      }
    })
    .catch((err) => {
      console.error("Failed to fetch server header:", err);
    });
})();
