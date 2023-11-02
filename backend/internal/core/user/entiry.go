package user

type User struct {
	ID         int     `json:"id"`
	Username   *string `json:"username"`
	Email      *string `json:"email"`
	IIN        *string `json:"iin"`
	BIN        *string `json:"bin"`
	Is_manager *bool   `json:"is_manager"`
}
