CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE application_status AS ENUM ('applied', 'reviewing', 'interview', 'accepted', 'rejected');

CREATE TABLE applications (
      id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
      student_id     UUID NOT NULL,
      listing_id     UUID NOT NULL,
      student_email  VARCHAR(255) NOT NULL DEFAULT '',
      company_email  VARCHAR(255) NOT NULL DEFAULT '',
      cover_letter   TEXT NOT NULL DEFAULT '',
      status         application_status NOT NULL DEFAULT 'applied',
      created_at     TIMESTAMP NOT NULL DEFAULT NOW(),
      updated_at     TIMESTAMP NOT NULL DEFAULT NOW(),
      UNIQUE (student_id, listing_id)
);

CREATE TABLE application_events (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    old_status     application_status,
    new_status     application_status NOT NULL,
    comment        TEXT,
    changed_at     TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_applications_student_id ON applications(student_id);
CREATE INDEX idx_applications_listing_id ON applications(listing_id);
CREATE INDEX idx_application_events_application_id ON application_events(application_id);