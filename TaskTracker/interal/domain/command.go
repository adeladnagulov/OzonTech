package domain

type Command string

const (
	HELP Command = "help"
	EXIT Command = "exit"
	ADD  Command = "add"
)
