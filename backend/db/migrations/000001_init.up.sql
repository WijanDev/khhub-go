CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX sessions_expires_at_idx ON sessions (expires_at);

CREATE TABLE congregation (
    id SMALLINT PRIMARY KEY DEFAULT 1 CHECK (id = 1),
    name TEXT NOT NULL DEFAULT '',
    number TEXT NOT NULL DEFAULT '',
    midweek_day SMALLINT NOT NULL DEFAULT 4 CHECK (midweek_day BETWEEN 0 AND 6),
    weekend_day SMALLINT NOT NULL DEFAULT 0 CHECK (weekend_day BETWEEN 0 AND 6),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO congregation (id) VALUES (1);

CREATE TABLE households (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    address TEXT NOT NULL DEFAULT '',
    notes TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE publishers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    household_id UUID REFERENCES households (id) ON DELETE SET NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    gender TEXT NOT NULL CHECK (gender IN ('male', 'female')),
    phone TEXT NOT NULL DEFAULT '',
    email TEXT NOT NULL DEFAULT '',
    baptism_date DATE,
    started_preaching_date DATE,
    spiritual_status TEXT NOT NULL DEFAULT 'publisher'
        CHECK (spiritual_status IN ('student', 'unbaptized_publisher', 'publisher')),
    is_elder BOOLEAN NOT NULL DEFAULT FALSE,
    is_ministerial_servant BOOLEAN NOT NULL DEFAULT FALSE,
    is_regular_pioneer BOOLEAN NOT NULL DEFAULT FALSE,
    is_special_pioneer BOOLEAN NOT NULL DEFAULT FALSE,
    activity_status TEXT NOT NULL DEFAULT 'inactive'
        CHECK (activity_status IN ('regular', 'irregular', 'inactive')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX publishers_household_id_idx ON publishers (household_id);
CREATE INDEX publishers_last_name_idx ON publishers (last_name, first_name);

CREATE TABLE field_service_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    publisher_id UUID NOT NULL REFERENCES publishers (id) ON DELETE CASCADE,
    year SMALLINT NOT NULL CHECK (year BETWEEN 2000 AND 2100),
    month SMALLINT NOT NULL CHECK (month BETWEEN 1 AND 12),
    shared_in_ministry BOOLEAN NOT NULL DEFAULT FALSE,
    bible_studies INTEGER NOT NULL DEFAULT 0 CHECK (bible_studies >= 0),
    hours DOUBLE PRECISION,
    auxiliary_pioneer BOOLEAN NOT NULL DEFAULT FALSE,
    late BOOLEAN NOT NULL DEFAULT FALSE,
    remarks TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (publisher_id, year, month)
);

CREATE INDEX field_service_reports_period_idx ON field_service_reports (year, month);

CREATE TABLE meeting_attendance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    meeting_date DATE NOT NULL,
    kind TEXT NOT NULL CHECK (kind IN ('midweek', 'weekend')),
    in_person INTEGER NOT NULL CHECK (in_person >= 0),
    online INTEGER CHECK (online IS NULL OR online >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (meeting_date, kind)
);

CREATE INDEX meeting_attendance_date_idx ON meeting_attendance (meeting_date);
