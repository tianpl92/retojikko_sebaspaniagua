CREATE DATABASE portal_plan_public_app;

\connect portal_plan_public_app

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id text PRIMARY KEY,
    first_name text NOT NULL,
    last_name text NOT NULL,
    gender text,
    email text NOT NULL,
    phone_number text,
    password text,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz,
    deleted_at timestamptz,
    CONSTRAINT users_email_not_blank CHECK (btrim(email) <> ''),
    CONSTRAINT users_id_not_blank CHECK (btrim(id) <> ''),
    CONSTRAINT users_first_name_not_blank CHECK (btrim(first_name) <> ''),
    CONSTRAINT users_last_name_not_blank CHECK (btrim(last_name) <> '')
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);

COMMENT ON TABLE users IS 'Stores portal login users and their authentication/contact details.';
COMMENT ON COLUMN users.id IS 'Primary identifier for a portal user.';
COMMENT ON COLUMN users.first_name IS 'User first name.';
COMMENT ON COLUMN users.last_name IS 'User last name.';
COMMENT ON COLUMN users.gender IS 'User gender, when provided.';
COMMENT ON COLUMN users.email IS 'User email address used for login and contact lookup.';
COMMENT ON COLUMN users.phone_number IS 'User phone number, when provided.';
COMMENT ON COLUMN users.password IS 'User password or password reference used for authentication.';
COMMENT ON COLUMN users.created_at IS 'Timestamp when the record was created.';
COMMENT ON COLUMN users.updated_at IS 'Timestamp when the record was last updated.';
COMMENT ON COLUMN users.deleted_at IS 'Timestamp when the record was soft-deleted, if applicable.';

CREATE TABLE IF NOT EXISTS public_calls_proposals (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    call_name text NOT NULL,
    call_code text,
    call_date timestamp,
    issuer_name text,
    issuer_id text,
    issuer_phone text,
    issuer_email text,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz,
    deleted_at timestamptz,
    CONSTRAINT public_calls_proposals_call_name_not_blank CHECK (btrim(call_name) <> '')
);

COMMENT ON TABLE public_calls_proposals IS 'Stores centralized public calls for proposals and issuer metadata.';
COMMENT ON COLUMN public_calls_proposals.id IS 'Primary identifier for a public call for proposals record.';
COMMENT ON COLUMN public_calls_proposals.call_name IS 'Name of the public call for proposals.';
COMMENT ON COLUMN public_calls_proposals.call_code IS 'Reference code or call identifier from the source.';
COMMENT ON COLUMN public_calls_proposals.call_date IS 'Timestamp indicating when the call for proposals was issued or recorded.';
COMMENT ON COLUMN public_calls_proposals.issuer_name IS 'Name of the person or entity issuing the call.';
COMMENT ON COLUMN public_calls_proposals.issuer_id IS 'Identification number of the person or entity issuing the call.';
COMMENT ON COLUMN public_calls_proposals.issuer_phone IS 'Phone number of the issuer.';
COMMENT ON COLUMN public_calls_proposals.issuer_email IS 'Email address of the issuer.';
COMMENT ON COLUMN public_calls_proposals.created_at IS 'Timestamp when the record was created.';
COMMENT ON COLUMN public_calls_proposals.updated_at IS 'Timestamp when the record was last updated.';
COMMENT ON COLUMN public_calls_proposals.deleted_at IS 'Timestamp when the record was soft-deleted, if applicable.';

CREATE TABLE IF NOT EXISTS user_public_calls_proposals_association (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    public_call_id uuid NOT NULL,
    user_id text NOT NULL,
    association_date timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz,
    deleted_at timestamptz,
    CONSTRAINT fk_user_public_calls_association_public_call
        FOREIGN KEY (public_call_id)
        REFERENCES public_calls_proposals (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
    CONSTRAINT fk_user_public_calls_association_user
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
    CONSTRAINT user_public_calls_proposals_association_unique_pair UNIQUE (public_call_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_user_public_calls_association_public_call_id
    ON user_public_calls_proposals_association (public_call_id);

CREATE INDEX IF NOT EXISTS idx_user_public_calls_association_user_id
    ON user_public_calls_proposals_association (user_id);

COMMENT ON TABLE user_public_calls_proposals_association IS 'Links registered users to public calls for proposals.';
COMMENT ON COLUMN user_public_calls_proposals_association.id IS 'Primary identifier for the user-to-call association record.';
COMMENT ON COLUMN user_public_calls_proposals_association.public_call_id IS 'Foreign key referencing the associated public call for proposals.';
COMMENT ON COLUMN user_public_calls_proposals_association.user_id IS 'Foreign key referencing the registered portal user.';
COMMENT ON COLUMN user_public_calls_proposals_association.association_date IS 'Timestamp when the association was created.';
COMMENT ON COLUMN user_public_calls_proposals_association.created_at IS 'Timestamp when the record was created.';
COMMENT ON COLUMN user_public_calls_proposals_association.updated_at IS 'Timestamp when the record was last updated.';
COMMENT ON COLUMN user_public_calls_proposals_association.deleted_at IS 'Timestamp when the record was soft-deleted, if applicable.';
