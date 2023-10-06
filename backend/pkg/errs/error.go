package errs

import "errors"

var (
	ErrDBConnectionIsNill = errors.New("database connection is nil")
	ErrInsertUser         = errors.New("can not insert user")
	ErrGettingUsersCount  = errors.New("can not get count of users")
	ErrGetAllRows         = errors.New("can not get count of users")
	ErrInsertingUser      = errors.New("can not insert user")
	ErrRowsAffected       = errors.New("number of affected row have to be 1")
)
