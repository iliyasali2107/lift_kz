package handler

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"mado/internal/core/petition"
	"mado/pkg/errs"
)

type PetitionService interface {
	GetPetitionPdfByID(doc_id string) (response *petition.PetitionData, err error)
	GeneratePetitionPDF(data *petition.PetitionData) (*string, error)
	GeneratePDF(ctx context.Context, t *template.Template, pageData interface{}, outFilePath string, templatePath string) error
	GenerateFinalPdf(filePath string) (string, error)
}

type petitionDeps struct {
	router          *gin.RouterGroup
	petitionService PetitionService
}

type petitionHandler struct {
	petitionService PetitionService
}

func newPetitionHandler(deps petitionDeps) {
	handler := petitionHandler{
		petitionService: deps.petitionService,
	}

	usersGroup := deps.router.Group("/petition")
	{
		usersGroup.GET("/download/:id", handler.GetPetitionPDF)
		usersGroup.POST("/generate", handler.GeneratePetitionPDFHandler)
	}

}

func (h petitionHandler) GetPetitionPDF(c *gin.Context) {
	response, err := h.petitionService.GetPetitionPdfByID(c.Param("id"))
	if err != nil {
		if errors.Is(err, errs.ErrInvalidID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		} else if errors.Is(err, errs.ErrPdfFileNotFound) {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	// Set the appropriate headers for file download
	c.Header("Content-Disposition", strings.Join([]string{"attachment; filename", response.FileName + ".pdf"}, "="))
	// c.Header("Content-Disposition", "attachment; filename=your-file-name.pdf")
	c.Header("Content-Type", "application/pdf")
	c.Data(http.StatusOK, "application/pdf", response.PdfData)
}

func (h petitionHandler) GeneratePetitionPDFHandler(c *gin.Context) {

	// Parse JSON request body into a PetitionData struct
	var requestData petition.PetitionData
	if err := c.BindJSON(&requestData); err != nil {
		fmt.Println("HERE")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentDir, _ := os.Getwd()

	outFileName := uuid.New().String() + ".pdf"
	htmlPath := currentDir + "/files/input_html/index.html"
	outPath := currentDir + "/files/output_pdf/" + outFileName
	temp, err := template.ParseFiles(htmlPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.petitionService.GeneratePDF(c, temp, requestData, outPath, htmlPath)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	finalFilePath, err := h.petitionService.GenerateFinalPdf(outPath)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Content-Type", "application/pdf")

	c.File(finalFilePath)
}
