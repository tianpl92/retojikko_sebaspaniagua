-- psql connection command
-- psql -d portal_plan_public_app

CREATE DATABASE portal_plan_public_app;

-- Connect to the database
-- \c portal_plan_public_app

-- Table: users
CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    gender VARCHAR(50),
    phone_number VARCHAR(50),
    password VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);

-- Table: public_calls_proposals
CREATE TABLE public_calls_proposals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nombre_del_procedimiento VARCHAR(500) NOT NULL,
    entidad VARCHAR(255) NOT NULL,
    nit_entidad VARCHAR(50),
    departamento_entidad VARCHAR(100),
    ciudad_entidad VARCHAR(100),
    ordenentidad VARCHAR(50),
    referencia_del_proceso VARCHAR(100),
    descripci_n_del_procedimiento TEXT,
    fase VARCHAR(100) NOT NULL,
    fecha_de_publicacion_del TIMESTAMP,
    fecha_de_ultima_publicaci TIMESTAMP,
    modalidad_de_contratacion VARCHAR(100),
    precio_base NUMERIC(15, 2),
    duracion INTEGER,
    unidad_de_duracion VARCHAR(50),
    fecha_de_recepcion_de TIMESTAMP,
    estado_del_procedimiento VARCHAR(100),
    adjudicado VARCHAR(10),
    nombre_del_proveedor VARCHAR(255),
    valor_total_adjudicacion NUMERIC(15, 2),
    urlproceso TEXT,
    codigo_principal_de_categoria VARCHAR(50),
    tipo_de_contrato VARCHAR(100),
    estado_de_apertura_del_proceso VARCHAR(100),
    estado_resumen VARCHAR(100),
    proveedores_invitados INTEGER,
    proveedores_que_manifestaron INTEGER,
    respuestas_al_procedimiento INTEGER,
    numero_de_lotes INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: user_call_associations
CREATE TABLE user_call_associations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    public_call_id UUID NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    association_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Foreign Keys
ALTER TABLE user_call_associations
    ADD CONSTRAINT fk_associations_public_call
    FOREIGN KEY (public_call_id) REFERENCES public_calls_proposals(id) ON DELETE RESTRICT;

ALTER TABLE user_call_associations
    ADD CONSTRAINT fk_associations_user
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT;

-- Indexes
CREATE INDEX idx_associations_public_call ON user_call_associations(public_call_id);
CREATE INDEX idx_associations_user ON user_call_associations(user_id);
CREATE UNIQUE INDEX idx_associations_unique ON user_call_associations(public_call_id, user_id);

-- Comments
COMMENT ON TABLE users IS 'Stores portal login users and their authentication/contact details.';
COMMENT ON COLUMN users.id IS 'Primary identifier for a portal user.';
COMMENT ON COLUMN users.first_name IS 'User first name.';
COMMENT ON COLUMN users.last_name IS 'User last name.';
COMMENT ON COLUMN users.gender IS 'User gender, when provided.';
COMMENT ON COLUMN users.email IS 'User email address used for login and contact lookup.';
COMMENT ON COLUMN users.phone_number IS 'User phone number, when provided.';
COMMENT ON COLUMN users.password IS 'User password or password reference used for authentication.';

COMMENT ON TABLE public_calls_proposals IS 'Stores centralized public calls for proposals from SECOP II with entity details, phases, and award data.';
COMMENT ON COLUMN public_calls_proposals.id IS 'Unique process ID from SECOP II platform, Primary identifier.';
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

COMMENT ON TABLE user_call_associations IS 'Links registered users to public calls for proposals.';
COMMENT ON COLUMN user_call_associations.id IS 'Primary identifier for the user-to-call association record.';
COMMENT ON COLUMN user_call_associations.public_call_id IS 'Foreign key referencing the associated public call for proposals.';
COMMENT ON COLUMN user_call_associations.user_id IS 'Foreign key referencing the registered portal user.';
COMMENT ON COLUMN user_call_associations.association_date IS 'Timestamp when the association was created.';
