let searchInput;

document.addEventListener("DOMContentLoaded", () => {
  searchInput = document.getElementById("search-input");

  const searchButton = document.getElementById("search-button");

  searchInput.focus();

  searchInput.addEventListener("keypress", (event) => {
    if (event.key === "Enter") {
      makeSearchRequest();
    }
  });

  searchButton.addEventListener("click", makeSearchRequest);
});

function makeSearchRequest() {
  const query = searchInput.value;
  const url = new URL(window.location.href);

  url.searchParams.set("q", encodeURIComponent(query));

  window.location.href = url.toString();
}
