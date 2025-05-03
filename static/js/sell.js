import { fetchWithCache } from "./cached_fetch.js";

const specifications = [
  "Name",
  "Description",
  "Brand",
  "CPU",
  "ScreenSize",
  "Camera",
  "Battery",
  "Storage",
  "Price",
  "IsUsed",
  "Issues",
  "Condition",
  "Image",
];

let currentSpecificationIndex = 0;
let phoneDetails = {};

function addSpecifications(event) {
  event.preventDefault();

  const specification = specifications[currentSpecificationIndex];
  let value;

  if (specification === "Image") {
    const fileInput = document.getElementById("image-upload");
    const file = fileInput.files[0];

    if (!file) {
      alert("Please select an image file.");
      return;
    }

    if (!file.type.match("image.*")) {
      alert("Please select a valid image file.");
      return;
    }

    phoneDetails.ImageFile = file;
    currentSpecificationIndex++;
  } else {
    value = document.getElementById("value").value.trim();

    if (!value) {
      alert("Please enter a value for the specification.");
      return;
    }

    if (specification === "IsUsed") {
      value = value.toLowerCase() === "true";
    } else if (specification === "Price") {
      value = parseFloat(value);
      if (isNaN(value)) {
        alert("Please enter a valid price.");
        return;
      }
    }

    phoneDetails[specification] = value;
    currentSpecificationIndex++;
  }

  updateTable();

  if (currentSpecificationIndex < specifications.length) {
    document.getElementById("specification-display").textContent =
      specifications[currentSpecificationIndex];

    if (specifications[currentSpecificationIndex] === "Image") {
      document.getElementById("value-input").style.display = "none";
      document.getElementById("image-upload-container").style.display = "block";
    } else {
      document.getElementById("value-input").style.display = "block";
      document.getElementById("image-upload-container").style.display = "none";
    }
  } else {
    document.getElementById("add-specification-form").style.display = "none";
    document.getElementById("submit-button").style.display = "block";
  }

  document.getElementById("value").value = "";
}

function updateTable() {
  const tableContainer = document.getElementById("table-container");
  tableContainer.innerHTML = "";

  Object.entries(phoneDetails).forEach(([key, value]) => {
    const row = document.createElement("div");
    row.className = "row";

    const keyDiv = document.createElement("div");
    keyDiv.className = "key";
    keyDiv.textContent = key;

    const valueDiv = document.createElement("div");
    valueDiv.className = "value";
    valueDiv.textContent = value;

    row.appendChild(keyDiv);
    row.appendChild(valueDiv);
    tableContainer.appendChild(row);
  });
}

async function confirmListing() {
  try {
    const formData = new FormData();

    for (const key in phoneDetails) {
      if (key !== "ImageFile") {
        formData.append(key, phoneDetails[key]);
      }
    }

    if (phoneDetails.ImageFile) {
      formData.append("image", phoneDetails.ImageFile);
    }

    const response = await fetch("/api/phones", {
      method: "POST",
      body: formData,
    });

    if (response.ok) {
      document.getElementById("success-message").style.display = "block";
      setTimeout(() => {
        document.getElementById("success-message").style.display = "none";
      }, 3000);

      resetForm();
      fetchPhones();
    } else {
      const text = await response.text();
      alert(text);
    }
  } catch (error) {
    alert(error);
  }
}

function resetForm() {
  phoneDetails = {};
  currentSpecificationIndex = 0;
  document.getElementById("specification-display").textContent =
    specifications[currentSpecificationIndex];
  document.getElementById("add-specification-form").style.display = "block";
  document.getElementById("submit-button").style.display = "none";
  document.getElementById("table-container").innerHTML = "";
  document.getElementById("value-input").style.display = "block";
  document.getElementById("image-upload-container").style.display = "none";
  document.getElementById("image-upload").value = "";
}

async function fetchPhones() {
  try {
    const phones = await fetchWithCache("/api/phones");
    displayPhones(phones);
  } catch (error) {
    console.error("Fetch error:", error);
    alert(error);
  }
}

function displayPhones(phones) {
  const finalTableContainer = document.getElementById("final-table-container");
  finalTableContainer.innerHTML = "";

  if (phones.length === 0) {
    finalTableContainer.style.display = "none";
    return;
  }

  finalTableContainer.style.display = "block";

  phones.forEach((phone) => {
    const phoneEntry = document.createElement("div");
    phoneEntry.className = "final-phone";

    for (const key in phone) {
      const value = phone[key];

      const detailRow = document.createElement("div");
      detailRow.className = "row";

      const keyDiv = document.createElement("div");
      keyDiv.className = "key";
      keyDiv.textContent = key;

      const valueDiv = document.createElement("div");
      valueDiv.className = "value";
      valueDiv.textContent = value;

      detailRow.appendChild(keyDiv);
      detailRow.appendChild(valueDiv);
      phoneEntry.appendChild(detailRow);
    }

    const deleteButton = document.createElement("button");
    deleteButton.className = "delete-button";
    deleteButton.textContent = "Delete";
    deleteButton.onclick = () => deletePhone(phone.ID);

    phoneEntry.appendChild(deleteButton);
    finalTableContainer.appendChild(phoneEntry);
  });
}

function deletePhone(phoneId) {
  console.log("deleting!", phoneId);

  fetch(`/api/phones/${phoneId}`, {
    method: "DELETE",
  })
    .then((response) => {
      if (response.ok) {
        fetchPhones();
      } else {
        alert("Failed to delete phone. Please try again.");
      }
    })
    .catch((error) => {
      console.error("Error deleting phone:", error);
      alert("An error occurred while deleting the phone.");
    });
}

document.addEventListener("DOMContentLoaded", function () {
  document.getElementById("specification-display").textContent =
    specifications[currentSpecificationIndex];

  document
    .getElementById("add-specification-form")
    .addEventListener("submit", addSpecifications);
  document
    .getElementById("submit-button")
    .addEventListener("click", confirmListing);

  fetchPhones();
});
