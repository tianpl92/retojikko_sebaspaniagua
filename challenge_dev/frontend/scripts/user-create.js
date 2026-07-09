// Public Calls Portal — User Create Script

async function handleUserCreate(event) {
  event.preventDefault();

  const userData = {
    id:           document.getElementById("documento").value.trim(),
    email:        document.getElementById("email").value.trim(),
    password:     document.getElementById("password").value,
    first_name:   document.getElementById("nombres").value.trim(),
    last_name:    document.getElementById("apellidos").value.trim(),
    gender:       document.getElementById("genero").value,
    phone_number: document.getElementById("telefono").value.trim(),
    status:       "AC",
  };

  // Basic validation
  const required = ["id", "email", "password", "first_name", "last_name"];
  const missing  = required.filter(f => !userData[f]);
  if (missing.length > 0) {
    alert("Por favor completa todos los campos obligatorios.");
    return;
  }

  try {
    const response = await fetch(BACKEND_URL + ENDPOINTS.USER_CREATE, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(userData),
    });

    const result = await response.json();

    if (response.ok) {
      alert("Usuario creado correctamente");
      window.location.href = "login.html";
    } else {
      alert(result.error || "Error al crear el usuario.");
    }
  } catch (err) {
    alert("Error de conexión con el servidor.");
  }
}
