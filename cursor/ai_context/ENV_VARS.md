# Environment Variables

Create `.env` based on this example:

```
DISCORD_BOT_TOKEN=xxxxx
DISCORD_APP_ID=xxxxx
DISCORD_GUILD_ID=xxxxx   # for dev scoped commands; optional in prod
LOG_LEVEL=info
TIMEZONE=America/Los_Angeles
MACGYM_URL=https://www.connect2mycloud.com/Widgets/Data/locationCount?type=circle&key=92833ff9-2797-43ed-98ab-8730784a147f&loc_status=false
FITNESS_URL=https://fitness.sjsu.edu/Facility/GetSchedule
REFRESH_MACGYM_CRON=@every 2m
REFRESH_EVENTS_CRON=@every 30m
ALERT_CHANNEL_ID=        # optional
```
