CREATE DATABASE servicecatalog_test;
\c servicecatalog_test;

CREATE TABLE IF NOT EXISTS services (
                                        id SERIAL PRIMARY KEY,
                                        name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
                             );

CREATE TABLE IF NOT EXISTS versions (
                                        id SERIAL PRIMARY KEY,
                                        service_id INTEGER REFERENCES services(id),
    number VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
                             );

CREATE INDEX idx_services_name ON services(name);
CREATE INDEX idx_versions_service_id ON versions(service_id);

-- Add to migrations or setup.sql
ALTER TABLE services ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE versions ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;

CREATE INDEX idx_services_deleted_at ON services(deleted_at);
CREATE INDEX idx_versions_deleted_at ON versions(deleted_at);