CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE students (
      id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
      user_id       UUID NOT NULL UNIQUE,
      first_name    VARCHAR(100) NOT NULL DEFAULT '',
      last_name     VARCHAR(100) NOT NULL DEFAULT '',
      phone         VARCHAR(20),
      city          VARCHAR(100),
      bio           TEXT,
      avatar_url    TEXT,
      resume_url    TEXT,
      github_url    TEXT,
      linkedin_url  TEXT,
      portfolio_url TEXT,
      created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
      updated_at    TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE educations (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id     UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    university     VARCHAR(255) NOT NULL,
    faculty        VARCHAR(255),
    specialization VARCHAR(255),
    degree         VARCHAR(50),
    start_year     INT,
    end_year       INT,
    gpa            NUMERIC(3,2),
    is_current     BOOLEAN NOT NULL DEFAULT false,
    created_at     TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE experiences (
     id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
     student_id   UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
     company_name VARCHAR(255) NOT NULL,
     position     VARCHAR(255) NOT NULL,
     description  TEXT,
     start_date   DATE,
     end_date     DATE,
     is_current   BOOLEAN NOT NULL DEFAULT false,
     format       VARCHAR(20),
     created_at   TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE projects (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id  UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    url         TEXT,
    techs       TEXT[] NOT NULL DEFAULT '{}',
    start_date  DATE,
    end_date    DATE,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE student_skills (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    skill      VARCHAR(100) NOT NULL,
    level      VARCHAR(20) NOT NULL DEFAULT 'beginner'
);

CREATE TABLE student_languages (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    language   VARCHAR(100) NOT NULL,
    level      VARCHAR(10) NOT NULL
);

CREATE TABLE companies (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id           UUID NOT NULL UNIQUE,
    name              VARCHAR(255) NOT NULL DEFAULT '',
    tagline           TEXT,
    description       TEXT,
    industry          VARCHAR(100),
    size              VARCHAR(50),
    founded_year      INT,
    website           TEXT,
    contact_email     VARCHAR(255),
    city              VARCHAR(100),
    country           VARCHAR(100),
    is_remote_friendly BOOLEAN NOT NULL DEFAULT false,
    logo_url          TEXT,
    created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_students_user_id  ON students(user_id);
CREATE INDEX idx_companies_user_id ON companies(user_id);