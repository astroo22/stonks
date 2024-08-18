package main

import (
	"database/sql"
	"log"
	"net/http"
	"stonks/handlers"
	"stonks/services"
	"stonks/utils"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Goal is to build this so it can set up a postgres server itself so non devs can use but will have too wait.
// nothing is tested I haven't got api key yet.

func main() {
	// Load server and tickers configuration
	serverConfig, tickersConfig := loadConfigs()

	// Ensure PostgreSQL is running
	ensurePostgresIsRunning(serverConfig)

	// Initialize services
	schwabService := initializeSchwabService(serverConfig.APIKey)
	dbService := initializeDBService(serverConfig.DBConnStr)

	// Set up HTTP handlers and start the server
	setupAndStartServer(serverConfig.Port, schwabService, dbService, tickersConfig)
}

func loadConfigs() (*utils.ServerConfig, *utils.TickersConfig) {
	// Load server config
	serverConfig, err := utils.LoadServerConfig("config/server_config.yaml")
	if err != nil {
		log.Fatal("Error loading server config:", err)
	}
	//fmt.Println(serverConfig)

	// Load tickers config
	tickersConfig, err := utils.LoadTickersConfig("config/tickers_config.yaml")
	if err != nil {
		log.Fatal("Error loading tickers config:", err)
	}
	log.Println("Loaded configs sucessfully")
	return serverConfig, tickersConfig
}

func ensurePostgresIsRunning(serverConfig *utils.ServerConfig) {
	isRunning, err := services.CheckIfPostgresIsRunning(serverConfig.DBConnStr)
	if err != nil {
		log.Fatal("Error checking PostgreSQL status:", err)
	}

	if !isRunning {
		log.Println("PostgreSQL is not running. Starting PostgreSQL server...")

		switch serverConfig.Postgres.Mode {
		case "startup":
			if err := services.StartPostgresServer(serverConfig.Postgres); err != nil {
				log.Fatal("Failed to start PostgreSQL server:", err)
			}
			utils.SaveServerConfig("config/server_config.yaml", serverConfig)
		case "load":
			log.Println("Loading PostgreSQL server...")
			if err := services.LoadPostgresServer(serverConfig.Postgres); err != nil {
				log.Fatal("Failed to load PostgreSQL server:", err)
			}
		default:
			log.Fatal("Invalid PostgreSQL mode in config")
		}
	} else {
		log.Println("PostgreSQL is already running. Connecting to the database...")
	}
}

func initializeSchwabService(apiKey string) *services.SchwabService {
	return services.NewSchwabService(apiKey)
}

func initializeDBService(dbConnStr string) *services.DBService {
	dbService, err := services.NewDBService(dbConnStr)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	// Perform a database connection check
	if err := checkDBConnection(dbService.DB()); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	return dbService
}

func setupAndStartServer(port string, schwabService *services.SchwabService, dbService *services.DBService, tickersConfig *utils.TickersConfig) {
	// Set up handlers
	schwabHandler := handlers.NewSchwabHandler(schwabService, dbService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/ticker", schwabHandler.GetTickerInfo)
	r.Post("/tickers", schwabHandler.GetMultipleTickersInfo)

	// Schedule jobs if enabled
	if tickersConfig.Job.Enabled {
		utils.ScheduleTask(time.Duration(tickersConfig.Job.RefreshIntervalMinutes)*time.Minute, func() {
			if shouldRunJob(tickersConfig.Job.TradingHoursOnly) {
				infos, err := schwabService.GetMultipleTickersInfo(tickersConfig.Tickers)
				if err != nil {
					log.Println("Error refreshing tickers:", err)
					return
				}
				dbService.SaveMultipleTickersInfo(infos)
			}
		})
	}

	// Start the HTTP server
	log.Printf("Server started on %s", port)
	log.Fatal(http.ListenAndServe(port, r))
}

func shouldRunJob(tradingHoursOnly bool) bool {
	if !tradingHoursOnly {
		return true
	}

	now := time.Now()
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		return false
	}

	startTime := time.Date(now.Year(), now.Month(), now.Day(), 9, 30, 0, 0, now.Location())
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 16, 0, 0, 0, now.Location())

	return now.After(startTime) && now.Before(endTime)
}

func checkDBConnection(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}
	log.Println("Successfully connected to the database")
	return nil
}
