package ai

import (
	"context"
	"os"
	"testing"

	"github.com/tmc/langchaingo/llms/openai"
)

func setupAIAnalyzer(t *testing.T) AIAnalyzer {
	llm, err := openai.New(
		openai.WithModel("gpt-5"),
		openai.WithToken(os.Getenv("OPENAI_API_KEY")),
	)
	if err != nil {
		t.Fatalf("failed to create OpenAI LLM: %v", err)
	}

	return NewAIAnalyzer(llm)
}

func TestAnalyzePRQuickSummary(t *testing.T) {
	analyzer := setupAIAnalyzer(t)

	req := &AnalyzePRRequest{
		Owner:  "argoproj",
		Repo:   "argo-cd",
		Number: 7539,
		Mode:   PromptModeQuickSummary,
	}

	result, err := analyzer.AnalyzePR(context.Background(), req)
	if err != nil {
		t.Fatalf("AnalyzePR failed: %v", err)
	}

	t.Logf("Analyze Mode: %s", req.Mode)
	t.Logf("AnalyzePR result: %s", result)
}

func TestAnalyzePRDeepAnalysis(t *testing.T) {
	analyzer := setupAIAnalyzer(t)

	req := &AnalyzePRRequest{
		Owner:  "argoproj",
		Repo:   "argo-cd",
		Number: 7539,
		Mode:   PromptModeDeepAnalysis,
	}

	result, err := analyzer.AnalyzePR(context.Background(), req)
	if err != nil {
		t.Fatalf("AnalyzePR failed: %v", err)
	}

	t.Logf("Analyze Mode: %s", req.Mode)
	t.Logf("AnalyzePR result: %s", result)
}
