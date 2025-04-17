function showForm(id) {
  const elem = document.getElementById(id);
  elem.style.display === "none"
    ? (elem.style.display = "block")
    : (elem.style.display = "none");
}
