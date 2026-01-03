package github

import (
	"context"
	"testing"
)

func TestFetchCodeChanges(t *testing.T) {
	exampleReq := FetchPRCodeChangesRequest{
		Owner:  "argoproj",
		Repo:   "argo-cd",
		Number: 7539,
	}

	resp, err := FetchPRCodeChanges(context.TODO(), exampleReq)
	if err != nil {
		t.Fatalf("FetchPRCodeChanges returned error: %v", err)
	}

	if len(resp.CodeChanges) == 0 {
		t.Fatalf("FetchPRCodeChanges returned zero code changes")
	}

	for _, change := range resp.CodeChanges {
		t.Logf("File: %s, RawURL: %s, BlobURL: %s", change.FileName, change.RawURL, change.BlobURL)
		t.Logf("Patch:\n%s", change.Patch)
	}
}
