package routes

import "net/http"

func RegisterRoutes() *http.ServeMux {

	mux := http.NewServeMux()
	// Health Checks
	mux.HandleFunc("GET /api/health-check-http", handleHealthCheck)

	//Auths
	mux.HandleFunc("POST /api/auth/register-email", handleEmailRegister)

	return mux
}
