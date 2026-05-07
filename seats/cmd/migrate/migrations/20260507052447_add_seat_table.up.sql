CREATE TABLE IF NOT EXISTS seats (
    id SERIAL PRIMARY KEY,
    seat_number VARCHAR(10) NOT NULL,
    column_number VARCHAR(10) NOT NULL,
    row_number VARCHAR(10) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);