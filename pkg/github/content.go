package github

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/tidwall/gjson"
)

// Content represents the content (description) of a github issue, pull request.
type Content struct {
	Title string
	Body  string
	URL   string
}

const FetchContentEndpoint = "https://api.github.com/repos/%s/%s/issues/%d"

type FetchContentRequest struct {
	Owner  string
	Repo   string
	Number int
}

// FetchContent fetches the content of a Github Issue or Pull Request
func FetchContent(ctx context.Context, req FetchContentRequest) (*Content, error) {
	u := fmt.Sprintf(FetchContentEndpoint, req.Owner, req.Repo, req.Number)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	data, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	title := gjson.GetBytes(data, "title").String()
	body := gjson.GetBytes(data, "body").String()
	url := gjson.GetBytes(data, "html_url").String()

	return &Content{
		Title: title,
		Body:  body,
		URL:   url,
	}, nil
}
