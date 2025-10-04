package grpcserver

import (
	"context"
	"database/sql"
	"monitoring/internal/grpcserver/gen"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// implements the gRPC server for handling metrics-related requests
type MetricsServer struct {
	gen.UnimplementedMetricsServiceServer
	Repo *Repository
}

func NewMetricsServer(repo *Repository) *MetricsServer {
	return &MetricsServer{
		Repo: repo,
	}
}

// returns a list of metrics based on filtering criteria such as hostname, time range, limit, and offset
func (s *MetricsServer) GetListMetrics(ctx context.Context, req *gen.GetListMetricsRequest) (*gen.GetListMetricsResponse, error) {
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

	// getting metrics from repository
	metrics, err := s.Repo.GetMetrics(ctx, req.Hostname, fromTime, toTime, limit, offset)
	if err != nil {
		return nil, err
	}

	return &gen.GetListMetricsResponse{Metrics: metrics}, nil
}

// returns the latest metric for a specific hostname
func (s *MetricsServer) GetLatestMetrics(ctx context.Context, req *gen.GetLatestMetricsRequest) (*gen.GetLatestMetricsResponse, error) {
	// validate required hostname parameter
	if req.Hostname == "" {
		return nil, status.Errorf(codes.InvalidArgument, "hostname is required")
	}

	// getting the latest metric from repository
	metrics, err := s.Repo.LatestMetrics(ctx, req.Hostname)
	if err != nil {
		return nil, err
	}

	if len(metrics) == 0 {
		return &gen.GetLatestMetricsResponse{}, nil
	}

	return &gen.GetLatestMetricsResponse{Metrics: metrics[0]}, nil
}

// returns all metrics for a specific hostname
func (s *MetricsServer) GetMetricsByHost(ctx context.Context, req *gen.GetMetricsByHostRequest) (*gen.GetMetricsByHostResponse, error) {
	// validate required hostname parameter
	if req.Hostname == "" {
		return nil, status.Errorf(codes.InvalidArgument, "hostname is required")
	}

	// getting metrics by host from repository
	metrics, err := s.Repo.MetricsByHost(ctx, req.Hostname)
	if err != nil {
		return nil, err
	}

	return &gen.GetMetricsByHostResponse{Metrics: metrics}, nil
}
