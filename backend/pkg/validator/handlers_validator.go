package validator

import (
	"strconv"

	"mado/pkg/errs"
)

func IdValidator(id string) (*int, error) {
	idv, err := strconv.Atoi(id)
	if err != nil || idv <= 0 {
		return nil, errs.ErrInvalidID
	}
	return &idv, nil
}
