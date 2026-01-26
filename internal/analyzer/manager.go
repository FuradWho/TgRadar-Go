package analyzer

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/FuradWho/TgRadar-Go/internal/ai"
	"github.com/FuradWho/TgRadar-Go/internal/config"
	"github.com/FuradWho/TgRadar-Go/internal/model"
	"github.com/FuradWho/TgRadar-Go/internal/notifier"
)

type Manager struct {
	cfg          *config.Config
	aiClient     *ai.Client
	notifier     notifier.Sender
	msgChan      chan model.MessageData
	windowBuffer map[int64][]model.MessageData
	mu           sync.Mutex
}

const globalSummaryBanner = "\n====== GLOBAL INTELLIGENCE SUMMARY ======\n%s\n========================================="

func NewManager(cfg *config.Config, aiClient *ai.Client, notifier notifier.Sender) *Manager {
	return &Manager{
		cfg:          cfg,
		aiClient:     aiClient,
		notifier:     notifier,
		msgChan:      make(chan model.MessageData, 1000),
		windowBuffer: make(map[int64][]model.MessageData),
	}
}

// AddMessage queues a message for analysis
func (m *Manager) AddMessage(msg model.MessageData) {
	select {
	case m.msgChan <- msg:
	default:
		m.debugf("[WARN] Message queue full, dropping message")
	}
}

// Start runs the main analysis loop
func (m *Manager) Start(ctx context.Context) {
	windowDuration := time.Duration(m.cfg.Monitor.WindowSeconds) * time.Second
	ticker := time.NewTicker(windowDuration)
	defer ticker.Stop()

	log.Printf("Analyzer started, monitor window: %v", windowDuration)

	for {
		select {
		case msg := <-m.msgChan:
			m.mu.Lock()
			m.windowBuffer[msg.GroupID] = append(m.windowBuffer[msg.GroupID], msg)
			m.mu.Unlock()

		case <-ticker.C:
			m.analyzeAndPrint(ctx, windowDuration)

		case <-ctx.Done():
			return
		}
	}
}

func (m *Manager) analyzeAndPrint(ctx context.Context, window time.Duration) {
	m.mu.Lock()
	currentBatch := m.windowBuffer
	m.windowBuffer = make(map[int64][]model.MessageData)
	m.mu.Unlock()

	if len(currentBatch) == 0 {
		return
	}

	m.debugf("--- Monitor Report for past %v ---", window)

	var wg sync.WaitGroup
	var summaries []string
	var summariesMu sync.Mutex

	for groupID, msgs := range currentBatch {
		wg.Add(1)
		go func(gid int64, messages []model.MessageData) {
			defer wg.Done()
			summary := m.processGroupBatch(ctx, gid, messages)
			if summary != "" {
				summariesMu.Lock()
				summaries = append(summaries, formatGroupReport(gid, summary))
				summariesMu.Unlock()
			}
		}(groupID, msgs)
	}
	wg.Wait()

	if len(summaries) > 0 {
		m.processGlobalSummary(ctx, summaries)
	}

	m.debugf("---------------------------")
}

func (m *Manager) processGlobalSummary(ctx context.Context, summaries []string) {
	combinedReport := strings.Join(summaries, "\n\n---\n\n")

	m.debugf("Generating Global Summary...")

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	summary, err := m.aiClient.AnalyzeSummary(ctxWithTimeout, combinedReport)
	if err != nil {
		log.Printf("Global summary failed: %v", err)
		return
	}

	log.Printf(globalSummaryBanner, summary)
	if m.notifier != nil {
		if err := m.notifier.Send(ctx, fmt.Sprintf("Global Summary:\n%s", summary)); err != nil {
			log.Printf("Notifier send failed: %v", err)
		}
	}
}

func (m *Manager) processGroupBatch(ctx context.Context, groupID int64, msgs []model.MessageData) string {
	// Simple stats
	m.debugf("Group %d: %d messages", groupID, len(msgs))

	// 1. Preprocessing
	var chatLogBuilder strings.Builder
	messageCount := 0
	for _, msg := range msgs {
		if len([]rune(msg.Text)) < 2 {
			continue
		}
		// Include sender ID for accurate unique counts
		if msg.SenderID != 0 {
			chatLogBuilder.WriteString(fmt.Sprintf("- U%d: %s\n", msg.SenderID, msg.Text))
		} else {
			chatLogBuilder.WriteString(fmt.Sprintf("- U?: %s\n", msg.Text))
		}
		messageCount++
	}

	if messageCount == 0 {
		m.debugf("Group %d: No valid discussion", groupID)
		return ""
	}

	chatLog := chatLogBuilder.String()
	m.debugf("[DEBUG] Group %d text to analyze:\n%s\n", groupID, chatLog)

	// 2. Call LLM for analysis
	// Set timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	analysis, err := m.aiClient.Analyze(ctxWithTimeout, chatLog)
	if err != nil {
		log.Printf("Group %d LLM analysis failed: %v", groupID, err)
		return ""
	}

	m.debugf(">>> Group %d Analysis Result:\n%s\n", groupID, analysis)
	return analysis
}

func (m *Manager) debugf(format string, args ...any) {
	if m.cfg.Monitor.Debug {
		log.Printf(format, args...)
	}
}

func formatGroupReport(groupID int64, summary string) string {
	return fmt.Sprintf("Group %d Report:\n%s", groupID, summary)
}
