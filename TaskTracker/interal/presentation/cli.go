package presentation

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"task-tracker/interal/domain"
)

type Cli struct {
	scan *bufio.Scanner
}

func NewCli() *Cli {
	return &Cli{
		scan: bufio.NewScanner(os.Stdin),
	}
}

func (c *Cli) Run() error {
	for {
		command, _, err := c.parseInput()
		if err != nil {
			return err
		}

		switch command {
		case string(domain.EXIT):
			return nil
		default:
			fmt.Println("неизвестная команда")
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
