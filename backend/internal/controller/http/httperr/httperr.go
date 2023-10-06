package httperr

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"mado/pkg/logger"
)

// InternalError is a helper function to respond with internal server error.
func InternalError(c *gin.Context, slug string, err error) {
	httpRespondWithError(c, err, slug, "Internal server error", http.StatusInternalServerError)
}

// UnprocessableEntity is a helper function to respond with unprocessable entity error.
func UnprocessableEntity(c *gin.Context, slug string, err error) {
	httpRespondWithError(c, err, slug, "Unprocessable entity", http.StatusUnprocessableEntity)
}

// Unauthorised is a helper function to respond with unauthorized error.
func Unauthorised(c *gin.Context, slug string, err error) {
	httpRespondWithError(c, err, slug, "Unauthorised", http.StatusUnauthorized)
}

// BadRequest is a helper function to respond with bad request error.
func BadRequest(c *gin.Context, slug string, err error) {
	httpRespondWithError(c, err, slug, "Bad request", http.StatusBadRequest)
}

// Errors is a list of http errors.
type Errors struct {
	Body []string `json:"body"`
}

// ErrorResponse is a error response.
type ErrorResponse struct {
	Errors Errors `json:"errors"`
}

func httpRespondWithError(c *gin.Context, err error, slug string, logMessage string, status int) {
	logger.FromRequest(c).Warn(logMessage, zap.Error(err), zap.String("error-slug", slug))
	c.AbortWithStatusJSON(status, ErrorResponse{
		Errors: Errors{
			Body: []string{slug},
		},
	})
}
