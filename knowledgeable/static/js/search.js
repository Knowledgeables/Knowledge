let searchInput;

document.addEventListener('DOMContentLoaded', () => {
  searchInput = document.getElementById("search-input");

  searchInput.focus();

  searchInput.addEventListener('keypress', (event) => {
    if (event.key === 'Enter') {
      makeSearchRequest();
    }
  });
});

function makeSearchRequest() {
  const query = searchInput.value;
  const url = new URL(window.location.href);
  url.searchParams.set('q', query);
  window.location.href = url.toString();
}

const searchButton = document.getElementById("search-button");

searchButton.addEventListener('click', () => makeSearchRequest())