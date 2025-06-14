package rss

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}

	req.Header.Set("User-Agent", "gator")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	rssResponse := RSSFeed{}
	if err := xml.Unmarshal(data, &rssResponse); err != nil {
		return &RSSFeed{}, err
	}

	rssResponse.Channel.Title = html.UnescapeString(rssResponse.Channel.Title)
	rssResponse.Channel.Description = html.UnescapeString(rssResponse.Channel.Description)
	for i, item := range rssResponse.Channel.Item {
		rssResponse.Channel.Item[i].Title = html.UnescapeString(item.Title)
		rssResponse.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}

	return &rssResponse, nil
}
