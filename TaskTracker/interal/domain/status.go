package domain

type Status string

const (
	planned     Status = "planned"
	in_progress Status = "in_progress"
	canceled    Status = "canceled"
	done        Status = "done"
)
