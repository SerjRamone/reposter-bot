// Package app ...
package app

import (
	"log"
	"math/rand"
	"time"

	"github.com/SerjRamone/reposter-bot/config"
	"github.com/SerjRamone/reposter-bot/internal/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// App ...
type App struct {
	Bot    *tgbotapi.BotAPI
	L      *zap.Logger
	Config *config.Config
}

// New ...
func New(cfg *config.Config) *App {
	a := &App{
		L:      logger.Get(),
		Config: cfg,
	}

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Panic(err)
	}
	a.L.Info("autorized on account", zap.String("bot_name", bot.Self.UserName))
	a.Bot = bot

	return a
}

// Run ...
func (a *App) Run() {
	for {
		// loop and check channels
		for i := range a.Config.Channels {
			if a.checkChannel(&a.Config.Channels[i]) {
				go a.processChannel(&a.Config.Channels[i])
			}
		}
		time.Sleep(time.Millisecond * 1000)
	}
}

// checkChannel returns true if channel need to be reposted
func (a *App) checkChannel(c *config.Channel) bool {
	a.L.Debug("check channel", zap.Int("channel id", c.ID))
	result := int(time.Since(c.PostedAt).Minutes()) > getRandInt(c.Delay.From, c.Delay.To) && !c.IsInProcess
	a.L.Debug("channel check status", zap.Bool("is need to process", result))
	return result
}

// processChannel repost in channel
func (a *App) processChannel(c *config.Channel) {
	c.IsInProcess = true
	defer func() {
		c.IsInProcess = false
	}()
	a.L.Info("process channel", zap.Int("channel id", c.ID))

	rndMessageID := getRandInt(c.Posts.From, c.Posts.To)
	cmConfig := tgbotapi.CopyMessageConfig{
		BaseChat:   tgbotapi.BaseChat{ChatID: int64(c.ID)},
		FromChatID: int64(c.ID),
		MessageID:  rndMessageID,
		ParseMode:  "HTML",
	}

	copiedMessageID, err := a.Bot.CopyMessage(cmConfig)
	if err != nil {
		a.L.Error("copy message error", zap.Error(err), zap.Int("chat id", c.ID), zap.Int("message id", rndMessageID))
		c.PostedAt = c.PostedAt.Add(time.Duration(a.Config.RetryDelay) * time.Second)
	} else {
		a.L.Info("message copied", zap.Int("chat id", c.ID), zap.Int("message id", copiedMessageID.MessageID))
		c.PostedAt = time.Now()
	}
	a.L.Info("process channel is done", zap.Int("channel id", c.ID))
}

// getRandInt return random int from range
func getRandInt(min, max int) int {
	wdRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return wdRand.Intn(max-min) + min
}
