-- Create enum type for relationship
CREATE TYPE relation AS ENUM ('Friend', 'Family', 'Colleague', 'School', 'Network', 'Services');

CREATE TABLE contacts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    relationship relation,
    industry VARCHAR(255),
    company VARCHAR(255),
    birthday VARCHAR(255),
    vip BOOLEAN DEFAULT FALSE,
    -- FamilyDetails fields
    spouse VARCHAR(255),
    children TEXT,
    -- ContactInfo fields
    location VARCHAR(255),
    phone_number VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    linked_in VARCHAR(255),
    instagram VARCHAR(255),
    x VARCHAR(255),
    -- Notes
    notes TEXT,
    -- VipInfo fields
    last_met VARCHAR(255),
    last_contacted VARCHAR(255),
    last_update VARCHAR(255),
    status VARCHAR(255),
    -- Calendar sync fields
    google_calendar_event_id VARCHAR(255),
    calendar_sync_enabled BOOLEAN DEFAULT FALSE,
    calendar_synced_at TIMESTAMP,
    -- Standard timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_contacts_user_id ON contacts(user_id);
CREATE INDEX idx_contacts_email ON contacts(email);
CREATE INDEX idx_contacts_phone_number ON contacts(phone_number);
CREATE INDEX idx_contacts_deleted_at ON contacts(deleted_at);
