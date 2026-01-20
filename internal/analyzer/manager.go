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
)

type Manager struct {
	cfg          *config.Config
	aiClient     *ai.Client
	msgChan      chan model.MessageData
	windowBuffer map[int64][]model.MessageData
	mu           sync.Mutex
}

func NewManager(cfg *config.Config, aiClient *ai.Client) *Manager {
	return &Manager{
		cfg:          cfg,
		aiClient:     aiClient,
		msgChan:      make(chan model.MessageData, 1000),
		windowBuffer: make(map[int64][]model.MessageData),
	}
}

// AddMessage queues a message for analysis
func (m *Manager) AddMessage(msg model.MessageData) {
	select {
	case m.msgChan <- msg:
	default:
		if m.cfg.Monitor.Debug {
			log.Println("[WARN] Message queue full, dropping message")
		}
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

	if m.cfg.Monitor.Debug {
		log.Printf("--- Monitor Report for past %v ---", window)
	}

	var wg sync.WaitGroup

	for groupID, msgs := range currentBatch {
		wg.Add(1)
		go func(gid int64, messages []model.MessageData) {
			defer wg.Done()
			m.processGroupBatch(ctx, gid, messages)
		}(groupID, msgs)
	}
	wg.Wait()
	if m.cfg.Monitor.Debug {
		log.Println("---------------------------")
	}
}

func (m *Manager) processGroupBatch(ctx context.Context, groupID int64, msgs []model.MessageData) {
	// Simple stats
	if m.cfg.Monitor.Debug {
		log.Printf("Group %d: %d messages", groupID, len(msgs))
	}

	// 1. Preprocessing
	var chatLogBuilder strings.Builder
	messageCount := 0
	for _, msg := range msgs {
		if len([]rune(msg.Text)) < 2 {
			continue
		}
		// Format: [UserID]: Content
		chatLogBuilder.WriteString(fmt.Sprintf("- %s\n", msg.Text))
		messageCount++
	}

	if messageCount == 0 {
		if m.cfg.Monitor.Debug {
			log.Printf("Group %d: No valid discussion", groupID)
		}
		return
	}

	chatLog := chatLogBuilder.String()
	if m.cfg.Monitor.Debug {
		log.Printf("[DEBUG] Group %d text to analyze:\n%s\n", groupID, chatLog)
	}

	// 2. Call LLM for analysis
	// Set timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	analysis, err := m.aiClient.Analyze(ctxWithTimeout, chatLog)
	if err != nil {
		log.Printf("Group %d LLM analysis failed: %v", groupID, err)
		return
	}

	log.Printf(">>> Group %d Analysis Result:\n%s\n", groupID, analysis)
}
