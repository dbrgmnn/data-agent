package grpcserver

import (
	"context"
	"database/sql"
	"encoding/json"
	"monitoring/internal/grpcserver/gen"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const baseSelect = `SELECT hostname, os, platform, platform_ver,
    kernel_ver, uptime, cpu, ram, disk::text, network::text, time
	FROM metrics 
	`

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

// retrieves a list of metrics with optional filtering by hostname and time range, and supports pagination
func (r *Repository) GetMetrics(ctx context.Context, hostname string, fromTime, toTime sql.NullTime, limit, offset int32) ([]*gen.Metric, error) {
	query := baseSelect + `
		WHERE hostname = $1
	  		AND ($2::timestamptz IS NULL OR time >= $2::timestamptz)
	  		AND ($3::timestamptz IS NULL OR time <= $3::timestamptz)
		ORDER BY time DESC
		LIMIT $4
		OFFSET $5;
	`

	metrics, err := r.queryMetrics(ctx, query, hostname, fromTime, toTime, limit, offset)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

// retrieves the latest metric entry for a given hostname
func (r *Repository) LatestMetrics(ctx context.Context, hostname string) ([]*gen.Metric, error) {
	query := baseSelect + `
		WHERE hostname = $1
		ORDER BY time DESC
		LIMIT 1;
	`

	metrics, err := r.queryMetrics(ctx, query, hostname)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

// retrieves all metric entries for a given hostname
func (r *Repository) MetricsByHost(ctx context.Context, hostname string) ([]*gen.Metric, error) {
	query := baseSelect + `
		WHERE hostname = $1
		ORDER BY time DESC;
	`

	metrics, err := r.queryMetrics(ctx, query, hostname)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

// helper for scanning and unmarshaling a metric row
func scanMetrics(rows *sql.Rows) ([]*gen.Metric, error) {
	var metrics []*gen.Metric
	for rows.Next() {
		var m gen.Metric
		var diskJSON, networkJSON string

		// scan row into metric fields
		if err := rows.Scan(&m.Hostname, &m.Os, &m.Platform, &m.PlatformVer,
			&m.KernelVer, &m.Uptime, &m.Cpu, &m.Ram, &diskJSON, &networkJSON, &m.Time); err != nil {
			return nil, status.Errorf(codes.Internal, "scanning metric row: %v", err)
		}

		// unmarshal JSON fields
		if err := json.Unmarshal([]byte(diskJSON), &m.Disk); err != nil {
			return nil, status.Errorf(codes.Internal, "unmarshaling disk JSON: %v", err)
		}
		if err := json.Unmarshal([]byte(networkJSON), &m.Network); err != nil {
			return nil, status.Errorf(codes.Internal, "unmarshaling network JSON: %v", err)
		}

		metrics = append(metrics, &m)
	}
	if err := rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "iterating metric rows: %v", err)
	}
	return metrics, nil
}

// helper for querying metrics
func (r *Repository) queryMetrics(ctx context.Context, query string, args ...any) ([]*gen.Metric, error) {
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "querying metrics: %v", err)
	}
	defer rows.Close()

	// scan rows into Metric structs
	metrics, err := scanMetrics(rows)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}
