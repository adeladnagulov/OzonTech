package domain

type Status string

const (
	PLANNED     Status = "planned"
	IN_PROGRESS Status = "in_progress"
	CANCELED    Status = "canceled"
	DONE        Status = "done"
)

func ValidStatus(status string) (Status, error) {
	switch status {
	case string(PLANNED):
		return PLANNED, nil
	case string(IN_PROGRESS):
		return IN_PROGRESS, nil
	case string(CANCELED):
		return CANCELED, nil
	case string(DONE):
		return DONE, nil
	default:
		return "", ErrIncorrectStatus
	}
}
