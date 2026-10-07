package main

import (
	"log"
	"task-tracker/interal/presentation"
	"task-tracker/interal/storage"
	"task-tracker/interal/usecase"
)

func main() {
	repo := storage.NewTaskRepo("tasks.json")
	service := usecase.NewTaskService(repo)
	cli := presentation.NewCli(service)
	err := cli.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
