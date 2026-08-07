package handler

import (
	"net/http"

	"github.com/ballinwza/scraper-dashboard-be/config"
	usecase_user "github.com/ballinwza/scraper-dashboard-be/internal/usecase/user"
	"github.com/ballinwza/scraper-dashboard-be/pkg/validator"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	authUsecase usecase_user.UserUsecase
	cfg         config.Config
}

func NewUserHandler(authUsecase usecase_user.UserUsecase, cfg config.Config) *UserHandler {

	return &UserHandler{
		authUsecase: authUsecase,
		cfg:         cfg,
	}
}

// Register
// @Summary Register a new user
// @Description Register a new user with username, password, and name
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body registerRequest true "Register Request"
// @Success 201 {object} registerResponse
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req registerRequest
	var validator = validator.New()

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if errs := validator.ValidateStruct(&req); len(errs) > 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	user, err := h.authUsecase.Register(
		c.Request.Context(),
		req.Username, req.Password, req.Name,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := registerResponse{
		Username: user.Username,
		Name:     user.Name,
		Role:     user.Role,
		Message:  "user registered successfully",
	}

	c.JSON(http.StatusCreated, res)
}

// Login
// @Summary Loging into system
// @Description Login with username, password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body loginRequest true "Login Request"
// @Success 201 {object} loginResponse
// @Failure 400 {object} map[string]string
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req loginRequest
	var validator = validator.New()

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if errs := validator.ValidateStruct(&req); len(errs) > 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	userAgent := c.Request.UserAgent()
	clientIP := c.ClientIP()

	access, refresh, err := h.authUsecase.Login(
		c.Request.Context(),
		req.Username, req.Password, userAgent, clientIP,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := loginResponse{
		AccessToken:  *access,
		RefreshToken: *refresh,
	}

	c.JSON(http.StatusCreated, res)
}

// Refresh Token
// @Summary Refresh all token of user
// @Description Receive refresh_token from cookie to refreshing refreshToken in data
// @Tags Auth
// @Accept json
// @Produce json
// @Security     CookieAuth
// @Param        refresh_token  header    string           false  "Cookie: refresh_token=..."
// @Success      200            {object}  refreshResponse  "Access Token succesfully"
// @Failure      400            {object}  errorResponse    "Bad input request"
// @Failure      401            {object}  errorResponse    "Cookie missing or Token Invalid"
// @Failure      500            {object}  errorResponse    "Internal Server Error"
// @Header       200            {string}  Set-Cookie       "refresh_token=...; Path=/; HttpOnly; Secure; SameSite=Strict"
// @Router /auth/refresh [post]
func (h *UserHandler) Refresh(c *gin.Context) {
	// var req RefreshClaimsRequest
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
	// 	return
	// }

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token cookie missing"})
		return
	}

	access, refresh, err := h.authUsecase.Refresh(
		c.Request.Context(),
		refreshToken,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	res := loginResponse{
		AccessToken:  *access,
		RefreshToken: *refresh,
	}

	c.JSON(http.StatusOK, res)
}

// @Router /auth/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {}
