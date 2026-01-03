package ai

import (
	"context"
	"os"
	"testing"

	"github.com/tmc/langchaingo/llms/openai"
)

func TestAnalyzePR(t *testing.T) {
	llm, err := openai.New(
		openai.WithModel("gpt-5"),
		openai.WithToken(os.Getenv("OPENAI_API_KEY")),
	)
	if err != nil {
		t.Fatalf("failed to create OpenAI LLM: %v", err)
	}

	analyzer := NewAIAnalyzer(llm)

	testcases := []struct {
		name string
		req  *AnalyzePRRequest
	}{
		{
			name: "Quick Summary",
			req: &AnalyzePRRequest{
				Owner:  "argoproj",
				Repo:   "argo-cd",
				Number: 7539,
				Mode:   PromptModeQuickSummary,
			},
		},
		{
			name: "Deep Analysis",
			req: &AnalyzePRRequest{
				Owner:  "argoproj",
				Repo:   "argo-cd",
				Number: 7539,
				Mode:   PromptModeDeepAnalysis,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := analyzer.AnalyzePR(context.Background(), tc.req)
			if err != nil {
				t.Fatalf("AnalyzePR failed: %v", err)
			}

			t.Logf("Analyze Mode: %s", tc.req.Mode)
			t.Logf("AnalyzePR result: %s", result)
		})
	}
}
