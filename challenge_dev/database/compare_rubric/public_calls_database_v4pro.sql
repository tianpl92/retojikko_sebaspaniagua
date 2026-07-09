CREATE DATABASE portal_plan_public_app;

\c portal_plan_public_app

CREATE TABLE users (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    gender VARCHAR(50),
    email VARCHAR(255) NOT NULL,
    phone_number VARCHAR(50),
    password VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

COMMENT ON TABLE users IS 'Stores portal login users and their authentication/contact details.';
COMMENT ON COLUMN users.id IS 'Primary identifier for a portal user.';
COMMENT ON COLUMN users.first_name IS 'User first name.';
COMMENT ON COLUMN users.last_name IS 'User last name.';
COMMENT ON COLUMN users.gender IS 'User gender, when provided.';
COMMENT ON COLUMN users.email IS 'User email address used for login and contact lookup.';
COMMENT ON COLUMN users.phone_number IS 'User phone number, when provided.';
COMMENT ON COLUMN users.password IS 'User password or password reference used for authentication.';
COMMENT ON COLUMN users.created_at IS 'Timestamp when the user record was created.';
COMMENT ON COLUMN users.updated_at IS 'Timestamp when the user record was last updated.';
COMMENT ON COLUMN users.deleted_at IS 'Timestamp when the user record was soft-deleted.';

CREATE INDEX idx_users_email ON users (email);

CREATE TABLE public_calls_proposals (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    id_del_proceso TEXT NOT NULL,
    nombre_del_procedimiento TEXT NOT NULL,
    entidad TEXT NOT NULL,
    nit_entidad TEXT,
    departamento_entidad TEXT,
    ciudad_entidad TEXT,
    ordenentidad TEXT,
    referencia_del_proceso TEXT,
    descripci_n_del_procedimiento TEXT,
    fase TEXT NOT NULL,
    fecha_de_publicacion_del TIMESTAMPTZ,
    fecha_de_ultima_publicaci TIMESTAMPTZ,
    modalidad_de_contratacion TEXT,
    precio_base NUMERIC,
    duracion INTEGER,
    unidad_de_duracion TEXT,
    fecha_de_recepcion_de TIMESTAMPTZ,
    estado_del_procedimiento TEXT,
    adjudicado TEXT,
    nombre_del_proveedor TEXT,
    valor_total_adjudicacion NUMERIC,
    urlproceso TEXT,
    codigo_principal_de_categoria TEXT,
    tipo_de_contrato TEXT,
    estado_de_apertura_del_proceso TEXT,
    estado_resumen TEXT,
    proveedores_invitados INTEGER,
    proveedores_que_manifestaron INTEGER,
    respuestas_al_procedimiento INTEGER,
    numero_de_lotes INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

COMMENT ON TABLE public_calls_proposals IS 'Stores centralized public calls for proposals from SECOP II with entity details, phases, and award data.';
COMMENT ON COLUMN public_calls_proposals.id IS 'Internal UUID primary key for the public call record.';
COMMENT ON COLUMN public_calls_proposals.id_del_proceso IS 'Process ID from SECOP II platform.';
COMMENT ON COLUMN public_calls_proposals.nombre_del_procedimiento IS 'Name of the procurement procedure.';
COMMENT ON COLUMN public_calls_proposals.entidad IS 'Entity publishing the procurement process.';
COMMENT ON COLUMN public_calls_proposals.nit_entidad IS 'Tax ID of the publishing entity.';
COMMENT ON COLUMN public_calls_proposals.departamento_entidad IS 'Department of the entity.';
COMMENT ON COLUMN public_calls_proposals.ciudad_entidad IS 'City of the entity.';
COMMENT ON COLUMN public_calls_proposals.ordenentidad IS 'Entity order (Nacional, Regional).';
COMMENT ON COLUMN public_calls_proposals.referencia_del_proceso IS 'Process reference from entity.';
COMMENT ON COLUMN public_calls_proposals.descripci_n_del_procedimiento IS 'Description of the procedure.';
COMMENT ON COLUMN public_calls_proposals.fase IS 'Current phase of the process.';
COMMENT ON COLUMN public_calls_proposals.fecha_de_publicacion_del IS 'Initial publication date.';
COMMENT ON COLUMN public_calls_proposals.fecha_de_ultima_publicaci IS 'Last publication date.';
COMMENT ON COLUMN public_calls_proposals.modalidad_de_contratacion IS 'Selection modality.';
COMMENT ON COLUMN public_calls_proposals.precio_base IS 'Base estimated price.';
COMMENT ON COLUMN public_calls_proposals.duracion IS 'Estimated duration.';
COMMENT ON COLUMN public_calls_proposals.unidad_de_duracion IS 'Duration unit.';
COMMENT ON COLUMN public_calls_proposals.fecha_de_recepcion_de IS 'Response submission deadline.';
COMMENT ON COLUMN public_calls_proposals.estado_del_procedimiento IS 'Current status.';
COMMENT ON COLUMN public_calls_proposals.adjudicado IS 'Whether awarded (Si/No).';
COMMENT ON COLUMN public_calls_proposals.nombre_del_proveedor IS 'Awarded supplier name.';
COMMENT ON COLUMN public_calls_proposals.valor_total_adjudicacion IS 'Total award value.';
COMMENT ON COLUMN public_calls_proposals.urlproceso IS 'URL to the process on SECOP II.';
COMMENT ON COLUMN public_calls_proposals.codigo_principal_de_categoria IS 'UNSPSC main category code.';
COMMENT ON COLUMN public_calls_proposals.tipo_de_contrato IS 'Contract type.';
COMMENT ON COLUMN public_calls_proposals.estado_de_apertura_del_proceso IS 'Information opening status.';
COMMENT ON COLUMN public_calls_proposals.estado_resumen IS 'Summary status.';
COMMENT ON COLUMN public_calls_proposals.proveedores_invitados IS 'Total invited suppliers.';
COMMENT ON COLUMN public_calls_proposals.proveedores_que_manifestaron IS 'Suppliers who expressed interest.';
COMMENT ON COLUMN public_calls_proposals.respuestas_al_procedimiento IS 'Total responses.';
COMMENT ON COLUMN public_calls_proposals.numero_de_lotes IS 'Number of item lots.';
COMMENT ON COLUMN public_calls_proposals.created_at IS 'Timestamp when the record was inserted.';
COMMENT ON COLUMN public_calls_proposals.updated_at IS 'Timestamp when the record was last updated.';
COMMENT ON COLUMN public_calls_proposals.deleted_at IS 'Timestamp when the record was soft-deleted.';

CREATE TABLE user_public_call_association (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    public_call_id UUID NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    association_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_association_public_call FOREIGN KEY (public_call_id) REFERENCES public_calls_proposals(id) ON DELETE CASCADE,
    CONSTRAINT fk_association_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT uq_association_user_call UNIQUE (public_call_id, user_id)
);

COMMENT ON TABLE user_public_call_association IS 'Links registered users to public calls for proposals.';
COMMENT ON COLUMN user_public_call_association.id IS 'Primary identifier for the user-to-call association record.';
COMMENT ON COLUMN user_public_call_association.public_call_id IS 'Foreign key referencing the associated public call for proposals.';
COMMENT ON COLUMN user_public_call_association.user_id IS 'Foreign key referencing the registered portal user.';
COMMENT ON COLUMN user_public_call_association.association_date IS 'Timestamp when the association was created.';
COMMENT ON COLUMN user_public_call_association.updated_at IS 'Timestamp when the association record was last updated.';
COMMENT ON COLUMN user_public_call_association.deleted_at IS 'Timestamp when the association was soft-deleted.';

CREATE INDEX idx_association_public_call_id ON user_public_call_association (public_call_id);
CREATE INDEX idx_association_user_id ON user_public_call_association (user_id);
