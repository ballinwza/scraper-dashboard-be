package handler

import (
	"net/http"

	"github.com/ballinwza/scraper-dashboard-be/config"
	"github.com/ballinwza/scraper-dashboard-be/internal/delivery/http/dto"
	domain_ai_estate_rag "github.com/ballinwza/scraper-dashboard-be/internal/domain/ai_estate_rag"
	usecase_ai_estate_rag "github.com/ballinwza/scraper-dashboard-be/internal/usecase/ai_estate_rag"
	"github.com/gin-gonic/gin"
)

type ChatSessionHandler struct {
	usecase usecase_ai_estate_rag.IAiEstateRagUsecase
}

func NewChatSessionHandler(usecase usecase_ai_estate_rag.IAiEstateRagUsecase) *ChatSessionHandler {
	return &ChatSessionHandler{
		usecase: usecase,
	}
}

// CreateChatSession godoc
// @Summary      Create Chat Session
// @Description  Create a new chat session for a specific chatbot
// @Tags         Chat Session
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateChatSessionReqDTO  true  "Create Chat Session Request"
// @Success      201      {object}  dto.CreateChatSessionResDTO
// @Failure      400      {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Router       /chat-sessions/create [post]
func (h *ChatSessionHandler) CreateChatSession(c *gin.Context) {
	userId, exists := c.Get(config.USER_ID_KEY)
	if !exists {
		c.JSON(http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
		return
	}
	userId, ok := userId.(string)
	if !ok || userId == "" {
		c.JSON(http.StatusUnauthorized, errorResponse{Error: "invalid user session"})
		return
	}

	var req dto.CreateChatSessionReqDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	res, err := h.usecase.CreateChatSession(
		c.Request.Context(),
		userId.(string),
		req.ChatbotID,
		req.SessionTitle,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

// GetChatSession godoc
// @Summary      Get Chat Session by ID
// @Description  Get specific chat session details and message history
// @Tags         Chat Session
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "Session ID"
// @Success      200      {object}  dto.GetChatSessionResDTO
// @Failure      400      {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Router       /chat-sessions/{id} [get]
func (h *ChatSessionHandler) GetChatSession(c *gin.Context) {
	userId, exists := c.Get(config.USER_ID_KEY)
	if !exists {
		c.JSON(http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
		return
	}
	userId, ok := userId.(string)
	if !ok || userId == "" {
		c.JSON(http.StatusUnauthorized, errorResponse{Error: "invalid user session"})
		return
	}

	var req dto.GetChatSessionReqDTO

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	res, err := h.usecase.GetChatSession(c.Request.Context(), req.ID, userId.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// ListChatSessions godoc
// @Summary      List Chat Sessions
// @Description  Get a paginated list of chat sessions for a user or chatbot
// @Tags         Chat Session
// @Accept       json
// @Produce      json
// @Param 		 request 	 body 	   dto.ListChatSessionsReqDTO 	true "Request"
// @Success      200         {object}  dto.ListChatSessionsResDTO
// @Failure      400         {object}  errorResponse
// @Failure      500         {object}  errorResponse
// @Router       /chat-sessions/list [post]
func (h *ChatSessionHandler) ListChatSessions(c *gin.Context) {
	userId, exists := c.Get(config.USER_ID_KEY)
	if !exists {
		c.JSON(http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
		return
	}
	userId, ok := userId.(string)
	if !ok || userId == "" {
		c.JSON(http.StatusUnauthorized, errorResponse{Error: "invalid user session"})
		return
	}

	var req dto.ListChatSessionsReqDTO

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	res, err := h.usecase.ListChatSessions(
		c.Request.Context(),
		userId.(string),
		req.ChatbotID,
		int(req.PageSize),
		int(req.PageToken),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// AddChatMessage godoc
// @Summary      Add Chat Message
// @Description  Add a new chat message to an existing session
// @Tags         Chat Session
// @Accept       json
// @Produce      json
// @Param        request  body      dto.AddChatMessageReqDTO  true  "Add Chat Message Request"
// @Success      201      {object}  dto.AddChatMessageResDTO
// @Failure      400      {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Router       /chat-sessions/messages [post]
func (h *ChatSessionHandler) AddChatMessage(c *gin.Context) {
	userId, exists := c.Get(config.USER_ID_KEY)
	if !exists {
		c.JSON(http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
		return
	}
	userId, ok := userId.(string)
	if !ok || userId == "" {
		c.JSON(http.StatusUnauthorized, errorResponse{Error: "invalid user session"})
		return
	}

	var req dto.AddChatMessageReqDTO

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	res, err := h.usecase.AddChatMessage(
		c.Request.Context(),
		req.SessionID,
		userId.(string),
		req.Content,
		domain_ai_estate_rag.ParseMessageRole(req.Role.String()),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

// DeleteChatSession godoc
// @Summary      Delete Chat Session
// @Description  Delete a chat session by ID
// @Tags         Chat Session
// @Accept       json
// @Produce      json
// @Param 		 request 	 body 	   dto.DeleteChatSessionReqDTO 	true "Request"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Router       /chat-sessions/delete [delete]
func (h *ChatSessionHandler) DeleteChatSession(c *gin.Context) {
	userId, exists := c.Get(config.USER_ID_KEY)
	if !exists {
		c.JSON(http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
		return
	}
	userId, ok := userId.(string)
	if !ok || userId == "" {
		c.JSON(http.StatusUnauthorized, errorResponse{Error: "invalid user session"})
		return
	}

	var req dto.DeleteChatSessionReqDTO

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	if err := h.usecase.DeleteChatSession(c.Request.Context(), req.ID, userId.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Chat session deleted successfully",
	})
}
