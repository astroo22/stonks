package services

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"stonks/models"
	"stonks/utils"
	"time"

	_ "github.com/lib/pq"
)

type DBService struct {
	db *sql.DB
}

func NewDBService(connStr string) (*DBService, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &DBService{db: db}, nil
}
func (d *DBService) DB() *sql.DB {
	return d.db
}

func (d *DBService) SaveTickerInfo(info models.TickerInfo) error {
	// placeholder query will fix after I know whats returned from schwab
	query := `INSERT INTO tickers (ticker, price) VALUES ($1, $2) ON CONFLICT (ticker) DO UPDATE SET price = $2`
	_, err := d.db.Exec(query, info.Ticker, info.Price)
	return err
}

func (d *DBService) SaveMultipleTickersInfo(infos []models.TickerInfo) error {
	for _, info := range infos {
		// placeholder query will make this performant with a batch query after I know whats returned from schwab
		if err := d.SaveTickerInfo(info); err != nil {
			return err
		}
	}
	return nil
}

func StartPostgresServer(config utils.PostgresConfig) error {
	// Check if the data directory is already initialized
	// Assuming the existence of postgresql.conf means the directory is initialized
	if _, err := os.Stat(config.DataDir + "/postgresql.conf"); err == nil {
		log.Println("PostgreSQL data directory already initialized. Skipping initialization.")
	} else {
		// Initialize the PostgreSQL data directory if not already done
		initCmd := exec.Command(config.BinaryPath+"/pg_ctl", "initdb", "-D", config.DataDir)
		initOutput, err := initCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to initialize postgres data directory: %v - %s", err, string(initOutput))
		}
		log.Println("PostgreSQL data directory initialized successfully.")
	}

	// Start the PostgreSQL server
	startCmd := exec.Command(config.BinaryPath+"/pg_ctl", "start", "-D", config.DataDir)
	startOutput, err := startCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start postgres server: %v - %s", err, string(startOutput))
	}

	log.Printf("PostgreSQL server start output: %s", string(startOutput))
	log.Println("PostgreSQL server started successfully.")

	// Give the server some time to start
	time.Sleep(5 * time.Second)

	// Verify the server is running by attempting to connect to it
	statusCmd := exec.Command(config.BinaryPath+"/pg_ctl", "status", "-D", config.DataDir)
	statusOutput, err := statusCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to check postgres server status: %v - %s", err, string(statusOutput))
	}

	log.Printf("PostgreSQL server status: %s", string(statusOutput))

	return nil
}

func LoadPostgresServer(config utils.PostgresConfig) error {
	cmd := exec.Command(config.BinaryPath+"/pg_ctl", "start", "-D", config.LoadPath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to load postgres server: %v", err)
	}

	time.Sleep(5 * time.Second) // Give some time for the server to start

	return nil
}

// CheckIfPostgresIsRunning attempts to connect to the database to check if PostgreSQL is running.
func CheckIfPostgresIsRunning(connStr string) (bool, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return false, fmt.Errorf("failed to open database connection: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return false, nil // If ping fails, it means the server isn't running
	}

	return true, nil // Server is running
}
