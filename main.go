package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"PromMySQL/utils"

	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// MySQL's connection details
	dbHost     = utils.GetEnvOrDefaultString("MYSQL_HOST", "127.0.0.1")
	dbPort     = utils.GetEnvOrDefaultInt("MYSQL_PORT", 3306)
	dbUsername = utils.GetEnvOrDefaultString("MYSQL_USERNAME", "username")
	dbPassword = utils.GetEnvOrDefaultString("MYSQL_PASSWORD", "password")
	dbName     = utils.GetEnvOrDefaultString("MYSQL_DB_NAME", "test_db")

	sleepIntervalMillisecond = utils.GetEnvOrDefaultInt("SLEEP_INTERVAL_MILLISECOND", 1000)
)

var (
	// Prometheus metrics
	connections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mysql_connections",
			Help: "Number of MySQL connections",
		},
		[]string{"status"},
	)

	throughput = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mysql_throughput",
			Help: "Number of data entries pushed to MySQL",
		},
		[]string{},
	)

	successCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "mysql_success_count",
			Help: "Number of successful data pushes to MySQL",
		},
	)

	failCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "mysql_fail_count",
			Help: "Number of failed data pushes to MySQL",
		},
	)
)

func init() {
	// Register metrics with Prometheus
	prometheus.MustRegister(connections)
	prometheus.MustRegister(throughput)
	prometheus.MustRegister(successCount)
	prometheus.MustRegister(failCount)
}

func main() {
	// MySQL Connection
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbUsername, dbPassword, dbHost, dbPort, dbName),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// HTTP server for Prometheus metrics endpoint
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Create the table if it doesn't exist
	err = createTableIfNotExists(db)

	if err != nil {
		log.Fatal("Error creating table:", err)
	}

	// Clear the table
	err = clearData(db)
	if err != nil {
		log.Fatal("Error clearing data:", err)
	}

	count := 0

	// Infinite loop for continuous data push
	for {
		// Example: Insert a record into the table
		err := insertData(db, count)
		if err != nil {
			log.Println("Error inserting data:", err)
			failCount.Inc()
		} else {
			successCount.Inc()
		}
		count++

		// Increment throughput metric
		throughput.WithLabelValues().Inc()

		// Update and export the number of active connections
		activeConnections, err := getActiveConnections(db)
		if err != nil {
			log.Println("Error getting active connections:", err)
		} else {
			connections.WithLabelValues("active").Set(float64(activeConnections))
		}

		// Sleep for a defined interval (e.g., 1 second)
		time.Sleep(time.Duration(sleepIntervalMillisecond) * time.Millisecond)
	}
}

func insertData(db *sql.DB, value int) error {
	_, err := db.Exec("INSERT INTO `table_test` (number) VALUES (?)", value)
	return err
}

func getActiveConnections(db *sql.DB) (int, error) {
	var connections int
	var variableName string
	err := db.QueryRow("SHOW STATUS LIKE 'Threads_connected'").Scan(&variableName, &connections)
	if err != nil {
		return 0, err
	}
	return connections, nil
}

// CREATE TABLE IF NOT EXISTS table_test
// (
//
//	id         INT PRIMARY KEY AUTO_INCREMENT,
//	number     INT,
//	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
//	deleted_at TIMESTAMP DEFAULT NULL
//
// );
func createTableIfNotExists(db *sql.DB) error {
	// Create the table if it doesn't exist
	_, err := db.Exec(
		fmt.Sprintf(
			`CREATE TABLE IF NOT EXISTS table_test
(
    id         INT PRIMARY KEY AUTO_INCREMENT,
    number     INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);`,
		),
	)
	return err
}

// clearData truncates the table and resets the auto-increment ID
func clearData(db *sql.DB) error {
	_, err := db.Exec("TRUNCATE TABLE `table_test`")
	return err
}
