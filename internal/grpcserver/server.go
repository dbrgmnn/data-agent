package grpcserver

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"monitoring/internal/grpcserver/gen"
)

// implements the gRPC server for handling metrics-related requests
type MetricsServer struct {
	gen.UnimplementedMetricsServiceServer
	DB *sql.DB
}

func NewMetricsServer(db *sql.DB) *MetricsServer {
	return &MetricsServer{DB: db}
}

// retrieves the latest metrics from the database
func (s *MetricsServer) GetLatestMetrics(ctx context.Context, req *gen.GetMetricsRequest) (*gen.GetMetricsResponse, error) {
	query := `SELECT id, hostname, os, platform, platform_ver,
	kernel_ver, uptime, cpu, ram, disk::text, network::text, time
	FROM metrics
	ORDER BY time DESC
	LIMIT $1;
	`
	rows, err := s.DB.QueryContext(ctx, query, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("querying latest metrics: %w", err)
	}
	defer rows.Close()

	metrics, err := scanMetrics(rows)
	if err != nil {
		return nil, err
	}

	return &gen.GetMetricsResponse{Metrics: metrics}, nil
}

// retrieves metrics for a specific hostname from the database
func (s *MetricsServer) GetMetricsByHost(ctx context.Context, req *gen.GetMetricsByHostRequest) (*gen.GetMetricsResponse, error) {
	query := `SELECT id, hostname, os , platform, platform_ver,
	kernel_ver, uptime, cpu, ram, disk::text, network::text, time
	FROM metrics
	WHERE hostname = $1
	ORDER BY time DESC
	LIMIT $2;
	`
	rows, err := s.DB.QueryContext(ctx, query, req.Hostname, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("querying metrics by host: %w", err)
	}
	defer rows.Close()

	metrics, err := scanMetrics(rows)
	if err != nil {
		return nil, err
	}

	return &gen.GetMetricsResponse{Metrics: metrics}, nil
}

// helper for scanning and unmarshaling a metric row
func scanMetrics(rows *sql.Rows) ([]*gen.Metric, error) {
	var metrics []*gen.Metric
	for rows.Next() {
		var m gen.Metric
		var id int
		var diskJSON, networkJSON string

		// scan row into metric fields
		if err := rows.Scan(&id, &m.Hostname, &m.Os, &m.Platform, &m.PlatformVer,
			&m.KernelVer, &m.Uptime, &m.Cpu, &m.Ram, &diskJSON, &networkJSON, &m.Time); err != nil {
			return nil, fmt.Errorf("scanning metric row: %w", err)
		}

		// unmarshal JSON fields
		if err := json.Unmarshal([]byte(diskJSON), &m.Disk); err != nil {
			return nil, fmt.Errorf("unmarshaling disk JSON: %w", err)
		}
		if err := json.Unmarshal([]byte(networkJSON), &m.Network); err != nil {
			return nil, fmt.Errorf("unmarshaling network JSON: %w", err)
		}

		metrics = append(metrics, &m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating metric rows: %w", err)
	}
	return metrics, nil
}
