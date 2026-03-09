package routes

import (
	"net/http"

	"github.com/DEV-BC/backend_chatapp/internal/utils"
)

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.JSON(w, http.StatusOK, true, "Api is running", nil)
}
