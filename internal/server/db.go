package server

import (
	"database/sql"
	"fmt"
	"log"
	"monitoring/internal/models"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("Error loading .env file, proceeding with system environment variables")
	}
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		return nil, fmt.Errorf("database configuration variables are not set properly")
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Error opening database: %v\n", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Printf("Error connecting to the database: %v\n", err)
		return nil, err
	}

	log.Println("Successfully connected to the database")
	return db, nil
}

func SaveMetric(db *sql.DB, metric models.Metric) error {
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
