package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/domain"
)

// PostgresSavedProposalRepository is a PostgreSQL-backed implementation of SavedProposalRepository.
type PostgresSavedProposalRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresSavedProposalRepository creates a new PostgresSavedProposalRepository.
func NewPostgresSavedProposalRepository(pool *pgxpool.Pool) *PostgresSavedProposalRepository {
	return &PostgresSavedProposalRepository{pool: pool}
}

// Save creates a new saved proposal association.
func (r *PostgresSavedProposalRepository) Save(ctx context.Context, userID, publicCallID string) (*domain.SavedProposal, error) {
	now := time.Now()
	query := `
		INSERT INTO public_call_user_associations
			(user_id, public_call_id, association_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var id string
	err := r.pool.QueryRow(ctx, query,
		userID, publicCallID, now, now, now,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("save saved proposal: %w", err)
	}

	return &domain.SavedProposal{
		ID:              id,
		PublicCallID:    publicCallID,
		UserID:          userID,
		AssociationDate: now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

// FindByUserID retrieves all saved proposals for a user.
func (r *PostgresSavedProposalRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.SavedProposal, error) {
	query := `
		SELECT id, public_call_id, user_id, association_date,
		       created_at, updated_at, deleted_at
		FROM public_call_user_associations
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("find saved proposals by user: %w", err)
	}
	defer rows.Close()

	var results []*domain.SavedProposal
	for rows.Next() {
		s := &domain.SavedProposal{}
		err := rows.Scan(
			&s.ID, &s.PublicCallID, &s.UserID,
			&s.AssociationDate, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan saved proposal: %w", err)
		}
		results = append(results, s)
	}
	return results, rows.Err()
}

// FindByUserIDWithProposals retrieves all saved proposals with full proposal data via LEFT JOIN.
func (r *PostgresSavedProposalRepository) FindByUserIDWithProposals(ctx context.Context, userID string) ([]*domain.SavedProposalWithProposal, error) {
	query := `
		SELECT
			a.id,
			a.public_call_id,
			a.user_id,
			a.association_date,
			a.created_at,
			a.updated_at,
			a.deleted_at,
			p.id,
			p.nombre_del_procedimiento,
			p.entidad,
			p.nit_entidad,
			p.departamento_entidad,
			p.ciudad_entidad,
			p.ordenentidad,
			p.referencia_del_proceso,
			p.descripci_n_del_procedimiento,
			p.fase,
			p.fecha_de_publicacion_del,
			p.fecha_de_ultima_publicaci,
			p.modalidad_de_contratacion,
			p.precio_base,
			p.duracion,
			p.unidad_de_duracion,
			p.fecha_de_recepcion_de,
			p.estado_del_procedimiento,
			p.adjudicado,
			p.nombre_del_proveedor,
			p.valor_total_adjudicacion,
			p.urlproceso,
			p.codigo_principal_de_categoria,
			p.tipo_de_contrato,
			p.estado_de_apertura_del_proceso,
			p.estado_resumen,
			p.proveedores_invitados,
			p.proveedores_que_manifestaron,
			p.respuestas_al_procedimiento,
			p.numero_de_lotes,
			p.created_at,
			p.updated_at,
			p.deleted_at
		FROM public_call_user_associations a
		LEFT JOIN public_calls_proposals p ON a.public_call_id = p.id AND p.deleted_at IS NULL
		WHERE a.user_id = $1 AND a.deleted_at IS NULL
		ORDER BY a.created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("find saved proposals with proposals: %w", err)
	}
	defer rows.Close()

	var results []*domain.SavedProposalWithProposal
	for rows.Next() {
		item := &domain.SavedProposalWithProposal{}
		var assocDeletedAt *time.Time
		var (
			pID                          *string
			pNombreDelProcedimiento      *string
			pEntidad                     *string
			pNitEntidad                  *string
			pDepartamentoEntidad         *string
			pCiudadEntidad               *string
			pOrdenEntidad                *string
			pReferenciaDelProceso        *string
			pDescripcionDelProcedimiento *string
			pFase                        *string
			pFechaDePublicacionDel       *time.Time
			pFechaDeUltimaPublicaci      *time.Time
			pModalidadDeContratacion     *string
			pPrecioBase                  *float64
			pDuracion                    *string
			pUnidadDeDuracion            *string
			pFechaDeRecepcionDe          *time.Time
			pEstadoDelProcedimiento      *string
			pAdjudicado                  *string
			pNombreDelProveedor          *string
			pValorTotalAdjudicacion      *float64
			pURLProceso                  *string
			pCodigoPrincipalDeCategoria  *string
			pTipoDeContrato              *string
			pEstadoDeAperturaDelProceso  *string
			pEstadoResumen               *string
			pProveedoresInvitados        *string
			pProveedoresQueManifestaron  *string
			pRespuestasAlProcedimiento   *string
			pNumeroDeLotes               *int64
			pCreatedAt                   *time.Time
			pUpdatedAt                   *time.Time
			pDeletedAt                   *time.Time
		)

		err := rows.Scan(
			&item.ID, &item.PublicCallID, &item.UserID,
			&item.AssociationDate, &item.CreatedAt, &item.UpdatedAt, &assocDeletedAt,
			&pID,
			&pNombreDelProcedimiento, &pEntidad, &pNitEntidad,
			&pDepartamentoEntidad, &pCiudadEntidad, &pOrdenEntidad,
			&pReferenciaDelProceso, &pDescripcionDelProcedimiento,
			&pFase,
			&pFechaDePublicacionDel, &pFechaDeUltimaPublicaci,
			&pModalidadDeContratacion, &pPrecioBase, &pDuracion,
			&pUnidadDeDuracion, &pFechaDeRecepcionDe,
			&pEstadoDelProcedimiento, &pAdjudicado, &pNombreDelProveedor,
			&pValorTotalAdjudicacion, &pURLProceso,
			&pCodigoPrincipalDeCategoria, &pTipoDeContrato,
			&pEstadoDeAperturaDelProceso, &pEstadoResumen,
			&pProveedoresInvitados, &pProveedoresQueManifestaron,
			&pRespuestasAlProcedimiento, &pNumeroDeLotes,
			&pCreatedAt, &pUpdatedAt, &pDeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan saved proposal with proposal: %w", err)
		}

		if pID != nil {
			item.ProposalNombreDelProcedimiento = nullStr(pNombreDelProcedimiento)
			item.ProposalEntidad = nullStr(pEntidad)
			item.ProposalNitEntidad = nullStr(pNitEntidad)
			item.ProposalDepartamentoEntidad = nullStr(pDepartamentoEntidad)
			item.ProposalCiudadEntidad = nullStr(pCiudadEntidad)
			item.ProposalOrdenEntidad = nullStr(pOrdenEntidad)
			item.ProposalReferenciaDelProceso = nullStr(pReferenciaDelProceso)
			item.ProposalDescripcionDelProcedimiento = nullStr(pDescripcionDelProcedimiento)
			item.ProposalFase = nullStr(pFase)
			item.ProposalFechaDePublicacionDel = nullTimeStr(pFechaDePublicacionDel)
			item.ProposalFechaDeUltimaPublicaci = nullTimeStr(pFechaDeUltimaPublicaci)
			item.ProposalModalidadDeContratacion = nullStr(pModalidadDeContratacion)
			item.ProposalPrecioBase = ptrToFloat64(pPrecioBase)
			item.ProposalDuracion = nullStr(pDuracion)
			item.ProposalUnidadDeDuracion = nullStr(pUnidadDeDuracion)
			item.ProposalFechaDeRecepcionDe = nullTimeStr(pFechaDeRecepcionDe)
			item.ProposalEstadoDelProcedimiento = nullStr(pEstadoDelProcedimiento)
			item.ProposalAdjudicado = nullStr(pAdjudicado)
			item.ProposalNombreDelProveedor = nullStr(pNombreDelProveedor)
			item.ProposalValorTotalAdjudicacion = ptrToFloat64(pValorTotalAdjudicacion)
			item.ProposalURLProceso = nullStr(pURLProceso)
			item.ProposalCodigoPrincipalDeCategoria = nullStr(pCodigoPrincipalDeCategoria)
			item.ProposalTipoDeContrato = nullStr(pTipoDeContrato)
			item.ProposalEstadoDeAperturaDelProceso = nullStr(pEstadoDeAperturaDelProceso)
			item.ProposalEstadoResumen = nullStr(pEstadoResumen)
			item.ProposalProveedoresInvitados = nullStr(pProveedoresInvitados)
			item.ProposalProveedoresQueManifestaron = nullStr(pProveedoresQueManifestaron)
			item.ProposalRespuestasAlProcedimiento = nullStr(pRespuestasAlProcedimiento)
			item.ProposalNumeroDeLotes = ptrToInt64(pNumeroDeLotes)
		}

		if assocDeletedAt != nil {
			item.DeletedAt = assocDeletedAt
		}

		results = append(results, item)
	}

	return results, rows.Err()
}

// nullStr converts a *string to string, returning empty string for nil.
func nullStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// nullTimeStr converts a *time.Time to its RFC3339 string representation.
func nullTimeStr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

// ptrToFloat64 converts a *float64 to float64, returning 0 for nil.
func ptrToFloat64(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

// ptrToInt64 converts a *int64 to int64, returning 0 for nil.
func ptrToInt64(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

// Delete removes a saved proposal association.
func (r *PostgresSavedProposalRepository) Delete(ctx context.Context, userID, publicCallID string) error {
	pcID, err := strconv.Atoi(publicCallID)
	if err != nil {
		return fmt.Errorf("invalid public call id: %w", err)
	}

	query := `
		DELETE FROM public_call_user_associations
		WHERE user_id = $1 AND public_call_id = $2
	`
	ct, err := r.pool.Exec(ctx, query, userID, pcID)
	if err != nil {
		return fmt.Errorf("delete saved proposal: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
