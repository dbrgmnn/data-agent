CREATE TABLE IF NOT EXISTS metrics (
    id SERIAL PRIMARY KEY,
    hostname TEXT,
    os TEXT,
    platform TEXT,
    platform_ver TEXT,
    kernel_ver TEXT,
    uptime BIGINT,
    cpu NUMERIC,
    ram NUMERIC,
    disk JSONB,
    network JSONB,
    time TIMESTAMPTZ
);
