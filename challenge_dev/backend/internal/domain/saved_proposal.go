package domain

import (
	"fmt"
	"time"
)

// SavedProposal maps to the public_call_user_associations table.
type SavedProposal struct {
	ID              string     `json:"id"`
	PublicCallID    string     `json:"public_call_id"`
	UserID          string     `json:"user_id"`
	AssociationDate time.Time  `json:"association_date"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

// SavedProposalWithProposal combines SavedProposal with the full PublicCallProposal data.
type SavedProposalWithProposal struct {
	ID                                  string     `json:"id"`
	PublicCallID                        string     `json:"public_call_id"`
	UserID                              string     `json:"user_id"`
	AssociationDate                     time.Time  `json:"association_date"`
	CreatedAt                           time.Time  `json:"created_at"`
	UpdatedAt                           time.Time  `json:"updated_at"`
	DeletedAt                           *time.Time `json:"deleted_at,omitempty"`
	ProposalNombreDelProcedimiento      string     `json:"nombre_del_procedimiento,omitempty"`
	ProposalEntidad                     string     `json:"entidad,omitempty"`
	ProposalNitEntidad                  string     `json:"nit_entidad,omitempty"`
	ProposalDepartamentoEntidad         string     `json:"departamento_entidad,omitempty"`
	ProposalCiudadEntidad               string     `json:"ciudad_entidad,omitempty"`
	ProposalOrdenEntidad                string     `json:"ordenentidad,omitempty"`
	ProposalReferenciaDelProceso        string     `json:"referencia_del_proceso,omitempty"`
	ProposalDescripcionDelProcedimiento string     `json:"descripci_n_del_procedimiento,omitempty"`
	ProposalFase                        string     `json:"fase,omitempty"`
	ProposalFechaDePublicacionDel       string     `json:"fecha_de_publicacion_del,omitempty"`
	ProposalFechaDeUltimaPublicaci      string     `json:"fecha_de_ultima_publicaci,omitempty"`
	ProposalModalidadDeContratacion     string     `json:"modalidad_de_contratacion,omitempty"`
	ProposalPrecioBase                  float64    `json:"precio_base,omitempty"`
	ProposalDuracion                    string     `json:"duracion,omitempty"`
	ProposalUnidadDeDuracion            string     `json:"unidad_de_duracion,omitempty"`
	ProposalFechaDeRecepcionDe          string     `json:"fecha_de_recepcion_de,omitempty"`
	ProposalEstadoDelProcedimiento      string     `json:"estado_del_procedimiento,omitempty"`
	ProposalAdjudicado                  string     `json:"adjudicado,omitempty"`
	ProposalNombreDelProveedor          string     `json:"nombre_del_proveedor,omitempty"`
	ProposalValorTotalAdjudicacion      float64    `json:"valor_total_adjudicacion,omitempty"`
	ProposalURLProceso                  string     `json:"urlproceso,omitempty"`
	ProposalCodigoPrincipalDeCategoria  string     `json:"codigo_principal_de_categoria,omitempty"`
	ProposalTipoDeContrato              string     `json:"tipo_de_contrato,omitempty"`
	ProposalEstadoDeAperturaDelProceso  string     `json:"estado_de_apertura_del_proceso,omitempty"`
	ProposalEstadoResumen               string     `json:"estado_resumen,omitempty"`
	ProposalProveedoresInvitados        string     `json:"proveedores_invitados,omitempty"`
	ProposalProveedoresQueManifestaron  string     `json:"proveedores_que_manifestaron,omitempty"`
	ProposalRespuestasAlProcedimiento   string     `json:"respuestas_al_procedimiento,omitempty"`
	ProposalNumeroDeLotes               int64      `json:"numero_de_lotes,omitempty"`
}

// SavedProposalFullData is the public response shape for GET /saved_proposals.
// It includes user_id plus all proposal fields from public_calls_proposals,
// with the proposal ID serialised as id_del_proceso.
type SavedProposalFullData struct {
	UserID                      string  `json:"user_id"`
	Entidad                     string  `json:"entidad"`
	NitEntidad                  *string `json:"nit_entidad,omitempty"`
	DepartamentoEntidad         *string `json:"departamento_entidad,omitempty"`
	CiudadEntidad               *string `json:"ciudad_entidad,omitempty"`
	OrdenEntidad                *string `json:"ordenentidad,omitempty"`
	IDDelProceso                string  `json:"id_del_proceso"`
	ReferenciaDelProceso        *string `json:"referencia_del_proceso,omitempty"`
	NombreDelProcedimiento      string  `json:"nombre_del_procedimiento"`
	DescripcionDelProcedimiento *string `json:"descripci_n_del_procedimiento,omitempty"`
	Fase                        string  `json:"fase"`
	FechaDePublicacionDel       *string `json:"fecha_de_publicacion_del,omitempty"`
	FechaDeUltimaPublicaci      *string `json:"fecha_de_ultima_publicaci,omitempty"`
	ModalidadDeContratacion     *string `json:"modalidad_de_contratacion,omitempty"`
	PrecioBase                  *string `json:"precio_base,omitempty"`
	Duracion                    *string `json:"duracion,omitempty"`
	UnidadDeDuracion            *string `json:"unidad_de_duracion,omitempty"`
	FechaDeRecepcionDe          *string `json:"fecha_de_recepcion_de,omitempty"`
	EstadoDelProcedimiento      *string `json:"estado_del_procedimiento,omitempty"`
	Adjudicado                  *string `json:"adjudicado,omitempty"`
	NombreDelProveedor          *string `json:"nombre_del_proveedor,omitempty"`
	ValorTotalAdjudicacion      *string `json:"valor_total_adjudicacion,omitempty"`
	URLProceso                  *string `json:"urlproceso,omitempty"`
	CodigoPrincipalDeCategoria  *string `json:"codigo_principal_de_categoria,omitempty"`
	TipoDeContrato              *string `json:"tipo_de_contrato,omitempty"`
	EstadoDeAperturaDelProceso  *string `json:"estado_de_apertura_del_proceso,omitempty"`
	EstadoResumen               *string `json:"estado_resumen,omitempty"`
	ProveedoresInvitados        *string `json:"proveedores_invitados,omitempty"`
	ProveedoresQueManifestaron  *string `json:"proveedores_que_manifestaron,omitempty"`
	RespuestasAlProcedimiento   *string `json:"respuestas_al_procedimiento,omitempty"`
	NumeroDeLotes               *int64  `json:"numero_de_lotes,omitempty"`
}

// ToFullData converts a SavedProposalWithProposal into the SavedProposalFullData response shape.
func (s *SavedProposalWithProposal) ToFullData() *SavedProposalFullData {
	return &SavedProposalFullData{
		UserID:                      s.UserID,
		Entidad:                     s.ProposalEntidad,
		NitEntidad:                  strPtr(s.ProposalNitEntidad),
		DepartamentoEntidad:         strPtr(s.ProposalDepartamentoEntidad),
		CiudadEntidad:               strPtr(s.ProposalCiudadEntidad),
		OrdenEntidad:                strPtr(s.ProposalOrdenEntidad),
		IDDelProceso:                s.PublicCallID,
		ReferenciaDelProceso:        strPtr(s.ProposalReferenciaDelProceso),
		NombreDelProcedimiento:      s.ProposalNombreDelProcedimiento,
		DescripcionDelProcedimiento: strPtr(s.ProposalDescripcionDelProcedimiento),
		Fase:                        s.ProposalFase,
		FechaDePublicacionDel:       strPtr(s.ProposalFechaDePublicacionDel),
		FechaDeUltimaPublicaci:      strPtr(s.ProposalFechaDeUltimaPublicaci),
		ModalidadDeContratacion:     strPtr(s.ProposalModalidadDeContratacion),
		PrecioBase:                  float64StrPtr(s.ProposalPrecioBase),
		Duracion:                    strPtr(s.ProposalDuracion),
		UnidadDeDuracion:            strPtr(s.ProposalUnidadDeDuracion),
		FechaDeRecepcionDe:          strPtr(s.ProposalFechaDeRecepcionDe),
		EstadoDelProcedimiento:      strPtr(s.ProposalEstadoDelProcedimiento),
		Adjudicado:                  strPtr(s.ProposalAdjudicado),
		NombreDelProveedor:          strPtr(s.ProposalNombreDelProveedor),
		ValorTotalAdjudicacion:      float64StrPtr(s.ProposalValorTotalAdjudicacion),
		URLProceso:                  strPtr(s.ProposalURLProceso),
		CodigoPrincipalDeCategoria:  strPtr(s.ProposalCodigoPrincipalDeCategoria),
		TipoDeContrato:              strPtr(s.ProposalTipoDeContrato),
		EstadoDeAperturaDelProceso:  strPtr(s.ProposalEstadoDeAperturaDelProceso),
		EstadoResumen:               strPtr(s.ProposalEstadoResumen),
		ProveedoresInvitados:        strPtr(s.ProposalProveedoresInvitados),
		ProveedoresQueManifestaron:  strPtr(s.ProposalProveedoresQueManifestaron),
		RespuestasAlProcedimiento:   strPtr(s.ProposalRespuestasAlProcedimiento),
		NumeroDeLotes:               int64Ptr(s.ProposalNumeroDeLotes),
	}
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func float64StrPtr(f float64) *string {
	if f == 0 {
		return nil
	}
	s := formatFloat64(f)
	return &s
}

func int64Ptr(i int64) *int64 {
	if i == 0 {
		return nil
	}
	return &i
}

func formatFloat64(f float64) string {
	if f == float64(int64(f)) {
		return fmt.Sprintf("%d", int64(f))
	}
	return fmt.Sprintf("%.2f", f)
}
