package handler

import (
	"context"
	"net/http"
	"strconv"

	// "strconv"

	"github.com/gin-gonic/gin"

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
}

type UserReq struct {
	UserID int `json:"user_id"`
}

type surveyDeps struct {
	router *gin.RouterGroup

	surveyService SurveyService
}

type surveyHandler struct {
	surveyService SurveyService
}

func newSurveyHandler(deps surveyDeps) {
	handler := surveyHandler{
		surveyService: deps.surveyService,
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
}
