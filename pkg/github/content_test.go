package github

import (
	"context"
	"testing"
)

func TestFetchContent(t *testing.T) {
	exampleReq := FetchContentRequest{
		Owner:  "argoproj",
		Repo:   "argo-cd",
		Number: 4972,
	}

	content, err := FetchContent(context.TODO(), exampleReq)
	if err != nil {
		t.Fatalf("FetchContent failed: %v", err)
	}
	if content.URL == "" {
		t.Errorf("Expected non-empty content URL")
	}
	if content.Body == "" {
		t.Errorf("Expected non-empty content body")
	}

	t.Logf("Content Title: %s", content.Title)
	t.Logf("Content URL: %s", content.URL)
	t.Logf("Content Body: %s", content.Body)
}
