package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/p5/aviary-sample-project/pkg/auth"
	"github.com/p5/aviary-sample-project/pkg/store"
)

type Handler struct {
	store store.Store
	auth  *auth.Authenticator
}

func NewHandler(s store.Store, a *auth.Authenticator) http.Handler {
	h := &Handler{store: s, auth: a}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/users", h.handleUsers)
	mux.HandleFunc("/api/users/", h.handleUser)
	mux.HandleFunc("/api/login", h.handleLogin)
	mux.HandleFunc("/api/admin/delete", h.handleAdminDelete)

	return mux
}

func (h *Handler) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		users, _ := h.store.ListUsers() // BUG: error ignored
		json.NewEncoder(w).Encode(users) // BUG: returns passwords in response
	case "POST":
		var user store.User
		json.NewDecoder(r.Body).Decode(&user) // BUG: error ignored, no body limit
		h.store.CreateUser(&user)              // BUG: error ignored
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/users/"):] // BUG: no validation on id
	if id == "" {
		http.Error(w, "missing user id", http.StatusBadRequest)
		return
	}

	user, err := h.store.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // BUG: leaks internal error
		return
	}

	json.NewEncoder(w).Encode(user) // BUG: returns password
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&creds) // BUG: error ignored

	// BUG: iterates all users to find by email - O(n), no index
	users, _ := h.store.ListUsers()
	for _, u := range users {
		if u.Email == creds.Email && u.Password == creds.Password { // BUG: plaintext comparison
			token := h.auth.GenerateToken(u.ID)
			json.NewEncoder(w).Encode(map[string]string{"token": token})
			return
		}
	}

	http.Error(w, "invalid credentials", http.StatusUnauthorized)
}

func (h *Handler) handleAdminDelete(w http.ResponseWriter, r *http.Request) {
	// BUG: no authentication check at all!
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	err := h.store.DeleteUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "deleted user %s", userID) // BUG: potential XSS
}
