package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/theokyle/blog_aggregator/internal/database"
)

type command struct {
	name string
	args []string
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

	username := cmd.args[0]

	user, err := s.db.GetUser(context.Background(), username)
	if err == sql.ErrNoRows {
		fmt.Println("Error: User not found")
		os.Exit(1)
	} else if err != nil {
		return err
	}

	err = s.config.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("User %s has been set\n", user.Name)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no username provided")
	}

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}

	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			fmt.Println("Error: Duplicate User")
			os.Exit(1)
		}
		return err
	}

	err = s.config.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("User: %v was created.", user.Name)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("Users successfully reset")
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name == s.config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}
