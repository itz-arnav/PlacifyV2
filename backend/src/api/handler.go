package api

import (
	"encoding/json"
	"net/http"
	"placify/backend/src/storage"
	"placify/backend/src/validate"

	"github.com/gorilla/mux"
)

type ApiResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func GreetHandler(w http.ResponseWriter, r *http.Request, userStorage *storage.UserStorage, userValidator *validate.UserValidator) {
	respondWithJSON(w, http.StatusOK, ApiResponse{Message: "Hello, welcome to Placify V2!"})
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request, userStorage *storage.UserStorage, userValidator *validate.UserValidator) {
	var user storage.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}
	if err := userValidator.ValidateUser(&user, false); err != nil {
		respondWithError(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return
	}
	if err := userStorage.CreateUser(&user); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user: "+err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, ApiResponse{Message: "User created successfully", Data: user})
}

func GetUserHandler(w http.ResponseWriter, r *http.Request, userStorage *storage.UserStorage) {
	vars := mux.Vars(r)
	userID := vars["id"]
	user, err := userStorage.GetUser(userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found: "+err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, ApiResponse{Message: "User fetched successfully", Data: user})
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request, userStorage *storage.UserStorage, userValidator *validate.UserValidator) {
	vars := mux.Vars(r)
	userID := vars["id"]
	var user storage.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}
	user.ID = userID
	if err := userValidator.ValidateUser(&user, true); err != nil {
		respondWithError(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return
	}
	if err := userStorage.UpdateUser(userID, &user); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update user: "+err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, ApiResponse{Message: "User updated successfully", Data: user})
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request, userStorage *storage.UserStorage) {
	vars := mux.Vars(r)
	userID := vars["id"]
	if err := userStorage.DeleteUser(userID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete user: "+err.Error())
		return
	}
	respondWithJSON(w, http.StatusNoContent, ApiResponse{Message: "User deleted successfully"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, ApiResponse{Message: message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload ApiResponse) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
