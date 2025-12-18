try {
  const raw = localStorage.getItem("homebox/preferences/location");
  const prefs = raw ? JSON.parse(raw) : null;
  const theme = prefs?.theme || "hatchbridge";
  document.documentElement.setAttribute("data-theme", theme);
  document.documentElement.classList.add("theme-" + theme);
} catch (e) {
  console.error("Failed to set theme", e);
}
