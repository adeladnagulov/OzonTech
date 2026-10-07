package main

import "task-tracker/interal/presentation"

func main() {
	cli := presentation.NewCli()
	cli.Run()
}
