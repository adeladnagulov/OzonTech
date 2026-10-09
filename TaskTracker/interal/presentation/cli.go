package presentation

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"task-tracker/interal/domain"
	"task-tracker/interal/usecase"
	"time"
)

type Cli struct {
	scan    *bufio.Scanner
	service *usecase.TaskService
}

func NewCli(service *usecase.TaskService) *Cli {
	return &Cli{
		scan:    bufio.NewScanner(os.Stdin),
		service: service,
	}
}

// help, remove, update и упростить
func (c *Cli) Run() error {
	for {
		command, args, err := c.parseInput()
		if err != nil {
			return err
		}

		switch command {
		case string(domain.EXIT):
			return nil
		case string(domain.ADD):
			if len(args) != 2 {
				fmt.Println("Формат ввода должен быть следующим: add задача_1 2006-01-02")
				continue
			}
			title := args[0]
			deadline, err := time.Parse("2006-01-02", args[1])
			if err != nil {
				fmt.Println("Некорректный формат дедлайна")
				continue
			}
			err = c.service.AddNewTask(title, deadline)
			if err != nil {
				if errors.Is(err, domain.ErrDeadlineInThePast) {
					fmt.Println("Указанная дата дедлайна уже прошла")
					continue
				}
				return err
			}
		case string(domain.GET_ALL):
			tasks, err := c.service.GetAllTasks()
			if err != nil {
				return err
			}
			for _, v := range tasks {
				fmt.Println(v)
				fmt.Println()
			}
		case string(domain.GET):
			if len(args) != 1 {
				fmt.Println("Формат ввода должен быть следующим: get 1")
				continue
			}
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("Переданный id некорретный")
				continue
			}
			task, err := c.service.GetTask(int64(id))
			if err != nil {
				if errors.Is(err, domain.ErrTaskNotFound) {
					fmt.Println("Задачи не существует")
					continue
				}
				return err
			}
			fmt.Println(task)
		case string(domain.CHANGE_STAUS):
			if len(args) != 2 {
				fmt.Println("Формат ввода должен быть следующим: change_status 1 done")
				continue
			}
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("Переданный id некорретный")
				continue
			}
			status := args[1]
			err = c.service.ChangeStatus(int64(id), status)
			if err != nil {
				if errors.Is(err, domain.ErrStatusAlreadySet) {
					fmt.Println("Задача уже имеет данный статус")
					continue
				}
				if errors.Is(err, domain.ErrConflictStatus) {
					fmt.Println("Задачу в статусе `done` или `canceled` менять нельзя")
					continue
				}
				if errors.Is(err, domain.ErrIncorrectStatus) {
					fmt.Println("Не существующий статус")
					continue
				}
				return err
			}
		case string(domain.DELETE):
			if len(args) != 1 {
				fmt.Println("Формат ввода должен быть следующим: delete 1")
				continue
			}
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("Переданный id некорретный")
				continue
			}
			err = c.service.DeleteTask(int64(id))
			if err != nil {
				if errors.Is(err, domain.ErrTaskNotFound) {
					fmt.Println("Задачи не существует")
					continue
				}
				return err
			}
		default:
			fmt.Println("Неизвестная команда")
		}
	}
}

func (c *Cli) parseInput() (string, []string, error) {
	if !c.scan.Scan() {
		return "", nil, c.scan.Err()
	}
	line := c.scan.Text()
	parts := strings.Fields(line)
	command, args := parts[0], parts[1:]
	return command, args, nil
}
