package psql

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"mado/internal/core/user"
	"mado/pkg/database/postgres"
	"mado/pkg/errs"
	"mado/pkg/logger"
)

// UserRepository is a user repository.
type UserRepository struct {
	db     *postgres.Postgres
	logger *zap.Logger
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *postgres.Postgres, logger *zap.Logger) UserRepository {
	return UserRepository{
		db:     db,
		logger: logger,
	}
}

func (ur UserRepository) CheckIfUserExistsByIIN(ctx context.Context, iin string) (bool, error) {
	// Ensure you have a valid database connection
	if ur.db == nil {
		ur.logger.Error(errs.ErrDBConnectionIsNill.Error())
		return false, errs.ErrDBConnectionIsNill
	}

	// Prepare the SQL statement to count the number of users with the given iin
	sqlStatement := `
        SELECT COUNT(*) FROM users WHERE iin = $1;
    `

	logger.FromContext(ctx).Debug("check user existence by iin query", zap.String("sql", sqlStatement), zap.String("iin", iin))

	var count int
	err := ur.db.Pool.QueryRow(ctx, sqlStatement, iin).Scan(&count)
	if err != nil {
		ur.logger.Error(errs.ErrGettingUsersCount.Error(), zap.Error(err))
		return false, fmt.Errorf("%w", err)
	}

	// If the count is greater than 0, it means a user with the given iin already exists
	return count > 0, nil
}

// TODO do it properly
func (ur UserRepository) Create(ctx context.Context, dto *user.User) (*user.User, error) {

	// Ensure you have a valid database connection
	if ur.db == nil {
		ur.logger.Error(errs.ErrDBConnectionIsNill.Error())
		return nil, errs.ErrDBConnectionIsNill
	}

	// Prepare the SQL statement
	sqlStatement := `
		INSERT INTO users (iin, email, bin, name, is_manager) 
		VALUES ($1, $2, $3, $4, $5);
		`
	logger.FromContext(ctx).Debug("create user query", zap.String("sql", sqlStatement), zap.Any("args", dto))

	// Execute the SQL statement
	result, err := ur.db.Pool.Exec(ctx, sqlStatement, dto.IIN, dto.Email, dto.BIN, dto.Username, false)
	if err != nil {
		ur.logger.Error(errs.ErrInsertingUser.Error(), zap.Error(err))
		return nil, fmt.Errorf("%w%w", errs.ErrInsertUser, err)
	}

	// Check the number of rows affected (usually for error checking)
	rowsAffected := result.RowsAffected()
	if rowsAffected != 1 {
		ur.logger.Error(errs.ErrRowsAffected.Error())
		return nil, fmt.Errorf("expected 1 row to be affected, but %d rows were affected", rowsAffected)
	}

	// Optionally, you can retrieve the newly inserted user if your database supports returning the inserted row.
	// Otherwise, you may want to fetch the user by some unique identifier (e.g., ID) and return it here.

	return dto, nil
}

func (ur UserRepository) GetAllRows(ctx context.Context) ([]*user.User, error) {
	if ur.db == nil {
		ur.logger.Error(errs.ErrDBConnectionIsNill.Error())
		return nil, errs.ErrDBConnectionIsNill
	}

	// Prepare the SQL statement
	sqlStatement := `SELECT * FROM users`

	// Execute the SQL statement and retrieve the result set
	rows, err := ur.db.Pool.Query(ctx, sqlStatement)
	if err != nil {
		ur.logger.Error(errs.ErrGetAllRows.Error(), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var users []*user.User

	// Iterate through the result set and scan each row into a user struct
	for rows.Next() {
		var u user.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.BIN, &u.Email); err != nil {
			ur.logger.Error(errs.ErrGetAllRows.Error(), zap.Error(err))
			return nil, err
		}
		users = append(users, &u)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		ur.logger.Error(errs.ErrGetAllRows.Error(), zap.Error(err))
		return nil, err
	}

	return users, nil
}
