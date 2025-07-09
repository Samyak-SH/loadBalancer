package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config represents the structure of config.json
type Config struct {
	Port                uint16   `json:"PORT"`
	Servers             []string `json:"Servers"`
	Algorithm           uint16   `json:"Algorithm"`
	VirtualNodeCount    int      `json:"VirtualNodeCount"`
	HealthCheckInterval int      `json:"HealthCheckInterval"`
	SecretKey           string   `json:"SECRET_KEY"`
}

func main() {
	// Read env variables
	portStr := os.Getenv("PORT")
	serversStr := os.Getenv("SERVERS")
	algorithmStr := os.Getenv("ALGORITHM")
	vnCountStr := os.Getenv("VIRTUAL_NODE_COUNT")
	hcIntervalStr := os.Getenv("HEALTH_CHECK_INTERVAL")
	secretKey := os.Getenv("SECRET_KEY")

	// Parse environment variables
	port, err := strconv.ParseUint(portStr, 10, 16)
	checkErr("PORT", err)

	algorithm, err := strconv.ParseUint(algorithmStr, 10, 16)
	checkErr("ALGORITHM", err)

	vnCount, err := strconv.Atoi(vnCountStr)
	checkErr("VIRTUAL_NODE_COUNT", err)

	hcInterval, err := strconv.Atoi(hcIntervalStr)
	checkErr("HEALTH_CHECK_INTERVAL", err)

	// Split servers by comma
	servers := splitAndTrim(serversStr)

	// Create config struct
	config := Config{
		Port:                uint16(port),
		Servers:             servers,
		Algorithm:           uint16(algorithm),
		VirtualNodeCount:    vnCount,
		HealthCheckInterval: hcInterval,
		SecretKey:           secretKey,
	}

	// Write to config.json
	file, err := os.Create("config.json") // writes in /app/config.json
	if err != nil {
		panic(fmt.Sprintf("Error creating config.json: %v", err))
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ") // pretty print
	err = encoder.Encode(config)
	if err != nil {
		panic(fmt.Sprintf("Error writing to config.json: %v", err))
	}

	fmt.Println("config.json generated successfully.")
}

// splitAndTrim turns comma-separated string into slice of trimmed strings
func splitAndTrim(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func checkErr(field string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid value for %s: %v\n", field, err)
		os.Exit(1)
	}
}
