// Public Calls Portal — Saved Proposals Logic

let savedCache = [];

// =============================================
// Init
// =============================================
document.addEventListener("DOMContentLoaded", async function () {
  const token = localStorage.getItem("jwt_token");
  if (!token) {
    window.location.href = "login.html";
    return;
  }

  // Show user name
  const name = localStorage.getItem("user_name") || "Usuario";
  const lastName = localStorage.getItem("user_last_name") || "";
  document.getElementById("userDisplay").textContent = name + " " + lastName;

  await loadSavedProposals();
});

// =============================================
// Go back to dashboard
// =============================================
function goToDashboard() {
  window.location.href = "dashboard.html";
}

// =============================================
// Fetch saved proposals
// =============================================
async function loadSavedProposals() {
  const token = localStorage.getItem("jwt_token");
  if (!token) {
    window.location.href = "login.html";
    return;
  }

  showLoading();

  try {
    const response = await fetch(BACKEND_URL + ENDPOINTS.SAVED_LIST, {
      headers: { Authorization: "Bearer " + token },
    });

    if (response.status === 401) {
      localStorage.clear();
      window.location.href = "login.html";
      return;
    }

    if (!response.ok) {
      hideLoading();
      document.getElementById("savedBody").innerHTML =
        '<tr><td colspan="4" style="text-align:center;padding:40px;color:#a94442;">Error al cargar convocatorias guardadas</td></tr>';
      return;
    }

    const data = await response.json();
    savedCache = data;

    hideLoading();
    renderGrid(data);
  } catch (err) {
    hideLoading();
    document.getElementById("savedBody").innerHTML =
      '<tr><td colspan="4" style="text-align:center;padding:40px;color:#a94442;">Error de conexión con el servidor</td></tr>';
  }
}

// =============================================
// Render grid
// =============================================
function renderGrid(items) {
  const tbody = document.getElementById("savedBody");

  if (!items || items.length === 0) {
    tbody.innerHTML =
      '<tr><td colspan="4" style="text-align:center;padding:40px;color:#7a7a8e;">No hay convocatorias guardadas</td></tr>';
    return;
  }

  let html = "";
  for (let i = 0; i < items.length; i++) {
    const p = items[i];
    const procId = p.public_call_id != null ? p.public_call_id : "—";
    const assocDate = formatDate(p.association_date);
    const createdDate = formatDate(p.created_at);

    html += "<tr>";
    html += "<td>" + escHtml(String(procId)) + "</td>";
    html += "<td>" + escHtml(assocDate) + "</td>";
    html += "<td>" + escHtml(createdDate) + "</td>";
    html += '<td><button class="btn-details" onclick="showDetail(' + i + ')">Detalles...</button></td>';
    html += "</tr>";
  }

  tbody.innerHTML = html;
}

// =============================================
// Format date helper
// =============================================
function formatDate(dateStr) {
  if (!dateStr) return "—";
  try {
    const d = new Date(dateStr);
    if (isNaN(d.getTime())) return dateStr;
    return d.toLocaleDateString("es-CO", {
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
    });
  } catch (e) {
    return dateStr;
  }
}

// =============================================
// Detail Modal
// =============================================
const LABEL_MAP = {
  id: "ID Asociación",
  public_call_id: "ID del Proceso",
  user_id: "ID Usuario",
  association_date: "Fecha de Asociación",
  created_at: "Fecha de Creación",
  updated_at: "Última Actualización",
  deleted_at: "Fecha de Eliminación",
};

function showDetail(index) {
  const p = savedCache[index];
  if (!p) return;

  document.getElementById("modalTitle").textContent =
    "Detalle — ID Proceso: " + (p.public_call_id != null ? p.public_call_id : "—");

  const grid = document.getElementById("detailGrid");
  grid.innerHTML = "";

  const fields = ["id", "public_call_id", "user_id", "association_date", "created_at", "updated_at", "deleted_at"];

  for (const field of fields) {
    const label = LABEL_MAP[field] || field;
    let value = p[field];

    if (value == null || value === "") value = "—";

    if (field === "association_date" || field === "created_at" || field === "updated_at" || field === "deleted_at") {
      value = formatDate(value);
    }

    const isLong = field === "public_call_id";
    const div = document.createElement("div");
    div.className = "detail-item" + (isLong ? " full-width" : "");
    div.innerHTML = '<div class="label">' + label + '</div><div class="value">' + escHtml(String(value)) + "</div>";
    grid.appendChild(div);
  }

  document.getElementById("detailModal").classList.add("active");
}

function closeDetail() {
  document.getElementById("detailModal").classList.remove("active");
}

// =============================================
// Loading overlay
// =============================================
function showLoading() {
  document.getElementById("loadingOverlay").classList.add("active");
}

function hideLoading() {
  document.getElementById("loadingOverlay").classList.remove("active");
}

// =============================================
// Helpers
// =============================================
function escHtml(str) {
  if (typeof str !== "string") return String(str || "");
  return str
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;");
}