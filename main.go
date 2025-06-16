package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/theokyle/blog_aggregator/internal/config"
	"github.com/theokyle/blog_aggregator/internal/database"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

func main() {
	gatorConfig, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	fileContents, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	db, err := sql.Open("postgres", fileContents.DbURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dbQueries := database.New(db)

	st := state{
		db:     dbQueries,
		config: &gatorConfig,
	}

	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAggregate)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerGetFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

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
}
