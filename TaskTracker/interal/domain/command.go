package domain

type Command string

const (
	HELP         Command = "help"
	EXIT         Command = "exit"
	ADD          Command = "add"
	GET_ALL      Command = "get_all"
	GET          Command = "get"
	CHANGE_STAUS Command = "change_status"
	DELETE       Command = "delete"
	EDIT         Command = "edit"
)
