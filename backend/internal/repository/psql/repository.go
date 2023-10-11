package psql

import (
	"go.uber.org/zap"

	"mado/pkg/database/postgres"
)

// Repositories is a collection of all repositories in the system.
type Repositories struct {
	User   UserRepository
	Survey SurveyrRepository
	Petition PetitionRepository
}

// NewRepositories returns a new instance of Repositories.
func NewRepositories(db *postgres.Postgres, logger *zap.Logger) Repositories {
	return Repositories{
		User:   NewUserRepository(db, logger),
		Survey: NewSurveyrRepository(db, logger),
		Petition: NewPetitionRepository(db, logger),
	}
}
