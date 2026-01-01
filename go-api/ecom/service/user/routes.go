package user

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/manaswini1869/interview-prep/go-api/ecom/service/auth"
	"github.com/manaswini1869/interview-prep/go-api/ecom/types"
	"github.com/manaswini1869/interview-prep/go-api/ecom/utils"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {

	// get json payload from request body
	// check if user already exists
	// if not create user in db
	var payload types.RegisterUserPayload
	err := utils.ParseJSON(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	if payload.Email == "" || payload.Password == "" || payload.FirstName == "" || payload.LastName == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid Payload"))
		return
	}

	_, err = h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("User Already Exists"))
		return
	}
	hashedPassword := auth.HashedPassword(payload.Password)
	if hashedPassword == "" {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Error adding the data to db"))
		return
	}
	err = h.store.CreateUser(&types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "User Created Successfully"})
}
