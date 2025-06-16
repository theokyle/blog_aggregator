package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/theokyle/blog_aggregator/internal/database"
	"github.com/theokyle/blog_aggregator/internal/rss"
)

func handlerAggregate(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no duration provided")
	}

	time_between_reqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %v\n", time_between_reqs)

	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			fmt.Printf("Error scraping feeds: %v\n", err)
		}
	}
}

func scrapeFeeds(s *state) error {
	next_feed, err := s.db.GetNextFeedToFetch(context.Background())
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("No feeds to fetch, waiting for the next tick…")
		return nil
	} else if err != nil {
		return err
	}

	rssFeed, err := rss.FetchFeed(context.Background(), next_feed.Url)
	if err != nil {
		fmt.Printf("feed at %s encountered error: %v\n", next_feed.Url, err)
		return nil
	}

	err = s.db.MarkFeedFetched(context.Background(), next_feed.ID)
	if err != nil {
		return err
	}

	for _, item := range rssFeed.Channel.Item {
		fmt.Println(item.Title)

		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		// Add post to database, ignore duplicate error but log others

		params := database.CreatePostParams{
			ID:          uuid.New(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: publishedAt,
			FeedID:      next_feed.ID,
		}

		_, err = s.db.CreatePost(context.Background(), params)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}

	return nil
}
