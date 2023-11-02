package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"mado/internal/config"
	"mado/internal/controller/http/middleware"
	"mado/internal/core"
)

// Deps is a http handler dependencies.
type Deps struct {
	Logger   *zap.Logger
	Services core.Services
}

// NewRouter returns a new http router.
func NewRouter(deps Deps) *gin.Engine {
	router := gin.New()

	if config.Get().IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}

	middleware.ApplyMiddlewares(router, deps.Logger)
	// Add the Gin logger middleware to log request information
	router.Use(gin.Logger())
	api := router.Group("/api")
	{

		newPetitionHandler(petitionDeps{
			router:          api,
			petitionService: deps.Services.Petition,
			userService:     deps.Services.User,
			surveyService:   deps.Services.Survey,
		})

		newUserHandler(userDeps{
			router: api,

			userService: deps.Services.User,
		})

		newSurveyHandler(
			surveyDeps{
				router:        api,
				surveyService: deps.Services.Survey,
			})
		newTestURLHandler(testURLDeps{
			router: router,
		})
	}

	return router
}
