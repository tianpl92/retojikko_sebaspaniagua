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
      '<tr><td colspan="5" style="text-align:center;padding:40px;color:#7a7a8e;">No hay convocatorias guardadas</td></tr>';
    return;
  }

  let html = "";
  for (let i = 0; i < items.length; i++) {
    const p = items[i];
    const estado = p.estado_del_procedimiento || "—";
    const badgeClass = badgeForStatus(estado);
    const safeNombre = escHtml(p.nombre_del_procedimiento || "—");
    const safeEntidad = escHtml(p.entidad || "—");
    const safeId = escHtml(p.id_del_proceso || (p.public_call_id != null ? String(p.public_call_id) : "—"));

    html += "<tr>";
    html += "<td>" + safeNombre + "</td>";
    html += "<td>" + safeEntidad + "</td>";
    html += "<td>" + safeId + "</td>";
    html += '<td><span class="badge ' + badgeClass + '">' + escHtml(estado) + "</span></td>";
    html += '<td><button class="btn-details" onclick="showDetail(' + i + ')">Detalles...</button></td>';
    html += "</tr>";
  }

  tbody.innerHTML = html;
}

// =============================================
// Detail Modal
// =============================================
const LABEL_MAP = {
  nombre_del_procedimiento: "Nombre del Procedimiento",
  entidad: "Entidad",
  nit_entidad: "NIT Entidad",
  departamento_entidad: "Departamento",
  ciudad_entidad: "Ciudad",
  ordenentidad: "Orden",
  id_del_proceso: "ID del Proceso",
  referencia_del_proceso: "Referencia",
  descripci_n_del_procedimiento: "Descripción",
  fase: "Fase",
  fecha_de_publicacion_del: "Fecha de Publicación",
  fecha_de_ultima_publicaci: "Última Publicación",
  fecha_de_recepcion_de: "Fecha de Recepción",
  modalidad_de_contratacion: "Modalidad",
  precio_base: "Precio Base",
  duracion: "Duración",
  unidad_de_duracion: "Unidad de Duración",
  estado_del_procedimiento: "Estado",
  estado_de_apertura_del_proceso: "Estado de Apertura",
  estado_resumen: "Estado Resumen",
  adjudicado: "Adjudicado",
  nombre_del_proveedor: "Proveedor",
  valor_total_adjudicacion: "Valor Total",
  codigo_principal_de_categoria: "Categoría",
  tipo_de_contrato: "Tipo de Contrato",
  urlproceso: "URL del Proceso",
  proveedores_invitados: "Proveedores Invitados",
  proveedores_que_manifestaron: "Proveedores Manifestados",
  respuestas_al_procedimiento: "Respuestas",
  numero_de_lotes: "Número de Lotes",
};

function showDetail(index) {
  const p = savedCache[index];
  if (!p) return;

  document.getElementById("modalTitle").textContent =
    p.nombre_del_procedimiento || "Detalle — ID Proceso: " + (p.id_del_proceso || "—");

  const grid = document.getElementById("detailGrid");
  grid.innerHTML = "";

  const fields = [
    "nombre_del_procedimiento", "entidad", "nit_entidad", "departamento_entidad",
    "ciudad_entidad", "ordenentidad", "id_del_proceso", "referencia_del_proceso",
    "descripci_n_del_procedimiento", "fase", "fecha_de_publicacion_del",
    "fecha_de_ultima_publicaci", "fecha_de_recepcion_de",
    "modalidad_de_contratacion", "precio_base", "duracion", "unidad_de_duracion",
    "estado_del_procedimiento", "estado_de_apertura_del_proceso", "estado_resumen",
    "adjudicado", "nombre_del_proveedor", "valor_total_adjudicacion",
    "codigo_principal_de_categoria", "tipo_de_contrato", "urlproceso",
    "proveedores_invitados", "proveedores_que_manifestaron",
    "respuestas_al_procedimiento", "numero_de_lotes",
  ];

  for (const field of fields) {
    const label = LABEL_MAP[field] || field;
    let value = p[field];

    if (value == null || value === "") value = "—";

    if (field === "precio_base" && value !== "—") {
      const num = Number(value);
      value = "$" + num.toLocaleString("es-CO", { minimumFractionDigits: 2 });
    }

    if (field === "valor_total_adjudicacion" && value !== "—") {
      const num = Number(value);
      value = "$" + num.toLocaleString("es-CO", { minimumFractionDigits: 2 });
    }

    if (field === "urlproceso" && value !== "—") {
      value = '<a href="' + escHtml(value) + '" target="_blank" rel="noopener">Ver en SECOP II &rarr;</a>';
    }

    const isLong = field === "descripci_n_del_procedimiento" || field === "nombre_del_procedimiento";
    const div = document.createElement("div");
    div.className = "detail-item" + (isLong ? " full-width" : "");
    div.innerHTML = '<div class="label">' + label + '</div><div class="value">' + value + "</div>";
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
function badgeForStatus(estado) {
  const s = (estado || "").toLowerCase();
  if (s.includes("public") || s.includes("abierto") || s.includes("presentaci"))
    return "badge-published";
  if (s.includes("cerrado") || s.includes("adjudicado") || s.includes("terminado"))
    return "badge-closed";
  return "badge-draft";
}

function escHtml(str) {
  if (typeof str !== "string") return String(str || "");
  return str
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;");
}