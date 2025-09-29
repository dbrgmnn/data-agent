package grpcserver

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"monitoring/internal/grpcserver/gen"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// implements the gRPC server for handling metrics-related requests
type MetricsServer struct {
	gen.UnimplementedMetricsServiceServer
	DB *sql.DB
}

func NewMetricsServer(db *sql.DB) *MetricsServer {
	return &MetricsServer{DB: db}
}

// retrieves a list of metrics with optional filtering by hostname and time range, and supports pagination
func (s *MetricsServer) ListMetrics(ctx context.Context, req *gen.ListMetricsRequest) (*gen.ListMetricsResponse, error) {
	query := `SELECT id, hostname, os, platform, platform_ver,
	kernel_ver, uptime, cpu, ram, disk::text, network::text, time
	FROM metrics
	WHERE hostname = $1
	  AND ($2::timestamptz IS NULL OR time >= $2::timestamptz)
	  AND ($3::timestamptz IS NULL OR time <= $3::timestamptz)
	ORDER BY time DESC
	LIMIT $4
	OFFSET $5;
	`

	// validate required hostname parameter
	if req.Hostname == "" {
		return nil, status.Errorf(codes.InvalidArgument, "hostname is required")
	}

	// parse optional time parameters
	var fromTime, toTime sql.NullTime
	if req.FromTime != "" {
		parsedTime, err := time.Parse(time.RFC3339, req.FromTime)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid from_time format: %v", err)
		}
		fromTime = sql.NullTime{Time: parsedTime, Valid: true}
	}
	if req.ToTime != "" {
		parsedTime, err := time.Parse(time.RFC3339, req.ToTime)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid to_time format: %v", err)
		}
		toTime = sql.NullTime{Time: parsedTime, Valid: true}
	}

	// set default limit
	limit := int32(10)
	if req.Limit > 0 {
		limit = req.Limit
	}

	// set default offset
	offset := int32(0)
	if req.Offset > 0 {
		offset = req.Offset
	}

	// execute query with parameters
	rows, err := s.DB.QueryContext(ctx, query, req.Hostname, fromTime, toTime, limit, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "querying metrics: %v", err)
	}
	defer rows.Close()

	// scan rows into Metric structs
	metrics, err := scanMetrics(rows)
	if err != nil {
		return nil, err
	}

	return &gen.ListMetricsResponse{Metrics: metrics}, nil
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
