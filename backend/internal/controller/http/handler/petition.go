package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"mado/internal/core/petition"
	"mado/pkg/errs"
)

type PetitionService interface {
	GetPetitionPdfByID(doc_id string) (response *petition.PetitionData, err error)
	GeneratePetitionPDF(data *petition.PetitionData) (*string, error)
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

	// Call the GeneratePetitionPDF function
	generatedData, err := h.petitionService.GeneratePetitionPDF(&requestData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set the Content-Type header to specify UTF-8 encoding
	c.Header("Content-Type", "application/json; charset=utf-8")

	// Return the generated PDF data in the response
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename="+*generatedData)
	// filePath := "internal/core/petition/outputs/" + *generatedData
	os.Chdir("../../")
	fmt.Println(os.Getwd())
	c.File(*generatedData)
}
