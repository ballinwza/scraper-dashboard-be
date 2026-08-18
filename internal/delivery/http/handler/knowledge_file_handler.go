package handler

import (
	"io"
	"net/http"

	"github.com/ballinwza/scraper-dashboard-be/config"
	"github.com/ballinwza/scraper-dashboard-be/internal/delivery/http/dto"
	usecase_ai_estate_rag "github.com/ballinwza/scraper-dashboard-be/internal/usecase/ai_estate_rag"
	"github.com/gin-gonic/gin"
)

type KnowledgeFileHandler struct {
	usecase usecase_ai_estate_rag.IAiEstateRagUsecase
}

func NewKnowledgeFileHandler(usecase usecase_ai_estate_rag.IAiEstateRagUsecase) *KnowledgeFileHandler {
	return &KnowledgeFileHandler{
		usecase: usecase,
	}
}

// GetKnowledgeFile godoc
// @Summary      Get Knowledge File by ID
// @Description  Get specific knowledge file details by its ID and User ID
// @Tags         Knowledge File
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "File ID"
// @Success      200      {object}  domain_ai_estate_rag.GetKnowledgeFileResponse
// @Failure      400      {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Router       /knowledge-files/{id} [get]
func (h *KnowledgeFileHandler) GetKnowledgeFile(c *gin.Context) {
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

	var req dto.GetKnowledgeFileReqDTO

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	res, err := h.usecase.GetKnowledgeFile(c.Request.Context(), req.ID, userId.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// ListKnowledgeFiles godoc
// @Summary      List Knowledge Files
// @Description  Get a paginated list of knowledge files for a chatbot
// @Tags         Knowledge File
// @Accept       json
// @Produce      json
// @Param 		 request 	 body 	   dto.ListKnowledgeFilesReqDTO 	true "Request"
// @Success      200         {object}  domain_ai_estate_rag.ListKnowledgeFilesResponse
// @Failure      400         {object}  errorResponse
// @Failure      500         {object}  errorResponse
// @Router       /knowledge-files/list [post]
func (h *KnowledgeFileHandler) ListKnowledgeFiles(c *gin.Context) {
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

	var req dto.ListKnowledgeFilesReqDTO

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	res, err := h.usecase.ListKnowledgeFiles(
		c.Request.Context(),
		req.ChatbotID,
		userId.(string),
		req.Limit,
		req.Offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// DeleteKnowledgeFile godoc
// @Summary      Delete Knowledge File
// @Description  Delete a knowledge file associated with a chatbot and user
// @Tags         Knowledge File
// @Accept       json
// @Produce      json
// @Param        request  body      dto.DeleteKnowledgeFileReqDTO  true  "Delete Knowledge File Request"
// @Success      200      {object}  domain_ai_estate_rag.DeleteKnowledgeFileResponse
// @Failure      400      {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Router       /knowledge-files/delete [delete]
func (h *KnowledgeFileHandler) DeleteKnowledgeFile(c *gin.Context) {
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

	var req dto.DeleteKnowledgeFileReqDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.usecase.DeleteKnowledgeFile(c.Request.Context(), req.ChatbotID, userId.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// UploadFileMultiTenant
// @Summary      Upload file (PDF/Image) to Multi-Tenant RAG
// @Description  Uploads a PDF or Image file using streaming gRPC underneath
// @Tags         RAG
// @Accept       multipart/form-data
// @Produce      json
// @Param 		request 	body 	dto.UploadFileMultiTenantReqDTO 	true 	"Upload Request"
// @Param        file        formData  file    true  "File to upload (PDF, PNG, JPG, WEBP)"
// @Success      200         {object}  domain_ai_estate_rag.UploadFileMultiTenantResponse
// @Failure      400         {object}  errorResponse "Invalid payload or file type"
// @Failure      500         {object}  errorResponse "Internal server error"
// @Router       /knowledge-files/upload [post]
func (h *KnowledgeFileHandler) UploadFileMultiTenant(c *gin.Context) {
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

	var reqDto dto.UploadFileMultiTenantReqDTO

	// 1. Bind Multipart Form Data
	if err := c.ShouldBind(&reqDto); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	// 2. Validate Content-Type / File Extension
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "Failed to get file from request: " + err.Error()})
		return
	}

	// 3. อ่านไฟล์เปลี่ยนเป็น []byte
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "Failed to open uploaded file"})
		return
	}
	defer file.Close()

	filename := fileHeader.Filename
	contentType := fileHeader.Header.Get("Content-Type")

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "Failed to read file content"})
		return
	}

	// 5. เรียกใช้ gRPC Stream Repository
	result, err := h.usecase.MultiTenantUploadFile(
		c.Request.Context(),
		userId.(string),
		reqDto.ChatbotId,
		filename,
		contentType,
		fileBytes,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	// 6. Return Response กลับ HTTP Client
	c.JSON(http.StatusOK, result)
}
