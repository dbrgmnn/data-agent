CREATE TABLE IF NOT EXISTS metrics (
    id SERIAL PRIMARY KEY,
    hostname TEXT,
    os TEXT,
    platform TEXT,
    platform_ver TEXT,
    kernel_ver TEXT,
    cpu NUMERIC,
    ram NUMERIC,
    uptime BIGINT,
    time TIMESTAMPTZ
);
