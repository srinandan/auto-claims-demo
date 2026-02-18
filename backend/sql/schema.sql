-- Schema for Claims App (PostgreSQL)

CREATE TABLE IF NOT EXISTS claims (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    policy_number TEXT,
    customer_name TEXT,
    status TEXT,
    description TEXT,
    accident_date TIMESTAMP WITH TIME ZONE,
    incident_city TEXT,
    incident_state TEXT,
    incident_type TEXT,
    collision_type TEXT,
    severity TEXT,
    total_loss_probability DOUBLE PRECISION,
    action_required BOOLEAN
);

CREATE INDEX IF NOT EXISTS idx_claims_deleted_at ON claims(deleted_at);

CREATE TABLE IF NOT EXISTS photos (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    claim_id INTEGER,
    url TEXT,
    CONSTRAINT fk_claims_photos FOREIGN KEY (claim_id) REFERENCES claims(id)
);

CREATE INDEX IF NOT EXISTS idx_photos_deleted_at ON photos(deleted_at);

CREATE TABLE IF NOT EXISTS analysis_results (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    photo_id INTEGER,
    quality_score TEXT,
    detections TEXT,
    parts_detected TEXT,
    severity TEXT,
    CONSTRAINT fk_photos_analysis_results FOREIGN KEY (photo_id) REFERENCES photos(id)
);

CREATE INDEX IF NOT EXISTS idx_analysis_results_deleted_at ON analysis_results(deleted_at);

CREATE TABLE IF NOT EXISTS estimates (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    claim_id INTEGER,
    total_amount DOUBLE PRECISION,
    items TEXT,
    source TEXT,
    CONSTRAINT fk_claims_estimates FOREIGN KEY (claim_id) REFERENCES claims(id)
);

CREATE INDEX IF NOT EXISTS idx_estimates_deleted_at ON estimates(deleted_at);

CREATE TABLE IF NOT EXISTS policy_holders (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    policy_number TEXT UNIQUE,
    first_name TEXT,
    last_name TEXT,
    months_as_customer INTEGER,
    age INTEGER,
    policy_bind_date TIMESTAMP WITH TIME ZONE,
    policy_state TEXT,
    policy_csl TEXT,
    policy_deductible INTEGER,
    policy_annual_premium DOUBLE PRECISION,
    umbrella_limit INTEGER,
    insured_zip INTEGER,
    insured_sex TEXT,
    insured_education_level TEXT,
    insured_occupation TEXT,
    insured_hobbies TEXT,
    insured_relationship TEXT,
    capital_gains INTEGER,
    capital_loss INTEGER,
    auto_make TEXT,
    auto_model TEXT,
    auto_year INTEGER
);

CREATE INDEX IF NOT EXISTS idx_policy_holders_deleted_at ON policy_holders(deleted_at);
