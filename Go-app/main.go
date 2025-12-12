package main

import (
    "fmt"
    "log"
    "os"
    "net/http"
    "time"
	"io" // Added for the worker loop only
)

// The constant for the Worker's internal listening port
const workerInternalPort = "8081" 

func main() {
    appRole := os.Getenv("APP_ROLE")

    if appRole == "worker" {
        startWorker()
    } else if appRole == "web" || appRole == "" {
        startWebServer()
    } else {
        log.Fatalf("FATAL: Unknown APP_ROLE environment variable: %s", appRole)
    }
}

// 1. Worker Process: Runs an internal HTTP server for status
func startWorker() {
    // Note: In a real app, this function would also contain your background job loop.
    
    // Internal API endpoint for the web app to query
    http.HandleFunc("/worker-status", func(w http.ResponseWriter, r *http.Request) {
        // Simple fixed string response
        fmt.Fprintf(w, "Worker is actively processing jobs.") 
    })
    
    fmt.Printf("Worker API listening internally on :%s\n", workerInternalPort)
    // Keep this running indefinitely
    log.Fatal(http.ListenAndServe(":"+workerInternalPort, nil))
}

// ... (imports remain the same, ensure you have "net/http", "io", and "time" imported)

// IMPORTANT: Replace 'your-service-name' with the 'service:' name from your deploy.yml
const workerURL = "http://13.127.117.248:8081/worker-status" 

func startWebServer() {
    fmt.Println("Starting Web Server...")
    
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        workerStatus, err := fetchWorkerStatus()
        
        fmt.Fprintf(w, "<h1>Hello from the Web Server!</h1>")

        if err != nil {
            // Worker is down or unreachable
            fmt.Fprintf(w, "<p>Background Worker Status: **WORKER NOT RUNNING**</p>")
            fmt.Fprintf(w, "<p>Error Detail: %v</p>", err)
        } else {
            // Worker is running and returned a status
            fmt.Fprintf(w, "<p>Background Worker Status: **%s**</p>", workerStatus)
        }
    })
    
    log.Fatal(http.ListenAndServe(":3001", nil)) 
}

// Function to fetch worker status (Simplified version)
func fetchWorkerStatus() (string, error) {
    // Use a short timeout to prevent the web server from hanging
    client := http.Client{Timeout: 500 * time.Millisecond} 
    
    resp, err := client.Get(workerURL)
    if err != nil {
        // Returns error if connection fails (worker is down/not started)
        return "", fmt.Errorf("connection failed: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        // Returns error if worker is running but returns a bad status code (e.g., 404)
        return "", fmt.Errorf("worker returned status code: %d", resp.StatusCode)
    }
    
    // Read the single line of data from the worker
    statusBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("failed to read response body: %w", err)
    }
    
    return string(statusBytes), nil
}