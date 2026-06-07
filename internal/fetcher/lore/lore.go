package lore

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/html/charset"
	"github.com/emmanuel326/gopher_intel/internal/fetcher"
)

type feed struct {
	Entries []entry `xml:"entry"`
}

type entry struct {
	ID      string `xml:"id"`
	Title   string `xml:"title"`
	Author  author `xml:"author"`
	Updated string `xml:"updated"`
	Link    link   `xml:"link"`
	Content string `xml:"content"`
}

type author struct {
	Name string `xml:"name"`
}

type link struct {
	Href string `xml:"href,attr"`
}

type LoreSource struct {
	list string
}

func New(list string) *LoreSource {
	return &LoreSource{list: list}
}

func (l *LoreSource) Name() string {
	return l.list
}

func (l *LoreSource) Fetch() ([]fetcher.Message, error) {
	url := fmt.Sprintf("https://lore.kernel.org/%s/new.atom", l.list)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("request %s: %w", l.list, err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", l.list, err)
	}
	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)
	decoder.CharsetReader = charset.NewReaderLabel

	var f feed
	if err := decoder.Decode(&f); err != nil {
		return nil, fmt.Errorf("decode %s: %w", l.list, err)
	}

	var messages []fetcher.Message
	for _, e := range f.Entries {
		date, _ := time.Parse(time.RFC3339, e.Updated)
		messages = append(messages, fetcher.Message{
			ID:      e.ID,
			Source:  l.list,
			Subject: e.Title,
			Author:  e.Author.Name,
			Date:    date,
			URL:     e.Link.Href,
			Body:    e.Content,
		})
	}

	return messages, nil
}
