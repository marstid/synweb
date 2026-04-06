package search

import "github.com/marstid/synweb/pkg/errors"

type SearchRequest struct {
	Query string `json:"query"`
}

type SearchResponse struct {
	Results []SearchResult `json:"results"`
}

type SearchResult struct {
	URL       string  `json:"url"`
	Title     string  `json:"title"`
	Text      string  `json:"text"`
	Published *string `json:"published,omitempty"`
}

type SearchParams struct {
	Query         string
	MaxTextLength int
}

func DefaultSearchParams() *SearchParams {
	return &SearchParams{
		Query:         "",
		MaxTextLength: 1000,
	}
}

func (p *SearchParams) Validate() error {
	if p.Query == "" {
		return errors.ErrMissingQuery
	}
	return nil
}

func (p *SearchParams) WithQuery(query string) *SearchParams {
	p.Query = query
	return p
}

func (p *SearchParams) WithMaxTextLength(max int) *SearchParams {
	p.MaxTextLength = max
	return p
}
