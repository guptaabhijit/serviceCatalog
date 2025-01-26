-- db/setup.sql

-- Create Database
CREATE DATABASE servicecatalog;

-- Connect to database
\c servicecatalog;

-- Create Services Table
CREATE TABLE IF NOT EXISTS services (
                                        id SERIAL PRIMARY KEY,
                                        name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
                             );

-- Create Versions Table
CREATE TABLE IF NOT EXISTS versions (
                                        id SERIAL PRIMARY KEY,
                                        service_id INTEGER REFERENCES services(id),
    number VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
                             );

-- Create Indexes
CREATE INDEX idx_services_name ON services(name);
CREATE INDEX idx_versions_service_id ON versions(service_id);

-- Insert Sample Data
INSERT INTO services (id, name, description) VALUES
                                                 (1, 'Authentication Service', 'Handles user authentication and authorization across all platforms'),
                                                 (2, 'Payment Gateway', 'Processes payments and manages transaction workflows'),
                                                 (3, 'User Management', 'Manages user profiles, preferences, and settings'),
                                                 (4, 'Notification Service', 'Handles email, SMS, and push notifications'),
                                                 (5, 'Analytics Service', 'Collects and processes user behavior and system metrics'),
                                                 (6, 'Content Delivery', 'Manages and serves static and dynamic content'),
                                                 (7, 'Search Service', 'Provides full-text search capabilities across all content'),
                                                 (8, 'Inventory Management', 'Tracks and manages product inventory in real-time'),
                                                 (9, 'Recommendation Engine', 'Provides personalized content and product recommendations'),
                                                 (10, 'Logging Service', 'Centralized logging and monitoring solution');

-- Insert Sample Versions
INSERT INTO versions (id, service_id, number) VALUES
                                                  (1, 1, '1.0.0'),
                                                  (2, 1, '1.1.0'),
                                                  (3, 1, '2.0.0'),
                                                  (4, 2, '1.0.0'),
                                                  (5, 2, '1.5.0'),
                                                  (6, 2, '2.0.0'),
                                                  (7, 2, '2.1.0'),
                                                  (8, 3, '1.0.0'),
                                                  (9, 3, '1.1.0'),
                                                  (10, 3, '1.2.0'),
                                                  (11, 4, '1.0.0'),
                                                  (12, 4, '2.0.0'),
                                                  (13, 5, '1.0.0'),
                                                  (14, 5, '1.1.0'),
                                                  (15, 5, '1.2.0'),
                                                  (16, 5, '2.0.0'),
                                                  (17, 6, '1.0.0'),
                                                  (18, 6, '1.5.0'),
                                                  (19, 7, '1.0.0'),
                                                  (20, 7, '1.1.0'),
                                                  (21, 7, '2.0.0'),
                                                  (22, 8, '1.0.0'),
                                                  (23, 8, '1.1.0'),
                                                  (24, 9, '1.0.0'),
                                                  (25, 9, '1.5.0'),
                                                  (26, 9, '2.0.0'),
                                                  (27, 10, '1.0.0'),
                                                  (28, 10, '1.1.0');