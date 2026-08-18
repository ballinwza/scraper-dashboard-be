package handler

import (
	"net/http"

	"github.com/ballinwza/scraper-dashboard-be/config"
	"github.com/ballinwza/scraper-dashboard-be/internal/delivery/http/dto"
	usecase_ai_estate_rag "github.com/ballinwza/scraper-dashboard-be/internal/usecase/ai_estate_rag"
	"github.com/gin-gonic/gin"
)

type ChatbotHandler struct {
	usecase usecase_ai_estate_rag.IAiEstateRagUsecase
}

func NewChatbotHandler(usecase usecase_ai_estate_rag.IAiEstateRagUsecase) *ChatbotHandler {
	return &ChatbotHandler{
		usecase: usecase,
	}
}

// CreateMultiTenantChatbot godoc
// @Summary      Create Multi-Tenant Chatbot
// @Description  Create a new chatbot blueprint for a specific user/tenant
// @Tags         Chatbot
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateMultiTenantChatbotReqDTO  true  "Create Chatbot Request"
// @Success      201      {object}  dto.CreateMultiTenantChatbotResDTO
// @Failure      400      {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Router       /chatbots/create [post]
func (h *ChatbotHandler) CreateMultiTenantChatbot(c *gin.Context) {
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

	var req dto.CreateMultiTenantChatbotReqDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	res, err := h.usecase.CreateMultiTenantChatbot(
		c.Request.Context(),
		userId.(string),
		req.Name,
		req.Description,
		req.SystemPrompt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

// GetMultiTenantChatbot godoc
// @Summary      Get Multi-Tenant Chatbot by ID
// @Description  Get specific chatbot blueprint details by its ID and User ID
// @Tags         Chatbot
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "Chatbot ID"
// @Success      200      {object}  dto.GetMultiTenantChatbotResDTO
// @Failure      400      {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Router       /chatbots/{id} [get]
func (h *ChatbotHandler) GetMultiTenantChatbot(c *gin.Context) {
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

	var req dto.GetMultiTenantChatbotReqDTO

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	res, err := h.usecase.GetMultiTenantChatbot(c.Request.Context(), req.ID, userId.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// ListMultiTenantChatbots godoc
// @Summary      List Multi-Tenant Chatbots
// @Description  Get a paginated list of chatbot blueprints for a specific user
// @Tags         Chatbot
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ListMultiTenantChatbotsReqDTO  true  "Create Chat Session Request"
// @Success      200         {object}  dto.ListMultiTenantChatbotsResDTO
// @Failure      400         {object}  errorResponse
// @Failure      500         {object}  errorResponse
// @Router       /chatbots/list [post]
func (h *ChatbotHandler) ListMultiTenantChatbots(c *gin.Context) {
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

	var req dto.ListMultiTenantChatbotsReqDTO

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	res, err := h.usecase.ListMultiTenantChatbots(
		c.Request.Context(),
		userId.(string),
		int(req.PageSize),
		int(req.PageToken),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// UpdateMultiTenantChatbot godoc
// @Summary      Update Multi-Tenant Chatbot
// @Description  Update details of an existing chatbot blueprint
// @Tags         Chatbot
// @Accept       json
// @Produce      json
// @Param        request  body      dto.UpdateMultiTenantChatbotReqDTO  true  "Update Chatbot Request"
// @Success      200      {object}  dto.UpdateMultiTenantChatbotResDTO
// @Failure      400      {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Router       /chatbots/update [post]
func (h *ChatbotHandler) UpdateMultiTenantChatbot(c *gin.Context) {
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

	var req dto.UpdateMultiTenantChatbotReqDTO

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	var name, description, systemPrompt string
	if req.Name != nil {
		name = *req.Name
	}
	if req.Description != nil {
		description = *req.Description
	}
	if req.SystemPrompt != nil {
		systemPrompt = *req.SystemPrompt
	}

	res, err := h.usecase.UpdateMultiTenantChatbot(
		c.Request.Context(),
		req.ID,
		userId.(string),
		name,
		description,
		systemPrompt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// DeleteMultiTenantChatbot godoc
// @Summary      Delete Multi-Tenant Chatbot
// @Description  Delete a chatbot blueprint by ID and User ID
// @Tags         Chatbot
// @Accept       json
// @Produce      json
// @Param        request  body      dto.DeleteMultiTenantChatbotReqDTO  true  "Update Chatbot Request"
// @Success      200      {object}  dto.DeleteMultiTenantChatbotResDTO
// @Failure      400      {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Router       /chatbots/delete [delete]
func (h *ChatbotHandler) DeleteMultiTenantChatbot(c *gin.Context) {
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

	var req dto.DeleteMultiTenantChatbotReqDTO

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	res, err := h.usecase.DeleteMultiTenantChatbot(c.Request.Context(), req.ID, userId.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
