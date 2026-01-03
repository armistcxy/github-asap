package github

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/tidwall/gjson"
)

// Comment represents a GitHub issue or pull request comment
type Comment struct {
	Body          string
	URL           string
	ReactionCount int
}

const FetchCommentsEndpoint = "https://api.github.com/repos/%s/%s/issues/%d/comments"

type FetchCommentsRequest struct {
	Owner    string
	Repo     string
	Number   int
	MaxFetch int
}

type FetchCommentsResponse struct {
	Comments []Comment
}

// FetchComments fetches comments from Github Issue or Pull Request
func FetchComments(ctx context.Context, req FetchCommentsRequest) (*FetchCommentsResponse, error) {
	u := fmt.Sprintf(FetchCommentsEndpoint, req.Owner, req.Repo, req.Number)

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

	comments := make([]Comment, 0)

	gjson.Parse(string(data)).ForEach(func(key, value gjson.Result) bool {
		body := value.Get("body").String()
		url := value.Get("html_url").String()
		reactionCount := value.Get("reactions.total_count").Int()

		comment := Comment{
			Body:          body,
			URL:           url,
			ReactionCount: int(reactionCount),
		}

		comments = append(comments, comment)

		if req.MaxFetch > 0 && len(comments) >= req.MaxFetch {
			return false
		}

		return true
	})

	return &FetchCommentsResponse{Comments: comments}, nil
}
