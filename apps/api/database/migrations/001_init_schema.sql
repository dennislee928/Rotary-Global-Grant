-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table (for triagers, admins, auditors)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'triager' CHECK (role IN ('admin', 'triager', 'auditor', 'educator')),
    display_name VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Reports table (community incident reports)
CREATE TABLE reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category VARCHAR(50) NOT NULL CHECK (category IN (
        'suspicious_item', 'suspicious_person', 'harassment_stalking',
        'scam_phishing', 'misinformation_panic', 'crowd_disorder',
        'infrastructure_hazard', 'other'
    )),
    severity_suggested VARCHAR(10) CHECK (severity_suggested IN ('S0', 'S1', 'S2', 'S3', 'S4')),
    area_hint VARCHAR(500) NOT NULL,
    time_window VARCHAR(100),
    description TEXT NOT NULL,
    evidence_refs JSONB DEFAULT '[]'::jsonb,
    reporter_contact_ref VARCHAR(255),
    status VARCHAR(50) NOT NULL DEFAULT 'submitted' CHECK (status IN (
        'submitted', 'under_review', 'triaged', 'escalated', 'closed', 'spam'
    )),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Triage decisions table
CREATE TABLE triage_decisions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    report_id UUID NOT NULL REFERENCES reports(id) ON DELETE CASCADE,
    decided_by UUID REFERENCES users(id) ON DELETE SET NULL,
    decision VARCHAR(50) NOT NULL CHECK (decision IN ('accept', 'reject', 'needs_more_info', 'escalate')),
    severity_final VARCHAR(10) NOT NULL CHECK (severity_final IN ('S0', 'S1', 'S2', 'S3', 'S4')),
    evidence_level VARCHAR(10) CHECK (evidence_level IN ('E0', 'E1', 'E2', 'E3')),
    rationale TEXT,
    audit_hash VARCHAR(64),
    decided_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Alerts table (CAP-ready alerts)
CREATE TABLE alerts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    report_id UUID REFERENCES reports(id) ON DELETE SET NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'approved', 'published', 'withdrawn')),
    event VARCHAR(255) NOT NULL,
    urgency VARCHAR(50) NOT NULL CHECK (urgency IN ('Immediate', 'Expected', 'Future', 'Past', 'Unknown')),
    severity VARCHAR(50) NOT NULL CHECK (severity IN ('Extreme', 'Severe', 'Moderate', 'Minor', 'Unknown')),
    certainty VARCHAR(50) NOT NULL CHECK (certainty IN ('Observed', 'Likely', 'Possible', 'Unlikely', 'Unknown')),
    area VARCHAR(500) NOT NULL,
    instruction TEXT NOT NULL,
    public_message TEXT,
    cap_xml TEXT,
    channels JSONB DEFAULT '[]'::jsonb,
    approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    published_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Training events table
CREATE TABLE training_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    event_date DATE NOT NULL,
    location VARCHAR(500) NOT NULL,
    audience VARCHAR(255),
    attendance_count INTEGER DEFAULT 0,
    pre_avg NUMERIC(5,2),
    post_avg NUMERIC(5,2),
    notes TEXT,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Training participants (for de-duplication tracking)
CREATE TABLE training_participants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL REFERENCES training_events(id) ON DELETE CASCADE,
    participant_hash VARCHAR(64) NOT NULL,
    pre_score NUMERIC(5,2),
    post_score NUMERIC(5,2),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(event_id, participant_hash)
);

-- Quiz results
CREATE TABLE quiz_results (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID REFERENCES training_events(id) ON DELETE SET NULL,
    participant_hash VARCHAR(64),
    quiz_type VARCHAR(50) NOT NULL CHECK (quiz_type IN ('pre', 'post')),
    score NUMERIC(5,2) NOT NULL,
    max_score NUMERIC(5,2) NOT NULL DEFAULT 100,
    answers JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Audit logs table (immutable)
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    actor_id UUID REFERENCES users(id) ON DELETE SET NULL,
    actor_ip VARCHAR(45),
    action VARCHAR(100) NOT NULL,
    object_type VARCHAR(50) NOT NULL,
    object_id UUID,
    diff JSONB,
    ts TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- API keys table (for basic auth)
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    key_hash VARCHAR(64) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    scopes JSONB DEFAULT '["read"]'::jsonb,
    expires_at TIMESTAMP WITH TIME ZONE,
    last_used_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    is_active BOOLEAN NOT NULL DEFAULT true
);

-- Indexes for performance
CREATE INDEX idx_reports_status ON reports(status);
CREATE INDEX idx_reports_category ON reports(category);
CREATE INDEX idx_reports_created_at ON reports(created_at DESC);
CREATE INDEX idx_triage_decisions_report_id ON triage_decisions(report_id);
CREATE INDEX idx_triage_decisions_decided_at ON triage_decisions(decided_at DESC);
CREATE INDEX idx_alerts_status ON alerts(status);
CREATE INDEX idx_alerts_created_at ON alerts(created_at DESC);
CREATE INDEX idx_training_events_date ON training_events(event_date DESC);
CREATE INDEX idx_audit_logs_ts ON audit_logs(ts DESC);
CREATE INDEX idx_audit_logs_object ON audit_logs(object_type, object_id);
CREATE INDEX idx_api_keys_key_hash ON api_keys(key_hash);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP INDEX IF EXISTS idx_api_keys_key_hash;
DROP INDEX IF EXISTS idx_audit_logs_object;
DROP INDEX IF EXISTS idx_audit_logs_ts;
DROP INDEX IF EXISTS idx_training_events_date;
DROP INDEX IF EXISTS idx_alerts_created_at;
DROP INDEX IF EXISTS idx_alerts_status;
DROP INDEX IF EXISTS idx_triage_decisions_decided_at;
DROP INDEX IF EXISTS idx_triage_decisions_report_id;
DROP INDEX IF EXISTS idx_reports_created_at;
DROP INDEX IF EXISTS idx_reports_category;
DROP INDEX IF EXISTS idx_reports_status;

DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS quiz_results;
DROP TABLE IF EXISTS training_participants;
DROP TABLE IF EXISTS training_events;
DROP TABLE IF EXISTS alerts;
DROP TABLE IF EXISTS triage_decisions;
DROP TABLE IF EXISTS reports;
DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS "uuid-ossp";
