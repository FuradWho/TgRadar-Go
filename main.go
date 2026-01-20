package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/FuradWho/TgRadar-Go/internal/ai"
	"github.com/FuradWho/TgRadar-Go/internal/analyzer"
	"github.com/FuradWho/TgRadar-Go/internal/config"
	"github.com/FuradWho/TgRadar-Go/internal/telegram"
)

func main() {
	// 1. Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// 2. Initialize AI client
	aiClient := ai.NewClient(cfg)

	// 3. Initialize Analyzer
	anal := analyzer.NewManager(cfg, aiClient)

	// 4. Initialize Telegram client
	// Bind message handler to analyzer's AddMessage
	tgClient := telegram.NewClient(cfg, anal.AddMessage)

	// 5. Start service
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Start analysis loop (background)
	go anal.Start(ctx)

	log.Println("Connecting to Telegram...")

	// Start Telegram client (blocking until ctx cancelled or error)
	if err := tgClient.Start(ctx); err != nil {
		// If error is due to ctx cancellation, it's not a fatal error
		if ctx.Err() == nil {
			log.Fatalf("Telegram client error: %v", err)
		}
	}

	log.Println("Service stopped")
	// Give some time for background tasks to clean up
	time.Sleep(time.Second)
}
