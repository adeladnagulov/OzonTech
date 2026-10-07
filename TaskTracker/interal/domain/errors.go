package domain

import "errors"

var (
	ErrIncurrectData     = errors.New("data obtained is incorrect")
	ErrDeadlineInThePast = errors.New("set deadline in the past")
)
