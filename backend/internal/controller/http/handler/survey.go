package handler

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"

	// "strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"mado/helpers"
	"mado/internal/core/petition"
	"mado/internal/core/survey"
	"mado/models"
)

type SurveyService interface {
	Create(*survey.SurveyRequirements) (*survey.SurveyRequirements, error)
	GetSurviesByUserID(user_id int, ctx *gin.Context) (response []*survey.SurveyResponse, err error)
	CloseSurvey(ctx context.Context, survey_id int) error
	GetSurveyById(ctx context.Context, surveyId int) (models.Survey, error)
	GetSurveySummary(ctx context.Context, survey_id int) (models.Survey, error)
	SaveSurvey(ctx context.Context, req models.SaveSurveyRequest) error
	GetSurveyByUserIdAndSurveyId(ctx context.Context, surveyId, userId int) (models.Survey, error)
}

type UserReq struct {
	UserID int `json:"user_id"`
}

type surveyDeps struct {
	router *gin.RouterGroup

	petitionService PetitionService
	surveyService   SurveyService
	userService     UserService
}

type surveyHandler struct {
	surveyService   SurveyService
	petitionService PetitionService
	userService     UserService
}

func newSurveyHandler(deps surveyDeps) {
	handler := surveyHandler{
		surveyService:   deps.surveyService,
		petitionService: deps.petitionService,
		userService:     deps.userService,
	}

	usersGroup := deps.router.Group("/survey")
	{
		usersGroup.POST("/create", handler.CreateSurvey)
		usersGroup.GET("/get/surveys/:id", handler.GetSurveis)
		usersGroup.POST("/close", handler.CloseSurvey)
		usersGroup.GET("/get/survey/:surveyId", handler.GetSurvey)
		usersGroup.GET("/summary/:survey_id", handler.GetSurveySummary)
		usersGroup.POST("/save", handler.SaveSurvey)
	}

}

type GetSurveySummaryRequest struct {
	SurveyId int `json:"survey_id"`
}

func (h surveyHandler) GetSurveySummary(c *gin.Context) {
	// var surveyId GetSurveySummaryRequest
	survey_id, err := strconv.Atoi(c.Param("survey_id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	surveys, err := h.surveyService.GetSurveySummary(c, survey_id)
	if err != nil {
		if err == survey.ErrNotFound {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, surveys)
}

func (h surveyHandler) CreateSurvey(c *gin.Context) {
	var request *survey.SurveyRequirements
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	resp, err := h.surveyService.Create(request)
	if err != nil {
		switch err {
		case survey.ErrSurvey:
			c.JSON(http.StatusBadRequest, gin.H{"error": "SurveyRequirements is nil"})
		case survey.ErrSurveyName:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Name field is empty"})
		case survey.ErrSurveyQuestion:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Questions field is empty"})
		default:
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"created_survey": resp})

}


func (h surveyHandler) GetSurveis(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	respnose, err := h.surveyService.GetSurviesByUserID(userID, c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, respnose)
}

type CloseSurveyRequestBody struct {
	SurveyId int `json:"survey_id"`
}

func (h surveyHandler) CloseSurvey(c *gin.Context) {
	var request *CloseSurveyRequestBody
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if request.SurveyId <= 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	surveyId := request.SurveyId

	err := h.surveyService.CloseSurvey(c, surveyId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h surveyHandler) GetSurvey(c *gin.Context) {
	surveyID, err := strconv.Atoi(c.Param("surveyId"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	if surveyID <= 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	res, err := h.surveyService.GetSurveyById(c, surveyID)
	if err != nil {
		if err == survey.ErrNotFound {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h surveyHandler) SaveSurvey(c *gin.Context) {
	var req models.SaveSurveyRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.surveyService.SaveSurvey(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	user, err := h.userService.GetUser(c, req.UserId)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	currentDir, _ := os.Getwd()

	outFileName := uuid.New().String() + ".pdf"
	htmlPath := currentDir + "/files/input_html/index.html"
	outPath := currentDir + "/files/output_pdf/" + outFileName
	temp, err := template.ParseFiles(htmlPath)
	if err != nil {
		fmt.Println(3)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pdfData := petition.PetitionData{}
	pdfData.CreationDate = helpers.CurrentDateModel()
	pdfData.Location = mockLocation
	pdfData.ResponsiblePerson = ""
	pdfData.OwnerName = *user.Username
	pdfData.OwnerAddress = "mock owner adress"
	pdfData.Questions = []petition.Question{}
	for i, question := range req.Questions {
		q := petition.Question{}
		q.Number = i
		q.Text = question.Text
		if question.AnswerId == 1 {
			q.Decision = "За"
		} else if question.AnswerId == 2 {
			q.Decision = "Против"
		} else if question.AnswerId == 3 {
			q.Decision = "Воздержусь"
		}

		pdfData.Questions = append(pdfData.Questions, q)
	}

	err = h.petitionService.GeneratePDF(c, temp, pdfData, outPath, htmlPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	finalFilePath, err := h.petitionService.GenerateFinalPdf(outPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Content-Type", "application/pdf")
	fmt.Println("file path: ", finalFilePath)
	
	c.File(finalFilePath)

}
