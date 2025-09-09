package main

import (
    "context"
    "log/slog"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/sjsu-badminton/badminton-discord-bot/internal/config"
    "github.com/sjsu-badminton/badminton-discord-bot/internal/discord"
)

func main() {
    // Set up structured logging
    logLevel := slog.LevelInfo
    if level := os.Getenv("LOG_LEVEL"); level != "" {
        switch level {
        case "debug":
            logLevel = slog.LevelDebug
        case "info":
            logLevel = slog.LevelInfo
        case "warn":
            logLevel = slog.LevelWarn
        case "error":
            logLevel = slog.LevelError
        }
    }
    
    handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: logLevel,
    })
    logger := slog.New(handler)
    slog.SetDefault(logger)
    
    slog.Info("Starting SJSU Badminton Discord Bot")
    
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        slog.Error("Failed to load configuration", "error", err)
        os.Exit(1)
    }
    
    slog.Info("Configuration loaded", 
        "appID", cfg.AppID,
        "guildID", cfg.GuildID,
        "timezone", cfg.TZ)

    // Set up graceful shutdown
    ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer cancel()

    // Create and start bot
    bot, err := discord.NewClient(ctx, cfg)
    if err != nil {
        slog.Error("Failed to create Discord client", "error", err)
        os.Exit(1)
    }
    
    if err := bot.Start(ctx); err != nil {
        slog.Error("Failed to start bot", "error", err)
        os.Exit(1)
    }
    
    slog.Info("Bot is running. Press Ctrl+C to stop.")
    
    // Wait for shutdown signal
    <-ctx.Done()
    
    slog.Info("Shutdown signal received, stopping bot...")
    
    // Give the bot time to clean up
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer shutdownCancel()
    
    go func() {
        bot.Stop()
        shutdownCancel()
    }()
    
    <-shutdownCtx.Done()
    
    if shutdownCtx.Err() == context.DeadlineExceeded {
        slog.Warn("Bot shutdown timed out")
    } else {
        slog.Info("Bot stopped gracefully")
    }
}
