package util

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

func Fetcher(ctx context.Context, fullURL string, out any) error {
	rq := &http.Client{
		Timeout: 5 * time.Millisecond,
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return err
	}

	resp, err := rq.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(out)
}
