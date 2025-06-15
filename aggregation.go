package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

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

	return nil
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
	}

	return nil
}
