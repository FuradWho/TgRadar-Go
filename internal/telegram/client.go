package telegram

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/FuradWho/TgRadar-Go/internal/config"
	"github.com/FuradWho/TgRadar-Go/internal/model"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/dcs"
	"github.com/gotd/td/tg"
	"golang.org/x/net/proxy"
)

// MessageHandler processes incoming messages
type MessageHandler func(model.MessageData)

type Client struct {
	client  *telegram.Client
	cfg     *config.Config
	handler MessageHandler
}

func NewClient(cfg *config.Config, handler MessageHandler) *Client {
	return &Client{
		cfg:     cfg,
		handler: handler,
	}
}

func (c *Client) Start(ctx context.Context) error {
	dispatcher := tg.NewUpdateDispatcher()
	dispatcher.OnNewMessage(c.onNewMessage)

	opts := telegram.Options{
		SessionStorage: &telegram.FileSessionStorage{Path: c.cfg.Telegram.SessionFile},
		UpdateHandler:  dispatcher,
	}

	if c.cfg.Telegram.Proxy != "" {
		dialer, err := proxy.SOCKS5("tcp", c.cfg.Telegram.Proxy, nil, proxy.Direct)
		if err != nil {
			return fmt.Errorf("proxy config error: %w", err)
		}
		opts.Resolver = dcs.Plain(dcs.PlainOptions{
			Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {
				if d, ok := dialer.(proxy.ContextDialer); ok {
					return d.DialContext(ctx, network, addr)
				}
				return dialer.Dial(network, addr)
			},
		})
		log.Printf("Proxy enabled: %s", c.cfg.Telegram.Proxy)
	}

	c.client = telegram.NewClient(
		c.cfg.Telegram.AppID,
		c.cfg.Telegram.AppHash,
		opts,
	)

	return c.client.Run(ctx, func(ctx context.Context) error {
		flow := auth.NewFlow(
			termAuth{phone: c.cfg.Telegram.Phone, pw: c.cfg.Telegram.Password},
			auth.SendCodeOptions{},
		)

		if err := c.client.Auth().IfNecessary(ctx, flow); err != nil {
			return err
		}

		me, _ := c.client.Self(ctx)
		log.Printf("Logged in as: %s (%s), monitoring started...", me.FirstName, me.Username)

		<-ctx.Done()
		return ctx.Err()
	})
}

func (c *Client) onNewMessage(ctx context.Context, e tg.Entities, update *tg.UpdateNewMessage) error {
	msg, ok := update.Message.(*tg.Message)
	if !ok || msg.Out {
		return nil
	}

	var groupID int64
	if peer, ok := msg.PeerID.(*tg.PeerChannel); ok {
		groupID = peer.ChannelID
	} else if peer, ok := msg.PeerID.(*tg.PeerChat); ok {
		groupID = peer.ChatID
	} else {
		// Only care about groups and channels for now
		return nil
	}

	if c.cfg.Monitor.Debug {
		log.Printf("[DEBUG] Received msg from group %d: %s", groupID, msg.Message)
	}

	if msg.Message != "" {
		c.handler(model.MessageData{
			GroupID:   groupID,
			SenderID:  0, // TODO: Resolve FromID
			Text:      msg.Message,
			Timestamp: time.Unix(int64(msg.Date), 0),
		})
	}
	return nil
}

type termAuth struct {
	phone string
	pw    string
}

func (t termAuth) Phone(_ context.Context) (string, error) {
	return t.phone, nil
}

func (t termAuth) Password(_ context.Context) (string, error) {
	return t.pw, nil
}

func (t termAuth) AcceptTermsOfService(_ context.Context, tos tg.HelpTermsOfService) error {
	return nil
}

func (t termAuth) SignUp(_ context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, fmt.Errorf("sign up not supported")
}

func (t termAuth) Code(_ context.Context, sentCode *tg.AuthSentCode) (string, error) {
	fmt.Print("Enter code: ")
	code, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(code), nil
}
