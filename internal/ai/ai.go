package ai

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/genai"
	"github.com/emmanuel326/gopher_intel/internal/fetcher"
)

type Summarizer struct {
	client *genai.Client
	model  string
}

func New(apiKey string) (*Summarizer, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("create gemini client: %w", err)
	}

	return &Summarizer{
		client: client,
		model:  "gemini-2.0-flash-lite",
	}, nil
}

// SummarizeAll sends all sources in a single API call
func (s *Summarizer) SummarizeAll(data map[string][]fetcher.Message) (string, error) {
	var sb strings.Builder

	sb.WriteString(`You are a senior Linux kernel engineer and systems software analyst with 20+ years of experience.
You have deep expertise in kernel internals, storage subsystems, virtualization, eBPF, and low-level systems programming.

You are analyzing today's activity across multiple Linux kernel and systems mailing lists simultaneously.

Your job is to produce a unified INTELLIGENCE BRIEF across all lists.
Think like a staff engineer doing their morning review — what is happening, what matters,
who is driving it, and what signals matter for the broader ecosystem.

Be direct. Be technical. Be opinionated where the data supports it. No filler.

For each mailing list, produce the following structure:

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
 [LIST NAME] INTELLIGENCE BRIEF
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

SITUATION OVERVIEW
What is the overall pulse? Active development, stabilization, bug-fixing, or architecture debate?
What does the volume and nature of traffic signal?

KEY THEMES & TECHNICAL THREADS
3-5 dominant themes. For each: core problem, key engineers, current state, real-world impact.

NOTABLE PATCHES & RFCs
Most significant patches. Problem solved, fix type (correctness/perf/feature), any controversy?

CRITICAL BUGS & STABILITY SIGNALS
Bug reports, regressions, CVE-adjacent fixes. Severity and blast radius. Who is affected?

PEOPLE & DYNAMICS
Most active engineers. Maintainer reviews, Acked-by chains, disagreements, newcomers?

ANALYST TAKE
Opinionated read. What to watch closely? What trends matter in 3-6 months?

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
After all lists, add a final section:

 CROSS-LIST SIGNALS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
What themes or engineers appear across multiple lists?
Any systemic patterns — security hardening wave, performance push, API churn?
What does the aggregate activity tell us about where the kernel is heading?
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Here is today's data:

`)

	for source, messages := range data {
		sb.WriteString(fmt.Sprintf("\n== %s ==\n", source))
		for _, m := range messages {
			sb.WriteString(fmt.Sprintf("- [%s] %s: %s\n", m.Date.Format("2006-01-02"), m.Author, m.Subject))
		}
	}

	sb.WriteString("\nDeliver the full brief now. No preamble. Start with the first list directly.")

	resp, err := s.client.Models.GenerateContent(
		context.Background(),
		s.model,
		genai.Text(sb.String()),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("gemini summarize: %w", err)
	}

	return resp.Text(), nil
}

// SummarizeAllWithRetry retries on rate limit errors
func (s *Summarizer) SummarizeAllWithRetry(data map[string][]fetcher.Message) (string, error) {
	backoff := []time.Duration{20 * time.Second, 45 * time.Second, 90 * time.Second}

	for attempt, wait := range backoff {
		result, err := s.SummarizeAll(data)
		if err == nil {
			return result, nil
		}

		if attempt < len(backoff)-1 {
			fmt.Printf("rate limited, retrying in %s (attempt %d/3)...\n", wait, attempt+1)
			time.Sleep(wait)
			continue
		}

		return "", err
	}

	return "", fmt.Errorf("exhausted retries")
}

// SummarizeAllWithFallback tries Gemini, falls back to local digest on failure
func (s *Summarizer) SummarizeAllWithFallback(data map[string][]fetcher.Message) string {
	backoff := []time.Duration{20 * time.Second, 45 * time.Second}

	for attempt, wait := range backoff {
		result, err := s.SummarizeAll(data)
		if err == nil {
			return result
		}

		fmt.Printf("AI unavailable (attempt %d/2): %v\n", attempt+1, err)

		if attempt < len(backoff)-1 {
			fmt.Printf("retrying in %s...\n", wait)
			time.Sleep(wait)
		}
	}

	fmt.Println("\ngemini quota exhausted — falling back to local digest...\n")
	return LocalDigest(data)
}
