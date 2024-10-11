package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// TestAPIEndpoint tests the /api/mirrorlist/:version/json endpoint.
func TestAPIEndpoint(t *testing.T) {
    // Create a request to pass to the handler.
    req, err := http.NewRequest("GET", "/api/mirrorlist/v1/json", nil)
    if err != nil {
        t.Fatalf("Could not create request: %v", err)
    }

    // Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
    rr := httptest.NewRecorder()

    // Create a new router using Gorilla Mux.
    r := setupRouter()

    // Serve the request
    r.ServeHTTP(rr, req)

    // Check if the status code is 200 OK.
    if rr.Code != http.StatusOK {
        t.Errorf("Expected status code 200, got %d", rr.Code)
    }

    // Check if the response is valid JSON.
    var response MirrorListResponse
    if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
        t.Errorf("Response was not valid JSON: %v", err)
    }

    // Check if the mirrors field has data.
    if len(response.Mirrors) == 0 {
        t.Errorf("Expected some mirrors, got none")
    }
}

// TestRootPath tests the root / path.
func TestRootPath(t *testing.T) {
    // Create a request to pass to the handler.
    req, err := http.NewRequest("GET", "/", nil)
    if err != nil {
        t.Fatalf("Could not create request: %v", err)
    }

    // Create a ResponseRecorder to record the response.
    rr := httptest.NewRecorder()

    // Create a new router using Gorilla Mux.
    r := setupRouter()

    // Serve the request
    r.ServeHTTP(rr, req)

    // Check if the status code is 200 OK.
    if rr.Code != http.StatusOK {
        t.Errorf("Expected status code 200, got %d", rr.Code)
    }

    // Check if the response body contains the content of the index.html file or any static file.
    body, err := ioutil.ReadAll(rr.Body)
    if err != nil {
        t.Fatalf("Could not read response body: %v", err)
    }

    expectedSubstring := "<html" // You can change this to something specific in your index.html.
    if !contains(body, expectedSubstring) {
        t.Errorf("Expected response body to contain '%s', got '%s'", expectedSubstring, string(body))
    }
}

// Helper function to check if the response body contains a string.
func contains(body []byte, substring string) bool {
    return string(body) != "" && substring != "" && string(body)[:len(substring)] == substring
}

// setupRouter sets up the router for testing.
func setupRouter() *mux.Router {
    r := mux.NewRouter()

    // Serve static files from the "static" directory
    fileServer := http.FileServer(http.Dir("./static"))
    r.PathPrefix("/").Handler(fileServer)

    // API endpoint for /api/mirrorlist/:version/json
    r.HandleFunc("/api/mirrorlist/{version}/json", MirrorListHandler).Methods("GET")

    return r
}
