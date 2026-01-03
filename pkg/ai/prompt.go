package ai

// 2 modes:
// 1. Quick summary:
// 2. Deep analysis:

type PromptMode int

const (
	PromptModeQuickSummary PromptMode = iota
	PromptModeDeepAnalysis
)

func (pm PromptMode) String() string {
	switch pm {
	case PromptModeQuickSummary:
		return `
		You are an expert code reviewer analyzing a GitHub Pull Request.
		**Task**: Provide a quick, actionable summary of the changes.
		**PR Title**: %s
		**PR Description**: %s
		**Code Changes**: %s
		**Comments**: %s

		**Required Analysis**:
		1. **Overview** (2-3 sentences): What is this PR trying to achieve?
		2. **Key Changes** (list each important change):
			- File: [filename]
			- Type: [bugfix/enhancement/feature/refactor]
			- Importance: [critical/high/medium/low]
			- Description: What changed and why
			- Details: How it solves the problem
		3. **Change Relationships**: How do these changes connect? What's the flow?
		4. **Impact Assessment**: 
			- What could break?
			- Performance implications?
		**Output Format**: Return as JSON:
		{
		"overview": "...",
		"key_changes": [
			{
			"file": "...",
			"type": "...",
			"importance": "...",
			"description": "...",
			"details": "..."
			}
		],
		"change_relationships": "...",
		"impact_assessment": "..."
		}
		`
	case PromptModeDeepAnalysis:
		return `
		You are a senior engineer mentoring developers through code review.
		**Task**: Identify learning opportunities and build knowledge around this PR.
		**PR Title**: %s
		**PR Description**: %s
		**Code Changes**: %s
		**Comments**: %s

		**Required Analysis**:
		1. **Highlights & Best Practices** (find 2-4 excellent techniques):
			- Topic: [Concept name]
			- Description: What's being done well
			- Code Example: Key snippet (if applicable)
			- Why it's good: Educational value
		2. **Knowledge Tree** (build around 1 core concept):
		- Core Concept: Main topic (e.g., "Lambda ImageConfig", "Error Handling Pattern")
		- Branches (2-3 sub-topics):
			* Topic: [sub-topic]
			* Current Usage: How it's used in this PR
			* Related Topics: [list of related concepts]
			* Advanced Concepts: [what to study next]
		3. **Learning Resources**:
			- Area: [Learning topic]
			- Suggestion: What to study and why
			- Resources: [relevant docs, articles, repos]
		
		**Output Format**: Return as JSON:
		{
			"highlights": [
				{
				"topic": "...",
				"description": "...",
				"code_example": "...",
				"why_good": "..."
				}
			],
			"knowledge_tree": {
				"core_concept": "...",
				"branches": [
				{
					"topic": "...",
					"current_usage": "...",
					"related_topics": ["...", "..."],
					"advanced_concepts": ["...", "..."]
				}
				]
			},
			"recommendations": [
				{
				"area": "...",
				"suggestion": "...",
				"resources": ["...", "..."]
				}
			]
		}
		`
	default:
		return "Unknown"
	}
}
