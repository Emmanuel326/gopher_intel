package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/emmanuel326/gopher_intel/internal/ai"
	"github.com/emmanuel326/gopher_intel/internal/fetcher"
	"github.com/emmanuel326/gopher_intel/internal/fetcher/lore"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY not set in .env")
	}

	summarizer, err := ai.New(apiKey)
	if err != nil {
		log.Fatalf("init ai: %v", err)
	}

	fetcher.Register(lore.New("lkml"))
	fetcher.Register(lore.New("linux-block"))
	fetcher.Register(lore.New("linux-nvme"))
	fetcher.Register(lore.New("io-uring"))
	fetcher.Register(lore.New("qemu-devel"))
	fetcher.Register(lore.New("bpf"))

	// Step 1: fetch all sources
	allData := make(map[string][]fetcher.Message)
	for name, src := range fetcher.All() {
		fmt.Printf("fetching %s...\n", name)
		messages, err := src.Fetch()
		if err != nil {
			log.Printf("error fetching %s: %v", name, err)
			continue
		}
		fmt.Printf("  got %d messages\n", len(messages))
		allData[name] = messages
	}

	// Step 2: single AI call with local fallback
	fmt.Printf("\nall sources fetched (%d lists). running analysis...\n\n", len(allData))
	brief := summarizer.SummarizeAllWithFallback(allData)
	fmt.Println(brief)
}
