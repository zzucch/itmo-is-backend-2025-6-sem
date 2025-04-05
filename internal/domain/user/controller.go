package user

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type Controller struct {
	service        Service
	loginTemplate  *template.Template
	signupTemplate *template.Template
}

func NewController(service Service, loginTmpl, signupTmpl *template.Template) *Controller {
	return &Controller{
		service:        service,
		loginTemplate:  loginTmpl,
		signupTemplate: signupTmpl,
	}
}

func (c *Controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	var inputData struct {
		Username string
		Email    string
		Password string
	}

	if err := json.NewDecoder(r.Body).Decode(&inputData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := HashPassword(inputData.Password)
	if err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	user := User{
		Username:     inputData.Username,
		Email:        inputData.Email,
		PasswordHash: hashedPassword,
	}

	if err := c.service.CreateUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(inputData)
}

func (c *Controller) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := c.service.GetUserByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (c *Controller) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.service.UpdateUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (c *Controller) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := c.service.DeleteUser(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := c.service.Login(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := c.service.GenerateJWT(user)
	if err != nil {
		http.Error(w, "Failed to generate token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(tokenExpiry),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		// TODO do i need it?
		"user":  user,
		"token": token,
	})
}

func (c *Controller) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "already logged out", http.StatusTeapot)
		return
	}

	if err := c.service.InvalidateToken(cookie.Value); err != nil {
		http.Error(w, "failed to logout:", http.StatusInternalServerError)

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})
}

func (c *Controller) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "authorization required", http.StatusUnauthorized)
				return
			}

			if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}
			cookie = &http.Cookie{Value: authHeader[7:]}
		}

		claims, err := c.service.ValidateToken(cookie.Value)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		next(w, r.WithContext(ctx))
	}
}

func (c *Controller) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := c.service.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	json.NewEncoder(w).Encode(response)
}

func (c *Controller) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var updateData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := c.service.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if updateData.Username != "" {
		user.Username = updateData.Username
	}
	if updateData.Email != "" {
		user.Email = updateData.Email
	}
	if updateData.Password != "" {
		hashedPassword, err := HashPassword(updateData.Password)
		if err != nil {
			http.Error(w, "Failed to update password", http.StatusInternalServerError)
			return
		}
		user.PasswordHash = hashedPassword
	}

	if err := c.service.UpdateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (c *Controller) DeleteCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := c.service.DeleteUser(userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if cookie, err := r.Cookie("token"); err == nil && cookie.Value != "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := c.loginTemplate.ExecuteTemplate(w, "layout", nil); err != nil {
		http.Error(w, "Failed to render login page", http.StatusInternalServerError)
	}
}

func (c *Controller) HandleSignupPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if cookie, err := r.Cookie("token"); err == nil && cookie.Value != "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := c.signupTemplate.ExecuteTemplate(w, "layout", nil); err != nil {
		http.Error(w, "Failed to render signup page", http.StatusInternalServerError)
	}
}
