package handler

import (
	"net/http"

	"github.com/ballinwza/scraper-dashboard-be/config"
	"github.com/ballinwza/scraper-dashboard-be/internal/delivery/http/dto"
	usecase_ai_estate_rag "github.com/ballinwza/scraper-dashboard-be/internal/usecase/ai_estate_rag"
	"github.com/gin-gonic/gin"
)

type RagMultiTenantHandler struct {
	usecase usecase_ai_estate_rag.IAiEstateRagUsecase
}

func NewRagMultiTenantHandler(usecase usecase_ai_estate_rag.IAiEstateRagUsecase) *RagMultiTenantHandler {
	return &RagMultiTenantHandler{
		usecase: usecase,
	}
}

// SearchSimilar godoc
// @Summary      Search Similar Knowledge (RAG)
// @Description  Search similar content from knowledge base using vector similarity search
// @Tags         RAG
// @Accept       json
// @Produce      json
// @Param        request  body      dto.SearchSimilarReqDTO  true  "Search Similar Request"
// @Success      200      {object}  dto.SearchSimilarResDTO
// @Failure      400      {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Router       /rag/search-similar [post]
func (h *RagMultiTenantHandler) SearchSimilar(c *gin.Context) {
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

	var req dto.SearchSimilarReqDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	var topK *int
	if req.TopK != nil {
		k := int(*req.TopK)
		topK = &k
	}

	res, err := h.usecase.SearchSimilar(
		c.Request.Context(),
		userId.(string),
		req.ChatbotID,
		req.QueryText,
		topK,
		req.KnowledgeFileID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
