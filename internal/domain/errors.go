package domain

import "errors"

// #Error handler for Usecase & HTTP Response

// ==========================================
// MongoDB connect
// ==========================================
var (
	ErrInvalidID    = errors.New("invalid entity id format")
	ErrInvalidInput = errors.New("invalid input")
	ErrDatabase     = errors.New("database operation failed")
)

// ==========================================
// Common / General Domain Errors
// ==========================================
var (
	ErrInternalServer = errors.New("internal server error occurred")
	ErrBadRequest     = errors.New("bad request parameters")
	ErrNotFound       = errors.New("requested resource not found")
	ErrUnauthorized   = errors.New("unauthorized access")
	ErrForbidden      = errors.New("access forbidden")
	ErrConflict       = errors.New("resource conflict or already exists")
)

// ==========================================
// Scraped Data & Scraper Job Errors
// ==========================================
var (
	ErrItemNotFound      = errors.New("scraped item not found")
	ErrDuplicateItem     = errors.New("scraped item with this URL already exists")
	ErrTargetNotFound    = errors.New("scraper target site configuration not found")
	ErrJobNotFound       = errors.New("scraper job not found")
	ErrJobAlreadyRunning = errors.New("scraper job is currently running")
	ErrJobFailedToStart  = errors.New("failed to initialize scraper job")
)

// ==========================================
// User & Authentication Errors
// ==========================================
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("email or username is already taken")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid or expired authentication token")
)

// ==========================================
// Custom Errors
// ==========================================
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}
