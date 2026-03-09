package routes

import (
	"encoding/json"
	"net/http"

	"github.com/DEV-BC/backend_chatapp/internal/models"
	"github.com/DEV-BC/backend_chatapp/internal/utils"
)

func handleEmailRegister(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid request body", nil)
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid credentials", nil)
		return
	}

	existingUser, _ := models.GetUserByEmail(req.Email)
	if existingUser != nil {
		utils.JSON(w, http.StatusConflict, false, "Email already in use", nil)
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Signup failed, please try again later", nil)
		return
	}

	user, err := models.CreateUserByEmail(req.Name, req.Email, hashedPassword)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Could not register", nil)
		return
	}

	utils.JSON(w, http.StatusCreated, true, "Sign up successful", user)
}
