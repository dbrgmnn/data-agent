package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"monitoring/internal/config"
	"monitoring/internal/models"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	// load config from .env file
	cfg := config.LoadConfig()
	if cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBUser == "" || cfg.DBPass == "" || cfg.DBName == "" {
		return nil, fmt.Errorf("database configuration variables are not set properly")
	}

	// make connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)

	// open connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Error opening database: %v\n", err)
		return nil, err
	}

	// test connection
	err = db.Ping()
	if err != nil {
		db.Close()
		log.Printf("Error connecting to the database: %v\n", err)
		return nil, err
	}

	log.Println("Successfully connected to the database")
	return db, nil
}

// insert host info into database
func SaveHostInfo(db *sql.DB, metric *models.MetricMessage) error {
	var HostID int64

	err := db.QueryRow(
		`INSERT INTO hosts (hostname, os, platform, platform_ver, kernel_ver) 
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (hostname) DO UPDATE SET
		   os = EXCLUDED.os,
		   platform = EXCLUDED.platform,
		   platform_ver = EXCLUDED.platform_ver,
		   kernel_ver = EXCLUDED.kernel_ver
		 RETURNING id`,
		metric.Host.Hostname,
		metric.Host.OS,
		metric.Host.Platform,
		metric.Host.PlatformVer,
		metric.Host.KernelVer,
	).Scan(&HostID)

	if err != nil {
		log.Printf("Error inserting host info: %v\n", err)
		return err
	}

	metric.Metric.HostID = HostID
	return nil
}

// insert metric into database
func SaveMetric(db *sql.DB, metric models.MetricMessage) error {
	// marshaling disk and network slices to JSON
	diskJSON, err := json.Marshal(metric.Metric.Disk)
	if err != nil {
		return fmt.Errorf("error marshaling disk metrics: %v", err)
	}
	networkJSON, err := json.Marshal(metric.Metric.Network)
	if err != nil {
		return fmt.Errorf("error marshaling network metrics: %v", err)
	}
	_, err = db.Exec(
		`INSERT INTO metrics (host_id, uptime, cpu, ram, disk, network, time)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		metric.Metric.HostID,
		metric.Metric.Uptime,
		metric.Metric.CPU,
		metric.Metric.RAM,
		diskJSON,
		networkJSON,
		metric.Metric.Time,
	)
	if err != nil {
		log.Printf("Error inserting metric: %v\n", err)
	}

	return err
}
