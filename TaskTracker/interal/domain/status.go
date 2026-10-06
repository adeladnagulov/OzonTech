package domain

type Status string

const (
	PLANNED     Status = "planned"
	IN_PROGRESS Status = "in_progress"
	CANCELED    Status = "canceled"
	DONE        Status = "done"
)
