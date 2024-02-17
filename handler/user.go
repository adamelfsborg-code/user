package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adamelfsborg-code/food/user/data"
	"github.com/google/uuid"
)

type UserHandler struct {
	Data data.DataConn
}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		http.Error(w, "Failed to decode json", http.StatusBadRequest)
		return
	}

	user, err := data.NewUserDto(body.Name, body.Password)
	if err != nil {
		fmt.Println("Failed to create user: ", err)
		http.Error(w, "Failed to create user", http.StatusBadRequest)
		return
	}

	err = u.Data.Register(*user)
	if err != nil {
		fmt.Println("Failed to create user: ", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Registerd"))
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		http.Error(w, "Failed to decode json", http.StatusBadRequest)
		return
	}

	token, err := u.Data.Login(body.Name, body.Password)
	if err != nil {
		fmt.Println("Failed to login user: ", err)
		http.Error(w, fmt.Sprintf("Failed to login user: %v", err), http.StatusBadRequest)
		return
	}

	response := map[string]string{"token": token}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Failed to encode token: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (u *UserHandler) Ping(w http.ResponseWriter, r *http.Request) {
	headerId := r.Header.Get("X-USER-ID")

	userId, err := uuid.Parse(headerId)
	if err != nil {
		fmt.Println("Failed to parse id: ", err)
		http.Error(w, "Failed to parse id", http.StatusBadRequest)
		return
	}

	user, err := u.Data.Ping(userId)
	if err != nil {
		fmt.Println("Failed to extract user details: ", err)
		http.Error(w, "Failed to extract user details", http.StatusBadRequest)
		return
	}

	jsonResponse, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Failed to encode user: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
