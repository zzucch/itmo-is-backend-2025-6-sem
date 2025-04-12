package user

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
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

// @Summary creates a new user account
// @Description creates a new user account
// @Tags users
// @Router /api/users [post]
// @Param request body user.CreateUserRequest true "User data"
// @Success 201 {object} user.UserResponse
func (c *Controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	var inputData CreateUserRequest

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

	user.PasswordHash = ""
	json.NewEncoder(w).Encode(user)
}

// @Summary retrieves user by ID
// @Description retrieves user by ID
// @Tags users
// @Router /api/users/{id} [get]
// @Param id path int true "User ID"
// @Security BearerAuth
// @Success 200 {object} user.UserResponse
// @Failure 400 {object} map[string]string "Invalid user ID"
func (c *Controller) GetUserByID(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	idStr := strings.Split(path, "/")[0]
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

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// @Summary updates user information
// @Description updates user information
// @Tags users
// @Router /api/users/{id} [put]
// @Param id path int true "User ID"
// @Security BearerAuth
// @Param request body user.User true "User update data"
// @Success 200 {object} user.UserResponse
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal server error"
func (c *Controller) UpdateUser(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	idStr := strings.Split(path, "/")[0]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	isAdmin, ok := r.Context().Value("is_admin").(bool)
	if !ok {
		isAdmin = false
	}

	if (userID != user.ID || userID != uint(id)) && !isAdmin {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := c.service.UpdateUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// @Summary retrieves all users
// @Description retrieves all users
// @Tags users
// @Router /api/users [get]
// @Security BearerAuth
// @Success 200 {array} user.UserResponse
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal server error"
func (c *Controller) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if _, ok := r.Context().Value("user_id").(uint); !ok {
		log.Println("meow")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if isAdmin, ok := r.Context().Value("is_admin").(bool); !isAdmin && !ok {
		log.Println("meow")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Println(":3")

	users, err := c.service.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("qwe", users)

	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Println("asdaqwe")
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
	log.Println("qeq")
	w.Header().Set("Content-Type", "application/json")
}

// @Summary removes a user account
// @Description removes a user account
// @Tags users
// @Router /api/users/{id} [delete]
// @Param id path int true "User ID"
// @Security BearerAuth
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
func (c *Controller) DeleteUser(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	idStr := strings.Split(path, "/")[0]
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

// @Summary authenticates user and returns token
// @Description authenticates user and returns token
// @Tags users
// @Router /api/users/login [post]
// @Param request body user.LoginRequest true "Credentials"
// @Success 200 {object} user.LoginResponse
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
		"user":  user,
		"token": token,
	})
}

// @Summary Get user tokens
// @Description Get all active tokens for a user
// @Tags users/me
// @Security ApiKeyAuth
// @Router /api/users/me/tokens [get]
// @Success 200 {array} user.Token
func (c *Controller) GetUserTokens(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tokens, err := c.service.GetUserTokens(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

// @Summary Invalidate all user tokens
// @Description Invalidates all active tokens for a user
// @Tags users/me
// @Security ApiKeyAuth
// @Router /api/users/me/invalidate-tokens [post]
// @Param id path uint true "User ID"
// @Success 200 {string} string "All tokens invalidated successfully"
func (c *Controller) InvalidateUserTokens(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := c.service.InvalidateAllUserTokens(userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All tokens invalidated successfully"))
}

// @Summary ends user session
// @Description ends user session
// @Tags users
// @Router /api/users/logout [post]
// @Security BearerAuth
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
		log.Println("authin")
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
			http.Error(
				w,
				"auth middleware: validation failed: "+err.Error(),
				http.StatusUnauthorized,
			)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		log.Println("daa")
		next(w, r.WithContext(ctx))
	}
}

// @Summary returns logged-in user's profile
// @Description returns logged-in user's profile
// @Tags users/me
// @Router /api/users/me [get]
// @Security BearerAuth
// @Success 200 {object} user.UserResponse
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

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// @Summary modifies user profile
// @Description modifies user profile
// @Tags users/me
// @Router /api/users/me [put]
// @Security BearerAuth
// @Param request body user.UpdateUserRequest false "Update data"
// @Success 200 {object} user.UpdateResponse
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

// @Summary removes the currently authenticated user
// @Description removes the currently authenticated user
// @Tags users/me
// @Router /api/users/me [delete]
// @Security BearerAuth
// @Success 204 "No Content"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
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

// @Summary Renders login page
// @Description Renders login page
// @Tags pages
// @Router /login [get]
// @Produce html
// @Success 200 "HTML login page"
// @Failure 302 "Redirect to home if already logged in"
// @Failure 405 "Method not allowed"
// @Failure 500 "Internal server error"
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

// @Summary Renders signup page
// @Description Renders signup page
// @Tags pages
// @Router /signup [get]
// @Produce html
// @Success 200 "HTML signup page"
// @Failure 302 "Redirect to home if already logged in"
// @Failure 405 "Method not allowed"
// @Failure 500 "Internal server error"
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

func (c *Controller) AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("adminin")
		userID, ok := r.Context().Value("user_id").(uint)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := c.service.GetUserByID(userID)
		if err != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		if user.Role != RoleAdmin {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "is_admin", true)
		log.Println("daada")
		next(w, r.WithContext(ctx))
	}
}
