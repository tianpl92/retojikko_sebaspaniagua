package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/domain"
)

var (
	ErrDuplicateDocument = errors.New("document already exists")
)

// mockUserRepository is an in-memory implementation of UserRepository for development/testing.
type mockUserRepository struct {
	mu      sync.RWMutex
	users   map[string]*domain.User // keyed by email
	byDoc   map[string]*domain.User // keyed by document (ID)
	counter int
}

// NewMockUserRepository creates a new mock user repository with a default test user.
func NewMockUserRepository() UserRepository {
	repo := &mockUserRepository{
		users: make(map[string]*domain.User),
		byDoc: make(map[string]*domain.User),
	}
	// Pre-seed with a default test user for development
	repo.users["test@example.com"] = &domain.User{
		ID:        "1",
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
		Password:  "$2a$10$IqHCg3FyYUqys/g/InYzouZ5sAwUG/3REhJ5NmWLQVocIsQaVRaQe",
		Status:    "active",
	}
	repo.byDoc["1"] = repo.users["test@example.com"]
	return repo
}

func (r *mockUserRepository) FindByID(_ context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.byDoc[id]
	if !ok {
		return nil, ErrNotFound
	}
	return user, nil
}

func (r *mockUserRepository) FindByEmail(_ context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[email]
	if !ok {
		return nil, ErrNotFound
	}
	return user, nil
}

func (r *mockUserRepository) Create(_ context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.users[user.Email]; exists {
		return ErrDuplicateEmail
	}
	if user.ID != "" {
		if _, exists := r.byDoc[user.ID]; exists {
			return ErrDuplicateDocument
		}
	}
	r.counter++
	if user.ID == "" {
		user.ID = fmt.Sprintf("%d", r.counter)
	}
	r.users[user.Email] = user
	r.byDoc[user.ID] = user
	return nil
}

func (r *mockUserRepository) ExistsByDocument(_ context.Context, document string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.byDoc[document]
	return ok, nil
}

func (r *mockUserRepository) ExistsByEmail(_ context.Context, email string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.users[email]
	return ok, nil
}

func (r *mockUserRepository) Update(_ context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.byDoc[user.ID]
	if !ok {
		return ErrNotFound
	}
	// If email changed, remove old email key
	if existing.Email != user.Email {
		delete(r.users, existing.Email)
	}
	r.users[user.Email] = user
	r.byDoc[user.ID] = user
	return nil
}

// mockProposalRepository is an in-memory implementation of ProposalRepository for development/testing.
type mockProposalRepository struct {
	mu        sync.RWMutex
	proposals []*domain.PublicCallProposal
	byID      map[string]*domain.PublicCallProposal // keyed by fmt.Sprint(ID)
}

// NewMockProposalRepository creates a new mock proposal repository.
func NewMockProposalRepository() ProposalRepository {
	return &mockProposalRepository{
		proposals: []*domain.PublicCallProposal{},
		byID:      make(map[string]*domain.PublicCallProposal),
	}
}

func (r *mockProposalRepository) List(_ context.Context, limit, offset int) ([]*domain.PublicCallProposal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if offset >= len(r.proposals) {
		return []*domain.PublicCallProposal{}, nil
	}
	end := offset + limit
	if end > len(r.proposals) || limit <= 0 {
		end = len(r.proposals)
	}
	if limit <= 0 {
		result := make([]*domain.PublicCallProposal, len(r.proposals))
		copy(result, r.proposals)
		return result, nil
	}
	result := make([]*domain.PublicCallProposal, end-offset)
	copy(result, r.proposals[offset:end])
	return result, nil
}

func (r *mockProposalRepository) ListWithFilters(_ context.Context, query, fase, entidad string, limit, offset int) ([]*domain.PublicCallProposal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	// Simplified mock: return all proposals (clientside filtering would happen in real impl)
	return r.List(context.TODO(), limit, offset)
}

func (r *mockProposalRepository) FindByID(_ context.Context, id string) (*domain.PublicCallProposal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if p, ok := r.byID[id]; ok {
		return p, nil
	}
	return nil, ErrNotFound
}

func (r *mockProposalRepository) Upsert(_ context.Context, proposal *domain.PublicCallProposal) (*domain.PublicCallProposal, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := fmt.Sprintf("%s", proposal.ID)
	if existing, ok := r.byID[key]; ok {
		// Update existing
		existing.NombreDelProcedimiento = proposal.NombreDelProcedimiento
		existing.Entidad = proposal.Entidad
		existing.NitEntidad = proposal.NitEntidad
		existing.DepartamentoEntidad = proposal.DepartamentoEntidad
		existing.CiudadEntidad = proposal.CiudadEntidad
		existing.OrdenEntidad = proposal.OrdenEntidad
		existing.ReferenciaDelProceso = proposal.ReferenciaDelProceso
		existing.DescripcionDelProcedimiento = proposal.DescripcionDelProcedimiento
		existing.Fase = proposal.Fase
		existing.FechaDePublicacionDel = proposal.FechaDePublicacionDel
		existing.FechaDeUltimaPublicaci = proposal.FechaDeUltimaPublicaci
		existing.ModalidadDeContratacion = proposal.ModalidadDeContratacion
		existing.PrecioBase = proposal.PrecioBase
		existing.Duracion = proposal.Duracion
		existing.UnidadDeDuracion = proposal.UnidadDeDuracion
		existing.FechaDeRecepcionDe = proposal.FechaDeRecepcionDe
		existing.EstadoDelProcedimiento = proposal.EstadoDelProcedimiento
		existing.Adjudicado = proposal.Adjudicado
		existing.NombreDelProveedor = proposal.NombreDelProveedor
		existing.ValorTotalAdjudicacion = proposal.ValorTotalAdjudicacion
		existing.URLProceso = proposal.URLProceso
		existing.CodigoPrincipalDeCategoria = proposal.CodigoPrincipalDeCategoria
		existing.TipoDeContrato = proposal.TipoDeContrato
		existing.EstadoDeAperturaDelProceso = proposal.EstadoDeAperturaDelProceso
		existing.EstadoResumen = proposal.EstadoResumen
		existing.ProveedoresInvitados = proposal.ProveedoresInvitados
		existing.ProveedoresQueManifestaron = proposal.ProveedoresQueManifestaron
		existing.RespuestasAlProcedimiento = proposal.RespuestasAlProcedimiento
		existing.NumeroDeLotes = proposal.NumeroDeLotes
		existing.UpdatedAt = time.Now()
		return existing, nil
	}
	// Insert new
	now := time.Now()
	proposal.CreatedAt = now
	proposal.UpdatedAt = now
	if proposal.ID == "" {
		proposal.ID = fmt.Sprintf("%d", len(r.proposals)+1)
	}
	key = proposal.ID
	r.proposals = append(r.proposals, proposal)
	r.byID[key] = proposal
	return proposal, nil
}

// mockSavedProposalRepository is an in-memory implementation of SavedProposalRepository.
type mockSavedProposalRepository struct {
	mu            sync.RWMutex
	saved         map[string]*domain.SavedProposal // keyed by ID
	byUser        map[string][]string              // userID -> list of saved IDs
	byUserAndCall map[string]string                // "userID:callID" -> saved ID
	proposalRepo  ProposalRepository               // for full proposal lookups
	counter       int
}

// NewMockSavedProposalRepository creates a new mock saved proposal repository.
// If proposalRepo is non-nil, FindByUserIDWithProposals will populate proposal data.
func NewMockSavedProposalRepository(proposalRepo ProposalRepository) SavedProposalRepository {
	return &mockSavedProposalRepository{
		saved:         make(map[string]*domain.SavedProposal),
		byUser:        make(map[string][]string),
		byUserAndCall: make(map[string]string),
		proposalRepo:  proposalRepo,
	}
}

func (r *mockSavedProposalRepository) Save(_ context.Context, userID, publicCallID string) (*domain.SavedProposal, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check for existing association (idempotent)
	key := userID + ":" + publicCallID
	if existingID, ok := r.byUserAndCall[key]; ok {
		if s, ok := r.saved[existingID]; ok {
			return s, nil
		}
	}

	r.counter++
	id := fmt.Sprintf("mock-uuid-%d", r.counter)
	saved := &domain.SavedProposal{
		ID:              id,
		PublicCallID:    publicCallID,
		UserID:          userID,
		AssociationDate: time.Now(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	r.saved[id] = saved
	r.byUser[userID] = append(r.byUser[userID], id)
	r.byUserAndCall[key] = id
	return saved, nil
}

func (r *mockSavedProposalRepository) FindByUserID(_ context.Context, userID string) ([]*domain.SavedProposal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := r.byUser[userID]
	result := make([]*domain.SavedProposal, 0, len(ids))
	for _, id := range ids {
		if s, ok := r.saved[id]; ok {
			result = append(result, s)
		}
	}
	return result, nil
}

func (r *mockSavedProposalRepository) FindByUserIDWithProposals(ctx context.Context, userID string) ([]*domain.SavedProposalWithProposal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := r.byUser[userID]
	proposals := make([]*domain.SavedProposal, 0, len(ids))
	for _, id := range ids {
		if s, ok := r.saved[id]; ok {
			proposals = append(proposals, s)
		}
	}

	result := make([]*domain.SavedProposalWithProposal, 0, len(proposals))
	for _, sp := range proposals {
		item := &domain.SavedProposalWithProposal{
			ID:              sp.ID,
			PublicCallID:    sp.PublicCallID,
			UserID:          sp.UserID,
			AssociationDate: sp.AssociationDate,
			CreatedAt:       sp.CreatedAt,
			UpdatedAt:       sp.UpdatedAt,
			DeletedAt:       sp.DeletedAt,
		}

		// Try to look up proposal data
		if r.proposalRepo != nil {
			callID := sp.PublicCallID
			proposal, lookupErr := r.proposalRepo.FindByID(ctx, callID)
			if lookupErr == nil && proposal != nil {
				if proposal.NombreDelProcedimiento.Valid {
					item.ProposalNombreDelProcedimiento = proposal.NombreDelProcedimiento.String
				}
				if proposal.Entidad.Valid {
					item.ProposalEntidad = proposal.Entidad.String
				}
				if proposal.NitEntidad.Valid {
					item.ProposalNitEntidad = proposal.NitEntidad.String
				}
				if proposal.DepartamentoEntidad.Valid {
					item.ProposalDepartamentoEntidad = proposal.DepartamentoEntidad.String
				}
				if proposal.CiudadEntidad.Valid {
					item.ProposalCiudadEntidad = proposal.CiudadEntidad.String
				}
				if proposal.OrdenEntidad.Valid {
					item.ProposalOrdenEntidad = proposal.OrdenEntidad.String
				}
				if proposal.ReferenciaDelProceso.Valid {
					item.ProposalReferenciaDelProceso = proposal.ReferenciaDelProceso.String
				}
				if proposal.DescripcionDelProcedimiento.Valid {
					item.ProposalDescripcionDelProcedimiento = proposal.DescripcionDelProcedimiento.String
				}
				if proposal.Fase.Valid {
					item.ProposalFase = proposal.Fase.String
				}
				if proposal.FechaDePublicacionDel.Valid {
					item.ProposalFechaDePublicacionDel = proposal.FechaDePublicacionDel.Time.Format(time.RFC3339)
				}
				if proposal.FechaDeUltimaPublicaci.Valid {
					item.ProposalFechaDeUltimaPublicaci = proposal.FechaDeUltimaPublicaci.Time.Format(time.RFC3339)
				}
				if proposal.ModalidadDeContratacion.Valid {
					item.ProposalModalidadDeContratacion = proposal.ModalidadDeContratacion.String
				}
				if proposal.PrecioBase.Valid {
					item.ProposalPrecioBase = proposal.PrecioBase.Float64
				}
				if proposal.Duracion.Valid {
					item.ProposalDuracion = proposal.Duracion.String
				}
				if proposal.UnidadDeDuracion.Valid {
					item.ProposalUnidadDeDuracion = proposal.UnidadDeDuracion.String
				}
				if proposal.FechaDeRecepcionDe.Valid {
					item.ProposalFechaDeRecepcionDe = proposal.FechaDeRecepcionDe.Time.Format(time.RFC3339)
				}
				if proposal.EstadoDelProcedimiento.Valid {
					item.ProposalEstadoDelProcedimiento = proposal.EstadoDelProcedimiento.String
				}
				if proposal.Adjudicado.Valid {
					item.ProposalAdjudicado = proposal.Adjudicado.String
				}
				if proposal.NombreDelProveedor.Valid {
					item.ProposalNombreDelProveedor = proposal.NombreDelProveedor.String
				}
				if proposal.ValorTotalAdjudicacion.Valid {
					item.ProposalValorTotalAdjudicacion = proposal.ValorTotalAdjudicacion.Float64
				}
				if proposal.URLProceso.Valid {
					item.ProposalURLProceso = proposal.URLProceso.String
				}
				if proposal.CodigoPrincipalDeCategoria.Valid {
					item.ProposalCodigoPrincipalDeCategoria = proposal.CodigoPrincipalDeCategoria.String
				}
				if proposal.TipoDeContrato.Valid {
					item.ProposalTipoDeContrato = proposal.TipoDeContrato.String
				}
				if proposal.EstadoDeAperturaDelProceso.Valid {
					item.ProposalEstadoDeAperturaDelProceso = proposal.EstadoDeAperturaDelProceso.String
				}
				if proposal.EstadoResumen.Valid {
					item.ProposalEstadoResumen = proposal.EstadoResumen.String
				}
				if proposal.ProveedoresInvitados.Valid {
					item.ProposalProveedoresInvitados = proposal.ProveedoresInvitados.String
				}
				if proposal.ProveedoresQueManifestaron.Valid {
					item.ProposalProveedoresQueManifestaron = proposal.ProveedoresQueManifestaron.String
				}
				if proposal.RespuestasAlProcedimiento.Valid {
					item.ProposalRespuestasAlProcedimiento = proposal.RespuestasAlProcedimiento.String
				}
				if proposal.NumeroDeLotes.Valid {
					item.ProposalNumeroDeLotes = proposal.NumeroDeLotes.Int64
				}
			}
		}

		result = append(result, item)
	}

	return result, nil
}

func (r *mockSavedProposalRepository) Delete(_ context.Context, userID, publicCallID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, s := range r.saved {
		if s.UserID == userID && s.PublicCallID == publicCallID {
			delete(r.saved, id)
			ids := r.byUser[userID]
			for i, sid := range ids {
				if sid == id {
					r.byUser[userID] = append(ids[:i], ids[i+1:]...)
					break
				}
			}
			key := userID + ":" + publicCallID
			delete(r.byUserAndCall, key)
			return nil
		}
	}
	return ErrNotFound
}

func parseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}

func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func nullFloat64(f float64) sql.NullFloat64 {
	return sql.NullFloat64{Float64: f, Valid: true}
}

func nullInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{Int64: i, Valid: true}
}

func nullTime(s string) sql.NullTime {
	if s == "" {
		return sql.NullTime{Valid: false}
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: t, Valid: true}
}

// mockSessionRepository is an in-memory implementation of SessionRepository for testing.
type mockSessionRepository struct {
	mu       sync.RWMutex
	sessions map[string]*domain.UserSession // keyed by token
	counter  int
}

// NewMockSessionRepository creates a new mock session repository.
func NewMockSessionRepository() SessionRepository {
	return &mockSessionRepository{
		sessions: make(map[string]*domain.UserSession),
	}
}

func (r *mockSessionRepository) Create(_ context.Context, userID, token string, expiresAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counter++
	session := &domain.UserSession{
		ID:        fmt.Sprintf("session-%d", r.counter),
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
	r.sessions[token] = session
	return nil
}

func (r *mockSessionRepository) FindByToken(_ context.Context, token string) (*domain.UserSession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	session, ok := r.sessions[token]
	if !ok {
		return nil, ErrNotFound
	}
	return session, nil
}

func (r *mockSessionRepository) DeleteExpired(_ context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	for token, session := range r.sessions {
		if session.ExpiresAt.Before(now) {
			delete(r.sessions, token)
		}
	}
	return nil
}
