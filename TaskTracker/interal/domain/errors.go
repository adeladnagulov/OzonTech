package domain

import "errors"

var (
	ErrIncorrectData     = errors.New("data obtained is incorrect")
	ErrDeadlineInThePast = errors.New("set deadline in the past")
	ErrTaskNotFound      = errors.New("task could not be found")
	ErrStatusAlreadySet  = errors.New("status has already been set")
	ErrConflictStatus    = errors.New("task with the status `done` or `canceled` cannot be changed")
	ErrIncorrectStatus   = errors.New("status obtained is incorrect")
)
