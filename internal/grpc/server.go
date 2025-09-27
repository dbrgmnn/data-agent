package grpcserver

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

// implements the gRPC server for handling metrics-related requests
type MetricsServer struct {
	UnimplementedMetricsServiceServer
	DB *sql.DB
}

func NewMetricsServer(db *sql.DB) *MetricsServer {
	return &MetricsServer{DB: db}
}

// retrieves the latest metrics from the database
func (s *MetricsServer) GetLatestMetrics(ctx context.Context, req *GetMetricsRequest) (*GetMetricsResponse, error) {
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

	var metrics []*Metric
	for rows.Next() {
		var m Metric
		var diskStr, networkStr string
		if err := rows.Scan(&m.Hostname, &m.Os, &m.Platform, &m.PlatformVer,
			&m.KernelVer, &m.Uptime, &m.Cpu, &m.Ram, &diskStr, &networkStr, &m.Time); err != nil {
			return nil, fmt.Errorf("scanning metric row: %w", err)
		}

		// unmarshal disk and network JSON strings into appropriate fields
		if err := json.Unmarshal([]byte(diskStr), &m.Disk); err != nil {
			return nil, fmt.Errorf("unmarshaling disk JSON: %w", err)
		}

		if err := json.Unmarshal([]byte(networkStr), &m.Network); err != nil {
			return nil, fmt.Errorf("unmarshaling network JSON: %w", err)
		}

		metrics = append(metrics, &m)
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("iterating metric rows: %w", err)
		}
	}

	return &GetMetricsResponse{Metrics: metrics}, nil
}
