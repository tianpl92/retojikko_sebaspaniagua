package domain

import (
	"database/sql"
	"time"
)

// PublicCallProposal maps to the public_calls_proposals table.
type PublicCallProposal struct {
	ID                          string          `json:"id"`
	NombreDelProcedimiento      sql.NullString  `json:"nombre_del_procedimiento"`
	Entidad                     sql.NullString  `json:"entidad"`
	NitEntidad                  sql.NullString  `json:"nit_entidad"`
	DepartamentoEntidad         sql.NullString  `json:"departamento_entidad"`
	CiudadEntidad               sql.NullString  `json:"ciudad_entidad"`
	OrdenEntidad                sql.NullString  `json:"ordenentidad"`
	ReferenciaDelProceso        sql.NullString  `json:"referencia_del_proceso"`
	DescripcionDelProcedimiento sql.NullString  `json:"descripci_n_del_procedimiento"`
	Fase                        sql.NullString  `json:"fase"`
	FechaDePublicacionDel       sql.NullTime    `json:"fecha_de_publicacion_del"`
	FechaDeUltimaPublicaci      sql.NullTime    `json:"fecha_de_ultima_publicaci"`
	ModalidadDeContratacion     sql.NullString  `json:"modalidad_de_contratacion"`
	PrecioBase                  sql.NullFloat64 `json:"precio_base"`
	Duracion                    sql.NullString  `json:"duracion"`
	UnidadDeDuracion            sql.NullString  `json:"unidad_de_duracion"`
	FechaDeRecepcionDe          sql.NullTime    `json:"fecha_de_recepcion_de"`
	EstadoDelProcedimiento      sql.NullString  `json:"estado_del_procedimiento"`
	Adjudicado                  sql.NullString  `json:"adjudicado"`
	NombreDelProveedor          sql.NullString  `json:"nombre_del_proveedor"`
	ValorTotalAdjudicacion      sql.NullFloat64 `json:"valor_total_adjudicacion"`
	URLProceso                  sql.NullString  `json:"urlproceso"`
	CodigoPrincipalDeCategoria  sql.NullString  `json:"codigo_principal_de_categoria"`
	TipoDeContrato              sql.NullString  `json:"tipo_de_contrato"`
	EstadoDeAperturaDelProceso  sql.NullString  `json:"estado_de_apertura_del_proceso"`
	EstadoResumen               sql.NullString  `json:"estado_resumen"`
	ProveedoresInvitados        sql.NullString  `json:"proveedores_invitados"`
	ProveedoresQueManifestaron  sql.NullString  `json:"proveedores_que_manifestaron"`
	RespuestasAlProcedimiento   sql.NullString  `json:"respuestas_al_procedimiento"`
	NumeroDeLotes               sql.NullInt64   `json:"numero_de_lotes"`
	CreatedAt                   time.Time       `json:"created_at"`
	UpdatedAt                   time.Time       `json:"updated_at"`
	DeletedAt                   *time.Time      `json:"deleted_at,omitempty"`
}
