package grpcserver

import (
	"context"
	"database/sql"
	"monitoring/proto"
)

type MetricService struct {
	proto.UnimplementedMetricServiceServer
	DB *sql.DB
}

// retrieves all metrics from the database
func (s *MetricService) ListMetrics(ctx context.Context, req *proto.MetricRequest) (*proto.MetricList, error) {
	// query all metrics for a specific hostname with a limit
	query := `
		SELECT m.id, m.host_id, m.uptime, m.cpu, m.ram, m.disk, m.network, m.time
		FROM metrics m
		JOIN hosts h ON m.host_id = h.id
		WHERE h.hostname = $1
		ORDER BY m.time DESC
		LIMIT $2
	`
	rows, err := s.DB.Query(query, req.Hostname, req.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []*proto.Metric
	for rows.Next() {
		var metric proto.Metric
		if err := rows.Scan(&metric.Id, &metric.HostId, &metric.Uptime, &metric.Cpu, &metric.Ram, &metric.Disk, &metric.Network, &metric.Time); err != nil {
			return nil, err
		}
		metrics = append(metrics, &metric)
	}

	return &proto.MetricList{Metrics: metrics}, nil
}

// retrieves the latest metrics
func (s *MetricService) GetLatestMetrics(ctx context.Context, _ *proto.Empty) (*proto.MetricList, error) {
	// query latest metrics for all hosts
	query := `
		SELECT DISTINCT ON (m.host_id) m.id, m.host_id, m.uptime, m.cpu, m.ram, m.disk, m.network, m.time
		FROM metrics m
		ORDER BY m.host_id, m.time DESC
	`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []*proto.Metric
	for rows.Next() {
		var metric proto.Metric
		if err := rows.Scan(&metric.Id, &metric.HostId, &metric.Uptime, &metric.Cpu, &metric.Ram, &metric.Disk, &metric.Network, &metric.Time); err != nil {
			return nil, err
		}
		metrics = append(metrics, &metric)
	}

	return &proto.MetricList{Metrics: metrics}, nil
}
