package main

import (
	"fmt"

	"github.com/theokyle/blog_aggregator/internal/config"
)

type command struct {
	name string
	args []string
}

type state struct {
	config *config.Config
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.handlers[cmd.name]
	if ok {
		err := f(s, cmd)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("command not found")
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no username provided")
	}

	user := cmd.args[0]

	err := s.config.SetUser(user)
	if err != nil {
		return err
	}

	fmt.Printf("User %s has been set\n", user)

	return nil
}
