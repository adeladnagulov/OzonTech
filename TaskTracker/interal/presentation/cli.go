package presentation

import (
	"bufio"
	"errors"
	"fmt"
	"os"
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
