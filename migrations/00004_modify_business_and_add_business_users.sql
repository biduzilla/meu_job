-- +goose Up
-- +goose StatementBegin

-- 1) Criar nova tabela pivô business_users
CREATE TABLE IF NOT EXISTS business_users (
    business_id BIGINT NOT NULL REFERENCES business(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (business_id, user_id)
);

-- 2) Migrar dados existentes da coluna business.user_id
INSERT INTO business_users (business_id, user_id)
SELECT id, user_id
FROM business
WHERE user_id IS NOT NULL;

-- 3) Remover constraints antigas que dependiam de user_id
ALTER TABLE business DROP CONSTRAINT IF EXISTS unique_user_business_name_deleted;
ALTER TABLE business DROP CONSTRAINT IF EXISTS unique_user_business_cnpj_deleted;
ALTER TABLE business DROP CONSTRAINT IF EXISTS unique_user_business_email_deleted;

-- 4) Remover coluna user_id da tabela business
ALTER TABLE business 
    DROP COLUMN IF EXISTS user_id;

-- 5) Criar novas constraints de unicidade sem user_id
ALTER TABLE business
    ADD CONSTRAINT unique_business_name_deleted UNIQUE (name, deleted);

ALTER TABLE business
    ADD CONSTRAINT unique_business_cnpj_deleted UNIQUE (cnpj, deleted);

ALTER TABLE business
    ADD CONSTRAINT unique_business_email_deleted UNIQUE (email, deleted);

-- 6) Criar índices úteis
CREATE INDEX IF NOT EXISTS idx_business_users_user ON business_users(user_id);
CREATE INDEX IF NOT EXISTS idx_business_users_business ON business_users(business_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- 1) Recriar coluna user_id na tabela business
ALTER TABLE business 
    ADD COLUMN user_id BIGINT REFERENCES users(id);

-- 2) Restaurar dados — pega o "primeiro" usuário vinculado à empresa
UPDATE business b
SET user_id = bu.user_id
FROM business_users bu
WHERE bu.business_id = b.id;

-- 3) Remover novas constraints
ALTER TABLE business DROP CONSTRAINT IF EXISTS unique_business_name_deleted;
ALTER TABLE business DROP CONSTRAINT IF EXISTS unique_business_cnpj_deleted;
ALTER TABLE business DROP CONSTRAINT IF EXISTS unique_business_email_deleted;

-- 4) Restaurar constraints antigas que dependiam de user_id
ALTER TABLE business
    ADD CONSTRAINT unique_user_business_name_deleted UNIQUE (user_id, name, deleted);

ALTER TABLE business
    ADD CONSTRAINT unique_user_business_cnpj_deleted UNIQUE (user_id, cnpj, deleted);

ALTER TABLE business
    ADD CONSTRAINT unique_user_business_email_deleted UNIQUE (user_id, email, deleted);

-- 5) Remover tabela pivô
DROP TABLE IF EXISTS business_users;

-- +goose StatementEnd
