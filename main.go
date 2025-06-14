package main

import (
	"fmt"
	"os"

	"github.com/theokyle/blog_aggregator/internal/config"
)

func main() {
	gatorConfig, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	st := state{
		config: &gatorConfig,
	}

	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("error: no command entered")
		os.Exit(1)
	}

	command := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	err = cmds.run(&st, command)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fileContents, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("db_url: %s current_user_name: %s\n", fileContents.DbURL, fileContents.CurrentUserName)
}
