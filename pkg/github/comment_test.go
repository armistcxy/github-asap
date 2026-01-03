package github

import (
	"context"
	"testing"
)

func TestFetchComments(t *testing.T) {
	exampleReq := FetchCommentsRequest{
		Owner:    "argoproj",
		Repo:     "argo-cd",
		Number:   4972,
		MaxFetch: 3,
	}

	resp, err := FetchComments(context.TODO(), exampleReq)
	if err != nil {
		t.Fatalf("FetchComments returned error: %v", err)
	}

	if len(resp.Comments) == 0 {
		t.Fatalf("FetchComments returned zero comments")
	}

	for _, comment := range resp.Comments {
		t.Logf("Comment URL: %s, ReactionCount: %d", comment.URL, comment.ReactionCount)
		t.Logf("Comment Body: %s", comment.Body)
	}
}
