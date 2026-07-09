// Public Calls Portal — Login Script

async function handleLogin(event) {
  event.preventDefault();

  const email    = document.getElementById("email").value.trim();
  const password = document.getElementById("password").value;

  if (!email || !password) {
    alert("Por favor completa todos los campos.");
    return;
  }

  try {
    const response = await fetch(BACKEND_URL + ENDPOINTS.LOGIN, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });

    if (response.ok) {
      const data = await response.json();
      // Store token for future use
      localStorage.setItem("jwt_token", data.token);
      localStorage.setItem("user_id", data.user.id);
      localStorage.setItem("user_name", data.user.first_name);
      localStorage.setItem("user_last_name", data.user.last_name);
      alert("Ingresaste correctamente.");
      window.location.href = "dashboard.html";
    } else {
      alert("Datos inválidos.");
    }
  } catch (err) {
    alert("Error de conexión con el servidor.");
  }
}
