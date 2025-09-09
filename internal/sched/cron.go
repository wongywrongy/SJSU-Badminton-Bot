package sched

import (
    "context"
    "log/slog"
    "math/rand"
    "time"

    "github.com/robfig/cron/v3"

    "github.com/sjsu-badminton/badminton-discord-bot/internal/config"
    "github.com/sjsu-badminton/badminton-discord-bot/internal/scrape"
    "github.com/sjsu-badminton/badminton-discord-bot/internal/store"
    "github.com/sjsu-badminton/badminton-discord-bot/internal/util"
)

type Cron struct {
    c     *cron.Cron
    store *store.MemoryStore
}

func Start(ctx context.Context, cfg config.Config, st *store.MemoryStore) *Cron {
    loc := util.MustLocation(cfg.TZ)
    
    // Create cron with location and logger
    c := cron.New(
        cron.WithLocation(loc),
        cron.WithLogger(cron.VerbosePrintfLogger(slog.NewLogLogger(slog.Default().Handler(), slog.LevelInfo))),
    )

    cronJob := &Cron{
        c:     c,
        store: st,
    }

    // Add Mac Gym refresh job with jitter
    c.AddFunc(cfg.CronMacGym, cronJob.refreshMacGym(cfg))

    // Add events refresh job with jitter
    c.AddFunc(cfg.CronEvents, cronJob.refreshEvents(cfg, loc))

    // Start with a small delay to avoid thundering herd
    go func() {
        jitter := time.Duration(rand.Intn(30)) * time.Second
        time.Sleep(jitter)
        c.Start()
        slog.Info("Cron scheduler started", 
            "macGymSchedule", cfg.CronMacGym,
            "eventsSchedule", cfg.CronEvents,
            "timezone", cfg.TZ)
    }()

    return cronJob
}

func (cr *Cron) Stop() {
    slog.Info("Stopping cron scheduler...")
    ctx := cr.c.Stop()
    <-ctx.Done()
    slog.Info("Cron scheduler stopped")
}

func (cr *Cron) refreshMacGym(cfg config.Config) func() {
    return func() {
        start := time.Now()
        
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        
        snap, err := scrape.FetchMacGym(ctx, cfg.MacGymURL)
        if err != nil {
            slog.Error("Failed to fetch Mac Gym data", 
                "error", err,
                "duration", time.Since(start))
            return
        }
        
        cr.store.SetMac(snap)
        
        slog.Info("Mac Gym data refreshed", 
            "capacity", snap.Capacity,
            "inUse", snap.InUse,
            "duration", time.Since(start))
    }
}

func (cr *Cron) refreshEvents(cfg config.Config, loc *time.Location) func() {
    return func() {
        start := time.Now()
        
        ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
        defer cancel()
        
        events, err := scrape.FetchBadmintonEvents(ctx, cfg.FitnessURL, loc)
        if err != nil {
            slog.Error("Failed to fetch fitness events", 
                "error", err,
                "duration", time.Since(start))
            return
        }
        
        cr.store.UpsertEvents(events)
        
        slog.Info("Fitness events refreshed", 
            "eventsFound", len(events),
            "totalEvents", cr.store.GetEventCount(),
            "duration", time.Since(start))
    }
}
