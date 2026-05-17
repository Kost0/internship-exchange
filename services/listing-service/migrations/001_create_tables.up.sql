CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE listing_status AS ENUM ('draft', 'active', 'closed');
CREATE TYPE listing_format AS ENUM ('office', 'remote', 'hybrid');
CREATE TYPE employment_type AS ENUM ('full_time', 'part_time', 'project');

CREATE TABLE companies (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id           UUID NOT NULL UNIQUE,
    name              VARCHAR(255) NOT NULL DEFAULT '',
    logo_url          TEXT,
    industry          VARCHAR(100),
    city              VARCHAR(100),
    created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE listings (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id      UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    title           VARCHAR(255) NOT NULL,
    description     TEXT NOT NULL DEFAULT '',
    requirements    TEXT NOT NULL DEFAULT '',
    what_we_offer   TEXT NOT NULL DEFAULT '',
    city            VARCHAR(100),
    format          listing_format NOT NULL DEFAULT 'office',
    employment_type employment_type NOT NULL DEFAULT 'full_time',
    salary_from     BIGINT,
    salary_to       BIGINT,
    salary_currency VARCHAR(10) DEFAULT 'RUB',
    deadline        DATE,
    status          listing_status NOT NULL DEFAULT 'draft',
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE listing_skills (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    listing_id  UUID NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    skill       VARCHAR(100) NOT NULL,
    is_required BOOLEAN NOT NULL DEFAULT true
);

CREATE INDEX idx_listings_company_id ON listings(company_id);
CREATE INDEX idx_listings_status     ON listings(status);
CREATE INDEX idx_listings_format     ON listings(format);
CREATE INDEX idx_listing_skills_listing_id ON listing_skills(listing_id);CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
