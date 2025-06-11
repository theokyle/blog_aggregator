package main

import (
	"fmt"

	"github.com/theokyle/blog_aggregator/internal/config"
)

func main() {
	gatorConfig, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	err = gatorConfig.SetUser("Kyle")
	if err != nil {
		fmt.Println(err)
	}

	fileContents, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("db_url: %s current_user_name: %s\n", fileContents.DbURL, fileContents.CurrentUserName)
}
