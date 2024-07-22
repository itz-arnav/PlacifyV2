package api

import (
	"net/http"
	"placify/backend/src/storage"
	"placify/backend/src/validate"

	"github.com/gorilla/mux"
)

func InitializeRouter(userStorage *storage.UserStorage, userValidator *validate.UserValidator) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/greet", func(w http.ResponseWriter, r *http.Request) {
		GreetHandler(w, r, userStorage, userValidator)
	}).Methods("POST")
	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		CreateUserHandler(w, r, userStorage, userValidator)
	}).Methods("POST")
	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		GetUserHandler(w, r, userStorage)
	}).Methods("GET")
	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		UpdateUserHandler(w, r, userStorage, userValidator)
	}).Methods("PUT")
	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		DeleteUserHandler(w, r, userStorage)
	}).Methods("DELETE")
	return router
}
