package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/domain"
)

// PostgresProposalRepository is a PostgreSQL-backed implementation of ProposalRepository.
type PostgresProposalRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresProposalRepository creates a new PostgresProposalRepository.
func NewPostgresProposalRepository(pool *pgxpool.Pool) *PostgresProposalRepository {
	return &PostgresProposalRepository{pool: pool}
}

// List retrieves proposals with pagination.
func (r *PostgresProposalRepository) List(ctx context.Context, limit, offset int) ([]*domain.PublicCallProposal, error) {
	query := `
		SELECT id, nombre_del_procedimiento, entidad, nit_entidad,
		       departamento_entidad, ciudad_entidad, ordenentidad,
		       referencia_del_proceso, descripci_n_del_procedimiento,
		       fase, fecha_de_publicacion_del, fecha_de_ultima_publicaci,
		       modalidad_de_contratacion, precio_base, duracion,
		       unidad_de_duracion, fecha_de_recepcion_de,
		       estado_del_procedimiento, adjudicado, nombre_del_proveedor,
		       valor_total_adjudicacion, urlproceso,
		       codigo_principal_de_categoria, tipo_de_contrato,
		       estado_de_apertura_del_proceso, estado_resumen,
		       proveedores_invitados, proveedores_que_manifestaron,
		       respuestas_al_procedimiento, numero_de_lotes,
		       created_at, updated_at, deleted_at
		FROM public_calls_proposals
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list proposals: %w", err)
	}
	defer rows.Close()
	return scanProposals(rows)
}

// ListWithFilters retrieves proposals with search/filter criteria and pagination.
func (r *PostgresProposalRepository) ListWithFilters(ctx context.Context, query, fase, entidad string, limit, offset int) ([]*domain.PublicCallProposal, error) {
	base := `
		SELECT id, nombre_del_procedimiento, entidad, nit_entidad,
		       departamento_entidad, ciudad_entidad, ordenentidad,
		       referencia_del_proceso, descripci_n_del_procedimiento,
		       fase, fecha_de_publicacion_del, fecha_de_ultima_publicaci,
		       modalidad_de_contratacion, precio_base, duracion,
		       unidad_de_duracion, fecha_de_recepcion_de,
		       estado_del_procedimiento, adjudicado, nombre_del_proveedor,
		       valor_total_adjudicacion, urlproceso,
		       codigo_principal_de_categoria, tipo_de_contrato,
		       estado_de_apertura_del_proceso, estado_resumen,
		       proveedores_invitados, proveedores_que_manifestaron,
		       respuestas_al_procedimiento, numero_de_lotes,
		       created_at, updated_at, deleted_at
		FROM public_calls_proposals
		WHERE deleted_at IS NULL
	`
	args := []interface{}{}
	argIdx := 1
	conditions := []string{}

	if query != "" {
		conditions = append(conditions, fmt.Sprintf(
			"(nombre_del_procedimiento ILIKE $%d OR referencia_del_proceso ILIKE $%d OR descripci_n_del_procedimiento ILIKE $%d)",
			argIdx, argIdx, argIdx,
		))
		args = append(args, "%"+query+"%")
		argIdx++
	}
	if fase != "" {
		conditions = append(conditions, fmt.Sprintf("fase = $%d", argIdx))
		args = append(args, fase)
		argIdx++
	}
	if entidad != "" {
		conditions = append(conditions, fmt.Sprintf("entidad ILIKE $%d", argIdx))
		args = append(args, "%"+entidad+"%")
		argIdx++
	}

	if len(conditions) > 0 {
		base += " AND " + strings.Join(conditions, " AND ")
	}

	base += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, base, args...)
	if err != nil {
		return nil, fmt.Errorf("list proposals with filters: %w", err)
	}
	defer rows.Close()
	return scanProposals(rows)
}

// FindByID retrieves a single proposal by its numeric ID.
func (r *PostgresProposalRepository) FindByID(ctx context.Context, id string) (*domain.PublicCallProposal, error) {
	query := `
		SELECT id, nombre_del_procedimiento, entidad, nit_entidad,
		       departamento_entidad, ciudad_entidad, ordenentidad,
		       referencia_del_proceso, descripci_n_del_procedimiento,
		       fase, fecha_de_publicacion_del, fecha_de_ultima_publicaci,
		       modalidad_de_contratacion, precio_base, duracion,
		       unidad_de_duracion, fecha_de_recepcion_de,
		       estado_del_procedimiento, adjudicado, nombre_del_proveedor,
		       valor_total_adjudicacion, urlproceso,
		       codigo_principal_de_categoria, tipo_de_contrato,
		       estado_de_apertura_del_proceso, estado_resumen,
		       proveedores_invitados, proveedores_que_manifestaron,
		       respuestas_al_procedimiento, numero_de_lotes,
		       created_at, updated_at, deleted_at
		FROM public_calls_proposals
		WHERE id = $1 AND deleted_at IS NULL
	`
	row := r.pool.QueryRow(ctx, query, id)
	proposal, err := scanProposal(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find proposal by id: %w", err)
	}
	return proposal, nil
}

// scanProposals scans all rows from a pgx.Rows into a slice of PublicCallProposal.
func scanProposals(rows pgx.Rows) ([]*domain.PublicCallProposal, error) {
	var results []*domain.PublicCallProposal
	for rows.Next() {
		p, err := scanProposalFromRows(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, p)
	}
	return results, rows.Err()
}

// scanProposal scans a single pgx.Row into a PublicCallProposal.
func scanProposal(row pgx.Row) (*domain.PublicCallProposal, error) {
	p := &domain.PublicCallProposal{}
	err := row.Scan(
		&p.ID, &p.NombreDelProcedimiento, &p.Entidad, &p.NitEntidad,
		&p.DepartamentoEntidad, &p.CiudadEntidad, &p.OrdenEntidad,
		&p.ReferenciaDelProceso, &p.DescripcionDelProcedimiento,
		&p.Fase, &p.FechaDePublicacionDel, &p.FechaDeUltimaPublicaci,
		&p.ModalidadDeContratacion, &p.PrecioBase, &p.Duracion,
		&p.UnidadDeDuracion, &p.FechaDeRecepcionDe,
		&p.EstadoDelProcedimiento, &p.Adjudicado, &p.NombreDelProveedor,
		&p.ValorTotalAdjudicacion, &p.URLProceso,
		&p.CodigoPrincipalDeCategoria, &p.TipoDeContrato,
		&p.EstadoDeAperturaDelProceso, &p.EstadoResumen,
		&p.ProveedoresInvitados, &p.ProveedoresQueManifestaron,
		&p.RespuestasAlProcedimiento, &p.NumeroDeLotes,
		&p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// scanProposalFromRows scans the current row from pgx.Rows into a PublicCallProposal.
func scanProposalFromRows(rows pgx.Rows) (*domain.PublicCallProposal, error) {
	p := &domain.PublicCallProposal{}
	err := rows.Scan(
		&p.ID, &p.NombreDelProcedimiento, &p.Entidad, &p.NitEntidad,
		&p.DepartamentoEntidad, &p.CiudadEntidad, &p.OrdenEntidad,
		&p.ReferenciaDelProceso, &p.DescripcionDelProcedimiento,
		&p.Fase, &p.FechaDePublicacionDel, &p.FechaDeUltimaPublicaci,
		&p.ModalidadDeContratacion, &p.PrecioBase, &p.Duracion,
		&p.UnidadDeDuracion, &p.FechaDeRecepcionDe,
		&p.EstadoDelProcedimiento, &p.Adjudicado, &p.NombreDelProveedor,
		&p.ValorTotalAdjudicacion, &p.URLProceso,
		&p.CodigoPrincipalDeCategoria, &p.TipoDeContrato,
		&p.EstadoDeAperturaDelProceso, &p.EstadoResumen,
		&p.ProveedoresInvitados, &p.ProveedoresQueManifestaron,
		&p.RespuestasAlProcedimiento, &p.NumeroDeLotes,
		&p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}
