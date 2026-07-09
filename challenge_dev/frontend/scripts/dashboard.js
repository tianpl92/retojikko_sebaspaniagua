// Public Calls Portal — Dashboard Logic

const PAGE_SIZE = 100;
let currentPage = 1;
let proposalsCache = [];
let totalCount = 0;
let hasMore = false;

// =============================================
// Init — runs on page load
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

  // Wire filter button
  document.getElementById("btnFiltrar").addEventListener("click", function () {
    currentPage = 1;
    loadProposals();
  });

  // Wire pagination buttons
  document.getElementById("btnPrev").addEventListener("click", function () {
    if (currentPage > 1) {
      currentPage--;
      loadProposals();
    }
  });
  document.getElementById("btnNext").addEventListener("click", function () {
    if (hasMore) {
      currentPage++;
      loadProposals();
    }
  });

  // Wire Guardar button
  document.getElementById("btnGuardar").addEventListener("click", saveSelectedProposals);

  // Checkbox change delegation — enable/disable Guardar button
  document.getElementById("proposalsBody").addEventListener("change", function (e) {
    if (e.target.type === "checkbox") {
      updateGuardarButton();
    }
  });

  // Load first page
  await loadProposals();
});

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
// Fetch proposals from backend
// =============================================
async function loadProposals() {
  const token = localStorage.getItem("jwt_token");
  if (!token) {
    window.location.href = "login.html";
    return;
  }

  // Clear grid and show loading
  document.getElementById("proposalsBody").innerHTML =
    '<tr><td colspan="6" style="text-align:center;padding:40px;color:#7a7a8e;">Buscando...</td></tr>';
  showLoading();

  // Build filters
  const nombre = document.getElementById("filtroNombre").value.trim();
  const codigo = document.getElementById("filtroCodigo").value.trim();
  const entidad = document.getElementById("filtroEntidad").value.trim();
  const nit = document.getElementById("filtroNit").value.trim();
  const ciudad = document.getElementById("filtroCiudad").value.trim();

  // Combine text filters into query param (backend searches nombre_del_procedimiento)
  const queryParts = [];
  if (nombre) queryParts.push(nombre);
  if (codigo) queryParts.push(codigo);
  if (nit) queryParts.push(nit);
  if (ciudad) queryParts.push(ciudad);
  const query = queryParts.join(" ");

  // Request 1 extra item to detect "more pages"
  const params = new URLSearchParams();
  params.set("limit", PAGE_SIZE + 1);
  params.set("offset", (currentPage - 1) * PAGE_SIZE);
  if (query) params.set("query", query);
  if (entidad) params.set("entidad", entidad);

  try {
    const response = await fetch(
      BACKEND_URL + ENDPOINTS.PROPOSALS + "?" + params.toString(),
      {
        headers: { Authorization: "Bearer " + token },
      }
    );

    if (response.status === 401) {
      localStorage.clear();
      window.location.href = "login.html";
      return;
    }

    if (!response.ok) {
      hideLoading();
      document.getElementById("proposalsBody").innerHTML =
        '<tr><td colspan="6" style="text-align:center;padding:40px;color:#a94442;">Error al cargar convocatorias</td></tr>';
      updatePagination(0);
      return;
    }

    const data = await response.json();

    // If we got PAGE_SIZE+1 items, there are more pages
    hasMore = data.length > PAGE_SIZE;
    const items = data.slice(0, PAGE_SIZE);
    proposalsCache = items;

    hideLoading();
    renderGrid(items);
    updatePagination(items.length);
  } catch (err) {
    hideLoading();
    document.getElementById("proposalsBody").innerHTML =
      '<tr><td colspan="6" style="text-align:center;padding:40px;color:#a94442;">Error de conexión con el servidor</td></tr>';
    updatePagination(0);
  }
}

// =============================================
// Render table rows
// =============================================
function renderGrid(items) {
  const tbody = document.getElementById("proposalsBody");

  if (!items || items.length === 0) {
    tbody.innerHTML =
      '<tr><td colspan="6" style="text-align:center;padding:40px;color:#7a7a8e;">No se encontraron convocatorias</td></tr>';
    return;
  }

  let html = "";
  for (let i = 0; i < items.length; i++) {
    const p = items[i];
    const estado = p.estado_del_procedimiento || "";
    const badgeClass = badgeForStatus(estado);
    const safeNombre = escHtml(p.nombre_del_procedimiento || "—");
    const safeEntidad = escHtml(p.entidad || "—");
    const safeId = escHtml(p.id_del_proceso || "—");

    html += "<tr>";
    html += '<td><input type="checkbox" data-index="' + i + '"></td>';
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
// Pagination controls
// =============================================
function updatePagination(count) {
  const info = document.getElementById("pageInfo");
  const btnPrev = document.getElementById("btnPrev");
  const btnNext = document.getElementById("btnNext");
  const pageNum = document.getElementById("pageNum");

  const start = totalCount === 0 && count > 0 ? 1 : (currentPage - 1) * PAGE_SIZE + 1;
  const end = (currentPage - 1) * PAGE_SIZE + count;

  info.textContent = "Mostrando " + start + "–" + end + " de " + (hasMore ? (end + "+") : end) + " resultados";
  btnPrev.disabled = currentPage <= 1;
  btnNext.disabled = !hasMore;
  pageNum.textContent = "Página " + currentPage;
}

// =============================================
// Detail Modal
// =============================================
const LABEL_MAP = {
  entidad: "Entidad",
  nit_entidad: "NIT Entidad",
  departamento_entidad: "Departamento",
  ciudad_entidad: "Ciudad",
  ordenentidad: "Orden",
  id_del_proceso: "ID del Proceso",
  referencia_del_proceso: "Referencia",
  nombre_del_procedimiento: "Nombre del Procedimiento",
  descripci_n_del_procedimiento: "Descripción",
  fase: "Fase",
  fecha_de_publicacion_del: "Fecha de Publicación",
  fecha_de_ultima_publicaci: "Última Publicación",
  modalidad_de_contratacion: "Modalidad",
  precio_base: "Precio Base",
  estado_del_procedimiento: "Estado",
  urlproceso: "URL del Proceso",
};

function showDetail(index) {
  const p = proposalsCache[index];
  if (!p) return;

  document.getElementById("modalTitle").textContent = p.nombre_del_procedimiento || "Detalle de Convocatoria";

  const grid = document.getElementById("detailGrid");
  grid.innerHTML = "";

  const fields = [
    "entidad", "nit_entidad", "departamento_entidad", "ciudad_entidad",
    "ordenentidad", "id_del_proceso", "referencia_del_proceso",
    "nombre_del_procedimiento", "descripci_n_del_procedimiento",
    "fase", "fecha_de_publicacion_del", "fecha_de_ultima_publicaci",
    "modalidad_de_contratacion", "precio_base", "estado_del_procedimiento",
    "urlproceso",
  ];

  for (const field of fields) {
    const label = LABEL_MAP[field] || field;
    let value = p[field];

    if (value == null || value === "") value = "—";

    if (field === "precio_base" && value !== "—") {
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

// Close modal on overlay click
document.addEventListener("DOMContentLoaded", function () {
  document.getElementById("detailModal").addEventListener("click", function (e) {
    if (e.target === this) closeDetail();
  });
});

// Close modal on Escape
document.addEventListener("keydown", function (e) {
  if (e.key === "Escape") closeDetail();
});

// =============================================
// Save proposals
// =============================================
function updateGuardarButton() {
  const btn = document.getElementById("btnGuardar");
  const checked = document.querySelectorAll('#proposalsBody input[type="checkbox"]:checked');
  btn.disabled = checked.length === 0;
}

async function saveSelectedProposals() {
  const token = localStorage.getItem("jwt_token");
  if (!token) {
    window.location.href = "login.html";
    return;
  }

  const checkedBoxes = document.querySelectorAll('#proposalsBody input[type="checkbox"]:checked');
  if (checkedBoxes.length === 0) return;

  // Show loading overlay with save message
  const overlay = document.getElementById("loadingOverlay");
  const loadText = overlay.querySelector(".loading-text");
  loadText.textContent = "Guardando " + checkedBoxes.length + " convocatorias...";
  overlay.classList.add("active");

  let successCount = 0;
  let errorCount = 0;

  for (const box of checkedBoxes) {
    const index = parseInt(box.getAttribute("data-index"), 10);
    const proposal = proposalsCache[index];
    if (!proposal || !proposal.id_del_proceso) continue;

    try {
      // Step 1: Save full proposal data locally
      // Map id_del_proceso → id (backend expects "id" field)
      const savePayload = { ...proposal, id: proposal.id_del_proceso };
      const saveResp = await fetch(BACKEND_URL + "/proposal-save", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: "Bearer " + token,
        },
        body: JSON.stringify(savePayload),
      });

      if (!saveResp.ok && saveResp.status !== 400) {
        errorCount++;
        continue;
      }

      // Step 2: Create user-proposal association
      const assocResp = await fetch(BACKEND_URL + ENDPOINTS.SAVED_SAVE, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: "Bearer " + token,
        },
        body: JSON.stringify({ public_call_id: proposal.id_del_proceso }),
      });

      if (assocResp.ok || assocResp.status === 400) {
        successCount++;
      } else {
        errorCount++;
      }
    } catch (err) {
      errorCount++;
    }
  }

  // Hide loading overlay
  overlay.classList.remove("active");
  loadText.textContent = "Cargando convocatorias...";

  // Uncheck all boxes and update button
  checkedBoxes.forEach(function (b) { b.checked = false; });
  updateGuardarButton();

  // Show result
  if (errorCount === 0) {
    alert("Guardado satisfactoriamente");
  } else {
    alert("Guardadas: " + successCount + ", errores: " + errorCount);
  }
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