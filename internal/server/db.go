package server

import (
	"database/sql"
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
		log.Printf("Error connecting to the database: %v\n", err)
		return nil, err
	}

	log.Println("Successfully connected to the database")
	return db, nil
}

func SaveMetric(db *sql.DB, metric models.Metric) error {
	// insert metric into database
	_, err := db.Exec(
		`INSERT INTO metrics (hostname, os, platform, platform_ver, kernel_ver, uptime, cpu, ram, time) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		metric.Hostname,
		metric.OS,
		metric.Platform,
		metric.PlatformVer,
		metric.KernelVer,
		metric.Uptime,
		metric.CPU,
		metric.RAM,
		metric.Time,
	)

	if err != nil {
		log.Printf("Error inserting metric: %v\n", err)
	}

	return err
}
