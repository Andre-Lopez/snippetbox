var navLinks = document.querySelectorAll("nav a")
for (const navLink of navLinks) {
  if (navLink.getAttribute("href") == window.location.pathname) {
    navLink.classList.add("live")
    break
  }
}
