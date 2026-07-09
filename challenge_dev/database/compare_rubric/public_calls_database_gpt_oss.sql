CREATE DATABASE portal_plan_public_app;
\connect portal_plan_public_app;

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id            VARCHAR(255) PRIMARY KEY,
    first_name    VARCHAR(255) NOT NULL,
    last_name     VARCHAR(255) NOT NULL,
    gender        VARCHAR(50),
    email         VARCHAR(255) NOT NULL,
    phone_number  VARCHAR(50),
    password      VARCHAR(255),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ
);

CREATE TABLE public_calls_proposals (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    call_name     VARCHAR(255) NOT NULL,
    call_code     VARCHAR(255),
    call_date     TIMESTAMPTZ,
    issuer_name   VARCHAR(255),
    issuer_id     VARCHAR(255),
    issuer_phone  VARCHAR(50),
    issuer_email  VARCHAR(255),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ
);

CREATE TABLE public_call_user_association (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    public_call_id    UUID NOT NULL,
    user_id           VARCHAR(255) NOT NULL,
    association_date  TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at        TIMESTAMPTZ,
    CONSTRAINT fk_public_call
        FOREIGN KEY (public_call_id)
        REFERENCES public_calls_proposals(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,
    CONSTRAINT uq_public_call_user UNIQUE (public_call_id, user_id)
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_assoc_public_call_id ON public_call_user_association(public_call_id);
CREATE INDEX idx_assoc_user_id ON public_call_user_association(user_id);

COMMENT ON TABLE users IS 'Stores portal login users and their authentication/contact details.';
COMMENT ON COLUMN users.id IS 'Primary identifier for a portal user.';
COMMENT ON COLUMN users.first_name IS 'User first name.';
COMMENT ON COLUMN users.last_name IS 'User last name.';
COMMENT ON COLUMN users.gender IS 'User gender, when provided.';
COMMENT ON COLUMN users.email IS 'User email address used for login and contact lookup.';
COMMENT ON COLUMN users.phone_number IS 'User phone number, when provided.';
COMMENT ON COLUMN users.password IS 'User password or password reference used for authentication.';
COMMENT ON COLUMN users.created_at IS 'Record creation timestamp.';
COMMENT ON COLUMN users.updated_at IS 'Record last update timestamp.';
COMMENT ON COLUMN users.deleted_at IS 'Record deletion timestamp, if soft deleted.';

COMMENT ON TABLE public_calls_proposals IS 'Stores centralized public calls for proposals and issuer metadata.';
COMMENT ON COLUMN public_calls_proposals.id IS 'Primary identifier for a public call for proposals record.';
COMMENT ON COLUMN public_calls_proposals.call_name IS 'Name of the public call for proposals.';
COMMENT ON COLUMN public_calls_proposals.call_code IS 'Reference code or call identifier from the source.';
COMMENT ON COLUMN public_calls_proposals.call_date IS 'Timestamp indicating when the call for proposals was issued or recorded.';
COMMENT ON COLUMN public_calls_proposals.issuer_name IS 'Name of the person or entity issuing the call.';
COMMENT ON COLUMN public_calls_proposals.issuer_id IS 'Identification number of the person or entity issuing the call.';
COMMENT ON COLUMN public_calls_proposals.issuer_phone IS 'Phone number of the issuer.';
COMMENT ON COLUMN public_calls_proposals.issuer_email IS 'Email address of the issuer.';
COMMENT ON COLUMN public_calls_proposals.created_at IS 'Record creation timestamp.';
COMMENT ON COLUMN public_calls_proposals.updated_at IS 'Record last update timestamp.';
COMMENT ON COLUMN public_calls_proposals.deleted_at IS 'Record deletion timestamp, if soft deleted.';

COMMENT ON TABLE public_call_user_association IS 'Links registered users to public calls for proposals.';
COMMENT ON COLUMN public_call_user_association.id IS 'Primary identifier for the user-to-call association record.';
COMMENT ON COLUMN public_call_user_association.public_call_id IS 'Foreign key referencing the associated public call for proposals.';
COMMENT ON COLUMN public_call_user_association.user_id IS 'Foreign key referencing the registered portal user.';
COMMENT ON COLUMN public_call_user_association.association_date IS 'Timestamp when the association was created.';
COMMENT ON COLUMN public_call_user_association.created_at IS 'Record creation timestamp.';
COMMENT ON COLUMN public_call_user_association.updated_at IS 'Record last update timestamp.';
COMMENT ON COLUMN public_call_user_association.deleted_at IS 'Record deletion timestamp, if soft deleted.';
