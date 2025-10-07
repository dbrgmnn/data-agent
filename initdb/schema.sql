CREATE TABLE IF NOT EXISTS hosts (
    id SERIAL PRIMARY KEY,
    hostname TEXT UNIQUE NOT NULL,
    os TEXT,
    platform TEXT,
    platform_ver TEXT,
    kernel_ver TEXT
);

CREATE TABLE IF NOT EXISTS metrics (
    id SERIAL PRIMARY KEY,
    host_id INT REFERENCES hosts(id) ON DELETE CASCADE,
    uptime BIGINT,
    cpu NUMERIC,
    ram NUMERIC,
    disk JSONB,
    network JSONB,
    time TIMESTAMPTZ
);
