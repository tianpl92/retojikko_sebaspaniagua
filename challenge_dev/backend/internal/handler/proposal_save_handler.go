package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/domain"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/service"
)

// proposalSaveRequest accepts all proposal fields matching the datos.gov.co response shape.
// Most fields are strings because datos.gov.co returns everything as strings.
// URLProceso is interface{} because datos.gov.co returns it as {"url": "..."}.
type proposalSaveRequest struct {
	ID                          string      `json:"id"`
	NombreDelProcedimiento      string      `json:"nombre_del_procedimiento"`
	Entidad                     string      `json:"entidad"`
	NitEntidad                  string      `json:"nit_entidad"`
	DepartamentoEntidad         string      `json:"departamento_entidad"`
	CiudadEntidad               string      `json:"ciudad_entidad"`
	OrdenEntidad                string      `json:"ordenentidad"`
	ReferenciaDelProceso        string      `json:"referencia_del_proceso"`
	DescripcionDelProcedimiento string      `json:"descripci_n_del_procedimiento"`
	Fase                        string      `json:"fase"`
	FechaDePublicacionDel       string      `json:"fecha_de_publicacion_del"`
	FechaDeUltimaPublicaci      string      `json:"fecha_de_ultima_publicaci"`
	ModalidadDeContratacion     string      `json:"modalidad_de_contratacion"`
	PrecioBase                  string      `json:"precio_base"`
	Duracion                    string      `json:"duracion"`
	UnidadDeDuracion            string      `json:"unidad_de_duracion"`
	FechaDeRecepcionDe          string      `json:"fecha_de_recepcion_de"`
	EstadoDelProcedimiento      string      `json:"estado_del_procedimiento"`
	Adjudicado                  string      `json:"adjudicado"`
	NombreDelProveedor          string      `json:"nombre_del_proveedor"`
	ValorTotalAdjudicacion      string      `json:"valor_total_adjudicacion"`
	URLProceso                  interface{} `json:"urlproceso"`
	CodigoPrincipalDeCategoria  string      `json:"codigo_principal_de_categoria"`
	TipoDeContrato              string      `json:"tipo_de_contrato"`
	EstadoDeAperturaDelProceso  string      `json:"estado_de_apertura_del_proceso"`
	EstadoResumen               string      `json:"estado_resumen"`
	ProveedoresInvitados        string      `json:"proveedores_invitados"`
	ProveedoresQueManifestaron  string      `json:"proveedores_que_manifestaron"`
	RespuestasAlProcedimiento   string      `json:"respuestas_al_procedimiento"`
	NumeroDeLotes               string      `json:"numero_de_lotes"`
}

// toDomain converts the request to a PublicCallProposal domain object.
func (r *proposalSaveRequest) toDomain() *domain.PublicCallProposal {
	now := time.Now()

	return &domain.PublicCallProposal{
		ID:                          r.ID,
		NombreDelProcedimiento:      toNullString(r.NombreDelProcedimiento),
		Entidad:                     toNullString(r.Entidad),
		NitEntidad:                  toNullString(r.NitEntidad),
		DepartamentoEntidad:         toNullString(r.DepartamentoEntidad),
		CiudadEntidad:               toNullString(r.CiudadEntidad),
		OrdenEntidad:                toNullString(r.OrdenEntidad),
		ReferenciaDelProceso:        toNullString(r.ReferenciaDelProceso),
		DescripcionDelProcedimiento: toNullString(r.DescripcionDelProcedimiento),
		Fase:                        toNullString(r.Fase),
		FechaDePublicacionDel:       toNullTime(r.FechaDePublicacionDel),
		FechaDeUltimaPublicaci:      toNullTime(r.FechaDeUltimaPublicaci),
		ModalidadDeContratacion:     toNullString(r.ModalidadDeContratacion),
		PrecioBase:                  parseNullFloat64(r.PrecioBase),
		Duracion:                    toNullString(r.Duracion),
		UnidadDeDuracion:            toNullString(r.UnidadDeDuracion),
		FechaDeRecepcionDe:          toNullTime(r.FechaDeRecepcionDe),
		EstadoDelProcedimiento:      toNullString(r.EstadoDelProcedimiento),
		Adjudicado:                  toNullString(r.Adjudicado),
		NombreDelProveedor:          toNullString(r.NombreDelProveedor),
		ValorTotalAdjudicacion:      parseNullFloat64(r.ValorTotalAdjudicacion),
		URLProceso:                  extractURLString(r.URLProceso),
		CodigoPrincipalDeCategoria:  toNullString(r.CodigoPrincipalDeCategoria),
		TipoDeContrato:              toNullString(r.TipoDeContrato),
		EstadoDeAperturaDelProceso:  toNullString(r.EstadoDeAperturaDelProceso),
		EstadoResumen:               toNullString(r.EstadoResumen),
		ProveedoresInvitados:        toNullString(r.ProveedoresInvitados),
		ProveedoresQueManifestaron:  toNullString(r.ProveedoresQueManifestaron),
		RespuestasAlProcedimiento:   toNullString(r.RespuestasAlProcedimiento),
		NumeroDeLotes:               parseNullInt64(r.NumeroDeLotes),
		CreatedAt:                   now,
		UpdatedAt:                   now,
	}
}

// ProposalSaveHandler handles POST /proposal-save.
type ProposalSaveHandler struct {
	proposalService *service.ProposalService
}

func NewProposalSaveHandler(proposalService *service.ProposalService) *ProposalSaveHandler {
	return &ProposalSaveHandler{proposalService: proposalService}
}

func (h *ProposalSaveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req proposalSaveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.ID == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	proposal := req.toDomain()
	saved, err := h.proposalService.SaveProposal(r.Context(), proposal)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Return clean JSON (sql.Null* types don't marshal well)
	clean := map[string]interface{}{
		"id":                               req.ID,
		"nombre_del_procedimiento":         req.NombreDelProcedimiento,
		"entidad":                          req.Entidad,
		"nit_entidad":                      req.NitEntidad,
		"departamento_entidad":             req.DepartamentoEntidad,
		"ciudad_entidad":                   req.CiudadEntidad,
		"ordenentidad":                     req.OrdenEntidad,
		"referencia_del_proceso":           req.ReferenciaDelProceso,
		"descripci_n_del_procedimiento":    req.DescripcionDelProcedimiento,
		"fase":                             req.Fase,
		"fecha_de_publicacion_del":         req.FechaDePublicacionDel,
		"fecha_de_ultima_publicaci":        req.FechaDeUltimaPublicaci,
		"modalidad_de_contratacion":        req.ModalidadDeContratacion,
		"precio_base":                      req.PrecioBase,
		"duracion":                         req.Duracion,
		"unidad_de_duracion":               req.UnidadDeDuracion,
		"fecha_de_recepcion_de":            req.FechaDeRecepcionDe,
		"estado_del_procedimiento":         req.EstadoDelProcedimiento,
		"adjudicado":                       req.Adjudicado,
		"nombre_del_proveedor":             req.NombreDelProveedor,
		"valor_total_adjudicacion":         req.ValorTotalAdjudicacion,
		"urlproceso":                       extractURLStringForResponse(req.URLProceso),
		"codigo_principal_de_categoria":    req.CodigoPrincipalDeCategoria,
		"tipo_de_contrato":                 req.TipoDeContrato,
		"estado_de_apertura_del_proceso":   req.EstadoDeAperturaDelProceso,
		"estado_resumen":                   req.EstadoResumen,
		"proveedores_invitados":            req.ProveedoresInvitados,
		"proveedores_que_manifestaron":     req.ProveedoresQueManifestaron,
		"respuestas_al_procedimiento":      req.RespuestasAlProcedimiento,
		"numero_de_lotes":                  req.NumeroDeLotes,
	}
	if saved != nil {
		clean["created_at"] = saved.CreatedAt.Format(time.RFC3339)
		clean["updated_at"] = saved.UpdatedAt.Format(time.RFC3339)
	}

	writeJSON(w, http.StatusCreated, clean)
}

// extractURLString handles urlproceso which can be a string or {"url": "..."} object.
func extractURLString(v interface{}) sql.NullString {
	if v == nil {
		return sql.NullString{Valid: false}
	}
	switch val := v.(type) {
	case string:
		if val == "" || val == "null" || val == "No Definido" {
			return sql.NullString{Valid: false}
		}
		return sql.NullString{String: val, Valid: true}
	case map[string]interface{}:
		if u, ok := val["url"]; ok {
			if s, ok := u.(string); ok && s != "" {
				return sql.NullString{String: s, Valid: true}
			}
		}
		return sql.NullString{Valid: false}
	}
	return sql.NullString{Valid: false}
}

// extractURLStringForResponse extracts a plain string from urlproceso for JSON responses.
func extractURLStringForResponse(v interface{}) string {
	ns := extractURLString(v)
	if ns.Valid {
		return ns.String
	}
	return ""
}

func toNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

// parseNullFloat64 accepts a string (as returned by datos.gov.co) and converts to sql.NullFloat64.
func parseNullFloat64(s string) sql.NullFloat64 {
	s = strings.TrimSpace(s)
	if s == "" || s == "null" {
		return sql.NullFloat64{Valid: false}
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return sql.NullFloat64{Valid: false}
	}
	return sql.NullFloat64{Float64: f, Valid: true}
}

// parseNullInt64 accepts a string (as returned by datos.gov.co) and converts to sql.NullInt64.
func parseNullInt64(s string) sql.NullInt64 {
	s = strings.TrimSpace(s)
	if s == "" || s == "null" {
		return sql.NullInt64{Valid: false}
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: i, Valid: true}
}

func toNullTime(s string) sql.NullTime {
	if s == "" {
		return sql.NullTime{Valid: false}
	}
	// Try several time formats
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05.000",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return sql.NullTime{Time: t, Valid: true}
		}
	}
	return sql.NullTime{Valid: false}
}