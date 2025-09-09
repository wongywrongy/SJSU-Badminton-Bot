package config

import (
    "errors"
    "os"
)

type Config struct {
    Token      string
    AppID      string
    GuildID    string
    TZ         string
    MacGymURL  string
    FitnessURL string
    CronMacGym string
    CronEvents string
    AlertChan  string
}

func get(k, def string) string { if v := os.Getenv(k); v != "" { return v }; return def }

func Load() (Config, error) {
    c := Config{
        Token:      os.Getenv("DISCORD_BOT_TOKEN"),
        AppID:      get("DISCORD_APP_ID", ""),
        GuildID:    get("DISCORD_GUILD_ID", ""),
        TZ:         get("TIMEZONE", "America/Los_Angeles"),
        MacGymURL:  get("MACGYM_URL", "https://www.connect2mycloud.com/Widgets/Data/locationCount?type=circle&key=92833ff9-2797-43ed-98ab-8730784a147f&loc_status=false"),
        FitnessURL: get("FITNESS_URL", "https://fitness.sjsu.edu/Facility/GetSchedule"),
        CronMacGym: get("REFRESH_MACGYM_CRON", "@every 2m"),
        CronEvents: get("REFRESH_EVENTS_CRON", "@every 30m"),
        AlertChan:  get("ALERT_CHANNEL_ID", ""),
    }
    if c.Token == "" { return c, errors.New("missing DISCORD_BOT_TOKEN") }
    return c, nil
}
