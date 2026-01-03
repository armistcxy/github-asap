package ai

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/tmc/langchaingo/llms"

	"github.com/armistcxy/github-asap/pkg/github"
)

// AIAnalyzer defines the interface for analyzing GitHub Pull Requests, Issues using AI.
type AIAnalyzer interface {
	// AnalyzePR analyzes a GitHub Pull Request and returns insights based on the specified mode.
	AnalyzePR(ctx context.Context, req *AnalyzePRRequest) (string, error)
}

type AnalyzePRRequest struct {
	URL    string
	Owner  string
	Repo   string
	Number int

	Model string
	Mode  PromptMode
}

type AIAnalyzerOption func(*AnalyzePRRequest)

func WithModel(model string) AIAnalyzerOption {
	return func(req *AnalyzePRRequest) {
		req.Model = model
	}
}

func WithMode(mode PromptMode) AIAnalyzerOption {
	return func(req *AnalyzePRRequest) {
		req.Mode = mode
	}
}

type implAIAnalyzer struct {
	llm llms.Model
}

func NewAIAnalyzer(llm llms.Model) AIAnalyzer {
	return &implAIAnalyzer{
		llm: llm,
	}
}

func (a *implAIAnalyzer) AnalyzePR(ctx context.Context, req *AnalyzePRRequest) (string, error) {
	if req.Mode != PromptModeQuickSummary && req.Mode != PromptModeDeepAnalysis {
		return "", fmt.Errorf("unsupported prompt mode: %v", req.Mode)
	}

	if req.URL != "" {
		owner, repo, number, err := parseGitHubPRURL(req.URL)
		if err != nil {
			return "", fmt.Errorf("failed to parse GitHub PR URL: %w", err)
		}
		req.Owner = owner
		req.Repo = repo
		req.Number = number
	}

	if req.Owner == "" || req.Repo == "" || req.Number == 0 {
		return "", fmt.Errorf("invalid PR identifier: owner=%s, repo=%s, number=%d", req.Owner, req.Repo, req.Number)
	}

	// Fetch PR content (title and description)
	contentResp, err := github.FetchContent(ctx, github.FetchContentRequest{
		Owner:  req.Owner,
		Repo:   req.Repo,
		Number: req.Number,
	})
	if err != nil {
		return "", fmt.Errorf("failed to fetch PR content: %w", err)
	}

	// Fetch code changes
	codeChangesResp, err := github.FetchPRCodeChanges(ctx, github.FetchPRCodeChangesRequest{
		Owner:  req.Owner,
		Repo:   req.Repo,
		Number: req.Number,
	})
	if err != nil {
		return "", fmt.Errorf("failed to fetch code changes: %w", err)
	}

	// Format code changes
	var codeChangesStr strings.Builder
	if codeChangesResp != nil {
		for _, change := range codeChangesResp.CodeChanges {
			codeChangesStr.WriteString(fmt.Sprintf("File: %s\n%s\n\n", change.FileName, change.Patch))
		}
	}

	// Fetch comments
	commentsResp, err := github.FetchComments(ctx, github.FetchCommentsRequest{
		Owner:    req.Owner,
		Repo:     req.Repo,
		Number:   req.Number,
		MaxFetch: 25, // Limit to 25 comments
	})
	if err != nil {
		return "", fmt.Errorf("failed to fetch comments: %w", err)
	}

	// Format comments
	var commentsStr strings.Builder
	if commentsResp != nil {
		for _, comment := range commentsResp.Comments {
			commentsStr.WriteString(fmt.Sprintf("Comment: %s\n\n", comment.Body))
		}
	}

	// Get the prompt template for the specified mode
	promptTemplate := req.Mode.String()

	// Format the prompt with actual data
	prompt := fmt.Sprintf(promptTemplate,
		contentResp.Title,
		contentResp.Body,
		codeChangesStr.String(),
		commentsStr.String(),
	)

	// Call the LLM
	result, err := a.llm.GenerateContent(ctx, []llms.MessageContent{
		{
			Role: llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{
				llms.TextContent{Text: prompt},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate content from LLM: %w", err)
	}

	// Extract the response text
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from LLM")
	}

	return result.Choices[0].Content, nil
}

// parseGitHubPRURL parses a GitHub PR URL.
// Example: https://github.com/argoproj/argo-cd/pull/7539
// Returns: owner, repo, prNumber
func parseGitHubPRURL(raw string) (string, string, int, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", "", 0, err
	}

	// Expect host = github.com
	if u.Host != "github.com" {
		return "", "", 0, fmt.Errorf("not a github.com URL")
	}

	// Path: /owner/repo/pull/number
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) != 4 || parts[2] != "pull" {
		return "", "", 0, fmt.Errorf("invalid GitHub PR URL format")
	}

	prNumber, err := strconv.Atoi(parts[3])
	if err != nil {
		return "", "", 0, fmt.Errorf("invalid PR number")
	}

	return parts[0], parts[1], prNumber, nil
}
