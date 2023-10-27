package core

import (
	"go.uber.org/zap"

	"mado/internal/core/petition"
	"mado/internal/core/survey"
	"mado/internal/core/user"
	"mado/internal/repository/psql"
)

// Services is a collection of all services in the system.
type Services struct {
	User     user.Service
	Survey   survey.Service
	Petition petition.Service
}

// NewServices returns a new instance of Services.
func NewServices(repositories psql.Repositories, userRepo psql.UserRepository, logger *zap.Logger) Services {
	return Services{
		User:     user.NewService(repositories.User, logger),
		Survey:   survey.NewService(repositories.Survey, logger),
		Petition: petition.NewService(repositories.Petition, userRepo, logger),
	}
}
