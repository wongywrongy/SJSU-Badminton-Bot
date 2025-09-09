# SJSU Badminton Discord Bot

A Discord bot for the SJSU Badminton club that provides real-time Mac Gym occupancy data and upcoming badminton events.

## Features

- **`/macgym`** - Shows current Mac Gym badminton court occupancy
- **`/badminton events [days]`** - Lists upcoming badminton events (default: 7 days)
- **`/subscribe [threshold]`** - Subscribe to alerts when occupancy crosses thresholds
- **`/unsubscribe`** - Unsubscribe from alerts
- Background jobs that refresh data every 2 minutes (Mac Gym) and 30 minutes (events)

## Quick Start

### Prerequisites

- Go 1.22 or later
- Discord Bot Token and Application ID
- Discord server with bot permissions

### Setup

1. **Clone and initialize:**
   ```bash
   git clone <your-repo>
   cd badminton-discord-bot
   go mod tidy
   ```

2. **Configure environment:**
   ```bash
   cp env.example .env
   # Edit .env with your Discord bot token and other settings
   ```

3. **Run locally:**
   ```bash
   make run
   ```

## Discord Bot Setup

### Creating a Discord Application

1. Go to [Discord Developer Portal](https://discord.com/developers/applications)
2. Click "New Application" and give it a name
3. Go to the "Bot" section and click "Add Bot"
4. Copy the bot token and add it to your `.env` file as `DISCORD_BOT_TOKEN`
5. Copy the Application ID and add it to your `.env` file as `DISCORD_APP_ID`

### Bot Permissions

Your bot needs these permissions:
- Send Messages
- Use Slash Commands
- Embed Links
- Read Message History

### Inviting the Bot

1. In the Discord Developer Portal, go to OAuth2 > URL Generator
2. Select scopes: `bot` and `applications.commands`
3. Select bot permissions: Send Messages, Use Slash Commands, Embed Links, Read Message History
4. Copy the generated URL and open it to invite the bot to your server

### Guild ID (for Development)

For development, you can scope commands to a specific guild:
1. Right-click your Discord server name
2. Click "Copy Server ID"
3. Add this to your `.env` file as `DISCORD_GUILD_ID`

For production, leave `DISCORD_GUILD_ID` empty to register global commands.

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DISCORD_BOT_TOKEN` | Discord bot token (required) | - |
| `DISCORD_APP_ID` | Discord application ID | - |
| `DISCORD_GUILD_ID` | Guild ID for dev commands (optional) | - |
| `LOG_LEVEL` | Logging level | `info` |
| `TIMEZONE` | Timezone for events | `America/Los_Angeles` |
| `MACGYM_URL` | Mac Gym occupancy API URL | (provided) |
| `FITNESS_URL` | SJSU Fitness schedule URL | (provided) |
| `REFRESH_MACGYM_CRON` | Mac Gym refresh schedule | `@every 2m` |
| `REFRESH_EVENTS_CRON` | Events refresh schedule | `@every 30m` |
| `ALERT_CHANNEL_ID` | Channel for alerts (optional) | - |

## Development

### Running Tests

```bash
make test
```

### Building

```bash
make build
```

### Docker

```bash
make docker-build
```

## Data Sources

### Mac Gym Occupancy
- **Source**: Connect2MyCloud API
- **URL**: `https://www.connect2mycloud.com/Widgets/Data/locationCount?type=circle&key=92833ff9-2797-43ed-98ab-8730784a147f&loc_status=false`
- **Format**: JSON
- **Refresh**: Every 2 minutes

### SJSU Fitness Schedule
- **Source**: SJSU Fitness website
- **URL**: `https://fitness.sjsu.edu/Facility/GetSchedule`
- **Format**: HTML/JSON (auto-detected)
- **Refresh**: Every 30 minutes
- **Filter**: Badminton-related events only

## Commands

### `/macgym`
Shows current Mac Gym badminton court occupancy with last updated timestamp.

### `/badminton events [days]`
Lists upcoming badminton events for the specified number of days (default: 7).

### `/subscribe [threshold]`
Subscribe to alerts. If threshold is specified, only alerts when occupancy is at or above that level.

### `/unsubscribe`
Remove your subscription to alerts.

## Architecture

```
badminton-discord-bot/
├─ cmd/bot/main.go              # Application entry point
├─ internal/
│  ├─ config/                   # Configuration management
│  ├─ discord/                  # Discord client and handlers
│  ├─ scrape/                   # Data scraping modules
│  ├─ store/                    # In-memory data store
│  ├─ sched/                    # Cron job scheduler
│  └─ util/                     # Utility functions
├─ cursor/                      # Development context files
└─ infra/                       # Deployment files
```

## Troubleshooting

### Bot Not Responding
1. Check that the bot token is correct in `.env`
2. Verify the bot has proper permissions in your Discord server
3. Check the console logs for error messages

### Commands Not Appearing
1. Ensure `DISCORD_APP_ID` is set correctly
2. For development, make sure `DISCORD_GUILD_ID` is set
3. Commands may take up to 1 hour to appear globally (without guild ID)

### Data Not Updating
1. Check the console logs for scraping errors
2. Verify the API endpoints are accessible
3. Check your internet connection

### Permission Errors
1. Ensure the bot has "Use Slash Commands" permission
2. Check that the bot role is above other roles that might restrict it

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is licensed under the MIT License.
