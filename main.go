package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

// MirrorListResponse represents the structure of the JSON response for the mirrorlist API.
type MirrorListResponse struct {
	Version string   `json:"version"`
	Mirrors []string `json:"mirrors"`
}

// LoadMirrors reads the mirrors from a JSON file and returns them as a slice of strings.
func LoadMirrors() ([]string, error) {
	// Construct the file path for mirrors.json
	filePath := filepath.Join("data", "mirrors.json")

	// Read the JSON file
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading mirrors file %s: %v", filePath, err)
		return nil, fmt.Errorf("failed to read mirrors file: %v", err)
	}

	// Create a struct to unmarshal the JSON data
	var mirrorsData struct {
		Mirrors []string `json:"mirrors"`
	}

	// Parse the JSON file
	err = json.Unmarshal(data, &mirrorsData)
	if err != nil {
		log.Printf("Error parsing mirrors JSON file: %v", err)
		return nil, fmt.Errorf("failed to parse mirrors JSON: %v", err)
	}

	return mirrorsData.Mirrors, nil
}

func main() {
	// Create a new router using Gorilla Mux
	r := mux.NewRouter()

	// Serve index.html on the root path "/"
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	// Serve all static files from the "static" directory
	staticFileDirectory := http.Dir("./static")
	staticFileHandler := http.StripPrefix("/static/", http.FileServer(staticFileDirectory))
	r.PathPrefix("/static/").Handler(staticFileHandler)

	// API endpoint for /api/mirrorlist/:version/json
	r.HandleFunc("/api/mirrorlist/{version}/json", MirrorListHandler).Methods("GET")

	// Start the server
	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// MirrorListHandler handles the /api/mirrorlist/:version/json endpoint
func MirrorListHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	version := vars["version"]

	// Load mirrors from the JSON file
	mirrors, err := LoadMirrors()
	if err != nil {
		log.Printf("Error loading mirrors: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create the response
	response := MirrorListResponse{
		Version: version,
		Mirrors: mirrors,
	}

	// Convert the response to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error generating JSON response: %v", err)
		http.Error(w, "Error generating JSON", http.StatusInternalServerError)
		return
	}

	// Set content-type to application/json and write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Printf("Error writing JSON response: %v", err)
	}
}
