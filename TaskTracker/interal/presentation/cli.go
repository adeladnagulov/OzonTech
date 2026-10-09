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
			err := c.add(args)
			if err != nil {
				return err
			}
		case string(domain.GET_ALL):
			err := c.getAll()
			if err != nil {
				return err
			}
		case string(domain.GET):
			err := c.get(args)
			if err != nil {
				return err
			}
		case string(domain.CHANGE_STAUS):
			err := c.changeStatus(args)
			if err != nil {
				return err
			}
		case string(domain.DELETE):
			err := c.delete(args)
			if err != nil {
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

func (c *Cli) add(args []string) error {
	if len(args) != 2 {
		fmt.Println("Формат ввода должен быть следующим: add задача_1 2006-01-02")
		return nil
	}
	title := args[0]
	deadline, err := time.Parse("2006-01-02", args[1])
	if err != nil {
		fmt.Println("Некорректный формат дедлайна")
		return nil
	}
	err = c.service.AddNewTask(title, deadline)
	if err != nil {
		if errors.Is(err, domain.ErrDeadlineInThePast) {
			fmt.Println("Указанная дата дедлайна уже прошла")
			return nil
		}
		return err
	}
	return nil
}

func (c *Cli) getAll() error {
	tasks, err := c.service.GetAllTasks()
	if err != nil {
		return err
	}
	for _, v := range tasks {
		fmt.Println(v)
		fmt.Println()
	}
	return nil
}

func (c *Cli) get(args []string) error {
	if len(args) != 1 {
		fmt.Println("Формат ввода должен быть следующим: get 1")
		return nil
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Переданный id некорретный")
		return nil
	}
	task, err := c.service.GetTask(int64(id))
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			fmt.Println("Задачи не существует")
			return nil
		}
		return err
	}
	fmt.Println(task)
	return nil
}

func (c *Cli) changeStatus(args []string) error {
	if len(args) != 2 {
		fmt.Println("Формат ввода должен быть следующим: change_status 1 done")
		return nil
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Переданный id некорретный")
		return nil
	}
	status := args[1]
	err = c.service.ChangeStatus(int64(id), status)
	if err != nil {
		if errors.Is(err, domain.ErrStatusAlreadySet) {
			fmt.Println("Задача уже имеет данный статус")
			return nil
		}
		if errors.Is(err, domain.ErrConflictStatus) {
			fmt.Println("Задачу в статусе `done` или `canceled` менять нельзя")
			return nil
		}
		if errors.Is(err, domain.ErrIncorrectStatus) {
			fmt.Println("Не существующий статус")
			return nil
		}
		return err
	}
	return nil
}

func (c *Cli) delete(args []string) error {
	if len(args) != 1 {
		fmt.Println("Формат ввода должен быть следующим: delete 1")
		return nil
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Переданный id некорретный")
		return nil
	}
	err = c.service.DeleteTask(int64(id))
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			fmt.Println("Задачи не существует")
			return nil
		}
		return err
	}
	return nil
}
