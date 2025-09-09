package discord

import (
    "context"
    "fmt"
    "log/slog"

    "github.com/bwmarrin/discordgo"

    "github.com/sjsu-badminton/badminton-discord-bot/internal/config"
    "github.com/sjsu-badminton/badminton-discord-bot/internal/sched"
    "github.com/sjsu-badminton/badminton-discord-bot/internal/store"
)

type Client struct {
    cfg   config.Config
    sess  *discordgo.Session
    store *store.MemoryStore
    cron  *sched.Cron
}

func NewClient(ctx context.Context, cfg config.Config) (*Client, error) {
    s, err := discordgo.New("Bot " + cfg.Token)
    if err != nil {
        return nil, fmt.Errorf("creating Discord session: %w", err)
    }
    
    // Set intents
    s.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages
    
    c := &Client{
        cfg:   cfg,
        sess:  s,
        store: store.NewMemoryStore(),
    }
    
    c.attachHandlers()
    
    slog.Info("Discord client created successfully")
    return c, nil
}

func (c *Client) Start(ctx context.Context) error {
    slog.Info("Starting Discord bot...")
    
    if err := c.sess.Open(); err != nil {
        return fmt.Errorf("opening Discord session: %w", err)
    }
    
    if err := c.registerCommands(); err != nil {
        return fmt.Errorf("registering commands: %w", err)
    }
    
    c.cron = sched.Start(ctx, c.cfg, c.store)
    
    slog.Info("Bot started successfully", 
        "guildID", c.cfg.GuildID,
        "appID", c.cfg.AppID)
    
    return nil
}

func (c *Client) Stop() {
    slog.Info("Stopping Discord bot...")
    
    if c.cron != nil {
        c.cron.Stop()
    }
    
    if c.sess != nil {
        c.sess.Close()
    }
    
    slog.Info("Bot stopped")
}

// SendAlert sends an alert to a user or channel
func (c *Client) SendAlert(userID string, message string) error {
    if c.cfg.AlertChan != "" {
        // Send to alert channel
        _, err := c.sess.ChannelMessageSend(c.cfg.AlertChan, fmt.Sprintf("<@%s> %s", userID, message))
        return err
    } else {
        // Send DM
        channel, err := c.sess.UserChannelCreate(userID)
        if err != nil {
            return fmt.Errorf("creating DM channel: %w", err)
        }
        
        _, err = c.sess.ChannelMessageSend(channel.ID, message)
        return err
    }
}
