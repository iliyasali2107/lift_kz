package user

import (
	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	Username   *string   `json:"username"`
	Email      *string   `json:"email"`
	IIN        *string   `json:"iin"`
	BIN        *string   `json:"bin"`
	Is_manager *bool     `json:"is_manager"`
}
