package github

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/tidwall/gjson"
)

// CodeChange represents a code change (a file patch) in a github pull request.
type CodeChange struct {
	FileName string

	// Patch is the diff/patch of the file change.
	Patch string

	// RawURL is the URL to fetch the raw content of the changed file.
	RawURL string

	BlobURL string
}

const FetchPRCodeChangesEndpoint = "https://api.github.com/repos/%s/%s/pulls/%d/files"

type FetchPRCodeChangesRequest struct {
	Owner  string
	Repo   string
	Number int
}

type FetchPRCodeChangesResponse struct {
	CodeChanges []CodeChange
}

// FetchPRCodeChanges fetches the code changes (file patches) from a Github Pull Request
func FetchPRCodeChanges(ctx context.Context, req FetchPRCodeChangesRequest) (*FetchPRCodeChangesResponse, error) {
	u := fmt.Sprintf(FetchPRCodeChangesEndpoint, req.Owner, req.Repo, req.Number)

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

	codeChanges := make([]CodeChange, 0)

	gjson.ParseBytes(data).ForEach(func(key, value gjson.Result) bool {
		fileName := value.Get("filename").String()
		patch := value.Get("patch").String()
		rawURL := value.Get("raw_url").String()
		blobURL := value.Get("blob_url").String()

		codeChange := CodeChange{
			FileName: fileName,
			Patch:    patch,
			RawURL:   rawURL,
			BlobURL:  blobURL,
		}
		codeChanges = append(codeChanges, codeChange)
		return true
	})

	resp := &FetchPRCodeChangesResponse{
		CodeChanges: codeChanges,
	}
	return resp, nil
}
