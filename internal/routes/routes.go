package routes

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {

	// Health Checks
	mux.HandleFunc("GET /api/health-check-http", handleHealthCheck)
}
