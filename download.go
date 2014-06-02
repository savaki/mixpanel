package mixpanel

import (
	"io"
	"net/url"
)

func (c *Client) Download(from, to, event string) (io.ReadCloser, error) {
	params := url.Values{}

	if from != "" {
		params.Add("from_date", from)
	}
	if to != "" {
		params.Add("to_date", to)
	}
	if event != "" {
		params.Add("event", event)
	}

	return c.get("/api/2.0/export/", params)
}
