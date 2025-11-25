-- +goose Up
-- +goose StatementBegin
CREATE TABLE business (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    cnpj TEXT NOT NULL,
    email TEXT NOT NULL,
    phone TEXT NOT NULL,

    user_id BIGINT NOT NULL REFERENCES users(id),

    version INT NOT NULL DEFAULT 1,
    deleted BOOLEAN NOT NULL DEFAULT FALSE,

    created_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by BIGINT,
    updated_at TIMESTAMPTZ,

    CONSTRAINT unique_user_business_name_deleted UNIQUE (user_id, name, deleted),
    CONSTRAINT unique_user_business_cnpj_deleted UNIQUE (user_id, cnpj, deleted),
    CONSTRAINT unique_user_business_email_deleted UNIQUE (user_id, email, deleted)
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS business;
-- +goose StatementEnd
