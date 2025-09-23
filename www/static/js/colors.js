const bulb = document.getElementById("bulb");
const moon = document.getElementById("moon");

const bulbMobile = document.getElementById("bulb-mobile");
const moonMobile = document.getElementById("moon-mobile");

if (localStorage.getItem("pref-theme") === "dark") {
  document.body.classList.add("dark");
  bulb.style.display = "inline";
  bulbMobile.style.display = "inline";
} else if (localStorage.getItem("pref-theme") === "light") {
  document.body.classList.add("light");
  moon.style.display = "inline";
  moonMobile.style.display = "inline";
} else if (
  window.matchMedia &&
  window.matchMedia("(prefers-color-scheme: dark)").matches
) {
  document.body.classList.add("dark");
  bulb.style.display = "inline";
  bulbMobile.style.display = "inline";
} else {
  document.body.classList.add("light");
  moonMobile.style.display = "inline";
}

function toggleTheme() {
  if (document.body.className.includes("dark")) {
    moonMobile.style.display = "inline";
    moon.style.display = "inline";
    bulb.style.display = "none";
    bulbMobile.style.display = "none";
    document.body.classList.remove("dark");
    document.body.classList.add("light");
    localStorage.setItem("pref-theme", "light");
  } else {
    bulb.style.display = "inline";
    bulbMobile.style.display = "inline";
    moon.style.display = "none";
    moonMobile.style.display = "none";
    document.body.classList.add("dark");
    document.body.classList.remove("light");
    localStorage.setItem("pref-theme", "dark");
  }
}

document.getElementById("theme-toggle-mobile").addEventListener("click", () => {
  toggleTheme();
});


document.getElementById("theme-toggle").addEventListener("click", () => {
  toggleTheme();
});
