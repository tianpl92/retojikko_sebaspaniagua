CREATE DATABASE portal_plan_public_app;

\connect portal_plan_public_app

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE OR REPLACE FUNCTION public.set_updated_at()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;

CREATE TABLE public.users (
    id varchar(128) PRIMARY KEY,
    first_name varchar(255) NOT NULL,
    last_name varchar(255) NOT NULL,
    gender varchar(100),
    email varchar(320) NOT NULL,
    phone_number varchar(100),
    password text,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);

CREATE INDEX idx_users_email ON public.users (email);

CREATE TRIGGER trg_users_set_updated_at
BEFORE UPDATE ON public.users
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_at();

COMMENT ON TABLE public.users IS 'Stores portal login users and their authentication/contact details.';
COMMENT ON COLUMN public.users.id IS 'Primary identifier for a portal user.';
COMMENT ON COLUMN public.users.first_name IS 'User first name.';
COMMENT ON COLUMN public.users.last_name IS 'User last name.';
COMMENT ON COLUMN public.users.gender IS 'User gender, when provided.';
COMMENT ON COLUMN public.users.email IS 'User email address used for login and contact lookup.';
COMMENT ON COLUMN public.users.phone_number IS 'User phone number, when provided.';
COMMENT ON COLUMN public.users.password IS 'User password or password reference used for authentication.';
COMMENT ON COLUMN public.users.created_at IS 'Timestamp when the user record was created.';
COMMENT ON COLUMN public.users.updated_at IS 'Timestamp when the user record was last updated.';
COMMENT ON COLUMN public.users.deleted_at IS 'Timestamp when the user record was soft deleted, when applicable.';

CREATE TABLE public.public_calls_proposals (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    id_del_proceso text NOT NULL,
    nombre_del_procedimiento text NOT NULL,
    entidad text NOT NULL,
    nit_entidad text,
    departamento_entidad text,
    ciudad_entidad text,
    ordenentidad text,
    referencia_del_proceso text,
    descripci_n_del_procedimiento text,
    fase text NOT NULL,
    fecha_de_publicacion_del timestamp with time zone,
    fecha_de_ultima_publicaci timestamp with time zone,
    modalidad_de_contratacion text,
    precio_base numeric,
    duracion integer,
    unidad_de_duracion text,
    fecha_de_recepcion_de timestamp with time zone,
    estado_del_procedimiento text,
    adjudicado text,
    nombre_del_proveedor text,
    valor_total_adjudicacion numeric,
    urlproceso text,
    codigo_principal_de_categoria text,
    tipo_de_contrato text,
    estado_de_apertura_del_proceso text,
    estado_resumen text,
    proveedores_invitados integer,
    proveedores_que_manifestaron integer,
    respuestas_al_procedimiento integer,
    numero_de_lotes integer,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    CONSTRAINT uq_public_calls_proposals_id_del_proceso UNIQUE (id_del_proceso)
);

CREATE INDEX idx_public_calls_proposals_id_del_proceso ON public.public_calls_proposals (id_del_proceso);
CREATE INDEX idx_public_calls_proposals_fase ON public.public_calls_proposals (fase);
CREATE INDEX idx_public_calls_proposals_entidad ON public.public_calls_proposals (entidad);

CREATE TRIGGER trg_public_calls_proposals_set_updated_at
BEFORE UPDATE ON public.public_calls_proposals
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_at();

COMMENT ON TABLE public.public_calls_proposals IS 'Stores centralized public calls for proposals from SECOP II with entity details, phases, and award data.';
COMMENT ON COLUMN public.public_calls_proposals.id IS 'UUID primary identifier for the public call for proposals record.';
COMMENT ON COLUMN public.public_calls_proposals.id_del_proceso IS 'Unique process ID from SECOP II platform.';
COMMENT ON COLUMN public.public_calls_proposals.nombre_del_procedimiento IS 'Name of the procurement procedure.';
COMMENT ON COLUMN public.public_calls_proposals.entidad IS 'Entity publishing the procurement process.';
COMMENT ON COLUMN public.public_calls_proposals.nit_entidad IS 'Tax ID of the publishing entity.';
COMMENT ON COLUMN public.public_calls_proposals.departamento_entidad IS 'Department of the entity.';
COMMENT ON COLUMN public.public_calls_proposals.ciudad_entidad IS 'City of the entity.';
COMMENT ON COLUMN public.public_calls_proposals.ordenentidad IS 'Entity order (Nacional, Regional).';
COMMENT ON COLUMN public.public_calls_proposals.referencia_del_proceso IS 'Process reference from entity.';
COMMENT ON COLUMN public.public_calls_proposals.descripci_n_del_procedimiento IS 'Description of the procedure.';
COMMENT ON COLUMN public.public_calls_proposals.fase IS 'Current phase of the process.';
COMMENT ON COLUMN public.public_calls_proposals.fecha_de_publicacion_del IS 'Initial publication date.';
COMMENT ON COLUMN public.public_calls_proposals.fecha_de_ultima_publicaci IS 'Last publication date.';
COMMENT ON COLUMN public.public_calls_proposals.modalidad_de_contratacion IS 'Selection modality.';
COMMENT ON COLUMN public.public_calls_proposals.precio_base IS 'Base estimated price.';
COMMENT ON COLUMN public.public_calls_proposals.duracion IS 'Estimated duration.';
COMMENT ON COLUMN public.public_calls_proposals.unidad_de_duracion IS 'Duration unit.';
COMMENT ON COLUMN public.public_calls_proposals.fecha_de_recepcion_de IS 'Response submission deadline.';
COMMENT ON COLUMN public.public_calls_proposals.estado_del_procedimiento IS 'Current status.';
COMMENT ON COLUMN public.public_calls_proposals.adjudicado IS 'Whether awarded (Si/No).';
COMMENT ON COLUMN public.public_calls_proposals.nombre_del_proveedor IS 'Awarded supplier name.';
COMMENT ON COLUMN public.public_calls_proposals.valor_total_adjudicacion IS 'Total award value.';
COMMENT ON COLUMN public.public_calls_proposals.urlproceso IS 'URL to the process on SECOP II.';
COMMENT ON COLUMN public.public_calls_proposals.codigo_principal_de_categoria IS 'UNSPSC main category code.';
COMMENT ON COLUMN public.public_calls_proposals.tipo_de_contrato IS 'Contract type.';
COMMENT ON COLUMN public.public_calls_proposals.estado_de_apertura_del_proceso IS 'Information opening status.';
COMMENT ON COLUMN public.public_calls_proposals.estado_resumen IS 'Summary status.';
COMMENT ON COLUMN public.public_calls_proposals.proveedores_invitados IS 'Total invited suppliers.';
COMMENT ON COLUMN public.public_calls_proposals.proveedores_que_manifestaron IS 'Suppliers who expressed interest.';
COMMENT ON COLUMN public.public_calls_proposals.respuestas_al_procedimiento IS 'Total responses.';
COMMENT ON COLUMN public.public_calls_proposals.numero_de_lotes IS 'Number of item lots.';
COMMENT ON COLUMN public.public_calls_proposals.created_at IS 'Timestamp when the public call record was created.';
COMMENT ON COLUMN public.public_calls_proposals.updated_at IS 'Timestamp when the public call record was last updated.';
COMMENT ON COLUMN public.public_calls_proposals.deleted_at IS 'Timestamp when the public call record was soft deleted, when applicable.';

CREATE TABLE public.public_call_user_associations (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    public_call_id uuid NOT NULL,
    user_id varchar(128) NOT NULL,
    association_date timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    CONSTRAINT fk_public_call_user_associations_public_call
        FOREIGN KEY (public_call_id)
        REFERENCES public.public_calls_proposals (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
    CONSTRAINT fk_public_call_user_associations_user
        FOREIGN KEY (user_id)
        REFERENCES public.users (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);

CREATE INDEX idx_public_call_user_associations_public_call_id ON public.public_call_user_associations (public_call_id);
CREATE INDEX idx_public_call_user_associations_user_id ON public.public_call_user_associations (user_id);

CREATE TRIGGER trg_public_call_user_associations_set_updated_at
BEFORE UPDATE ON public.public_call_user_associations
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_at();

COMMENT ON TABLE public.public_call_user_associations IS 'Links registered users to public calls for proposals.';
COMMENT ON COLUMN public.public_call_user_associations.id IS 'Primary identifier for the user-to-call association record.';
COMMENT ON COLUMN public.public_call_user_associations.public_call_id IS 'Foreign key referencing the associated public call for proposals.';
COMMENT ON COLUMN public.public_call_user_associations.user_id IS 'Foreign key referencing the registered portal user.';
COMMENT ON COLUMN public.public_call_user_associations.association_date IS 'Timestamp when the association was created.';
COMMENT ON COLUMN public.public_call_user_associations.created_at IS 'Timestamp when the association record was created.';
COMMENT ON COLUMN public.public_call_user_associations.updated_at IS 'Timestamp when the association record was last updated.';
COMMENT ON COLUMN public.public_call_user_associations.deleted_at IS 'Timestamp when the association record was soft deleted, when applicable.';
