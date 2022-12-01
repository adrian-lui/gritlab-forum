window.onload = checkSession();
async function checkSession() {
  await fetch("/checkSession")
    .then((response) => response.json())
    .then((json) => {
      if (json.status) {
        window.location.replace("/loginSuccess");
        return;
      }
    });
}
