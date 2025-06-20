package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
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

	fmt.Printf("User: %v was created.\n", user.Name)
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

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("error: missing name/url")
	}

	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			fmt.Println("Error: Duplicate Feed")
			os.Exit(1)
		}
		return err
	}

	follow_params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), follow_params)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			fmt.Println("Error: Duplicate Follow")
			os.Exit(1)
		}
		return err
	}

	fmt.Printf("Feed %s was created by User %s.\n", feed.Name, user.Name)
	return nil
}

func handlerGetFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("- Feed: %s, URL: %s, User: %s\n", feed.Name, feed.Url, feed.Username)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("error: missing url")
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			fmt.Println("Error: Duplicate Follow")
			os.Exit(1)
		}
		return err
	}

	fmt.Printf("Feed %s now followed by User %s.", feed.Name, user.Name)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	fmt.Println("Current user is following:")

	for _, follow := range follows {
		fmt.Printf("- Feed: %s\n", follow.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("error: missing url")
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	params := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	err = s.db.DeleteFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("Feed at %s successfully deleted.\n", feed.Url)
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.args) == 1 {
		if specifiedLimit, err := strconv.Atoi(cmd.args[0]); err == nil {
			limit = specifiedLimit
		} else {
			return fmt.Errorf("invalid limit: %w", err)
		}
	}

	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}

	posts, err := s.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("couldn't get posts for user: %w", err)
	}

	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.Name)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}

	return nil
}
