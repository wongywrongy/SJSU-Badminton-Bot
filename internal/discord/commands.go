package discord

import (
    "fmt"
    "log/slog"

    "github.com/bwmarrin/discordgo"
)

func (c *Client) registerCommands() error {
    cmds := []*discordgo.ApplicationCommand{
        {
            Name:        "macgym",
            Description: "Show current Mac Gym badminton occupancy",
        },
        {
            Name:        "badminton",
            Description: "Badminton information and events",
            Options: []*discordgo.ApplicationCommandOption{
                {
                    Type:        discordgo.ApplicationCommandOptionSubCommand,
                    Name:        "events",
                    Description: "List upcoming badminton events",
                    Options: []*discordgo.ApplicationCommandOption{
                        {
                            Type:        discordgo.ApplicationCommandOptionInteger,
                            Name:        "days",
                            Description: "Number of days to look ahead (default: 7)",
                            Required:    false,
                        },
                    },
                },
            },
        },
        {
            Name:        "subscribe",
            Description: "Subscribe to badminton alerts",
            Options: []*discordgo.ApplicationCommandOption{
                {
                    Type:        discordgo.ApplicationCommandOptionInteger,
                    Name:        "threshold",
                    Description: "Alert when occupancy reaches this level (default: 0)",
                    Required:    false,
                },
            },
        },
        {
            Name:        "unsubscribe",
            Description: "Unsubscribe from badminton alerts",
        },
    }

    for _, cmd := range cmds {
        var createdCmd *discordgo.ApplicationCommand
        var err error
        
        if c.cfg.GuildID != "" {
            // Guild-scoped command (for development)
            createdCmd, err = c.sess.ApplicationCommandCreate(c.cfg.AppID, c.cfg.GuildID, cmd)
            slog.Info("Created guild command", "name", cmd.Name, "guild", c.cfg.GuildID)
        } else {
            // Global command (for production)
            createdCmd, err = c.sess.ApplicationCommandCreate(c.cfg.AppID, "", cmd)
            slog.Info("Created global command", "name", cmd.Name)
        }
        
        if err != nil {
            return fmt.Errorf("creating command %s: %w", cmd.Name, err)
        }
        
        slog.Debug("Command created", "id", createdCmd.ID, "name", createdCmd.Name)
    }
    
    return nil
}

func (c *Client) attachHandlers() {
    c.sess.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
        if i.Type != discordgo.InteractionApplicationCommand {
            return
        }
        
        commandName := i.ApplicationCommandData().Name
        slog.Info("Command received", 
            "command", commandName, 
            "user", i.Member.User.Username,
            "guild", i.GuildID)
        
        switch commandName {
        case "macgym":
            c.handleMacGym(s, i)
        case "badminton":
            c.handleBadminton(s, i)
        case "subscribe":
            c.handleSubscribe(s, i)
        case "unsubscribe":
            c.handleUnsubscribe(s, i)
        default:
            c.ephemeral(s, i, "Unknown command: "+commandName)
        }
    })
    
    // Add ready handler
    c.sess.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
        slog.Info("Bot is ready", "user", r.User.Username, "guilds", len(r.Guilds))
    })
}

func (c *Client) ephemeral(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
    err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Content: content,
            Flags:   discordgo.MessageFlagsEphemeral,
        },
    })
    
    if err != nil {
        slog.Error("Failed to send ephemeral response", "error", err)
    }
}

func (c *Client) respond(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
    err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Content: content,
        },
    })
    
    if err != nil {
        slog.Error("Failed to send response", "error", err)
    }
}

func (c *Client) respondWithEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed) {
    err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Embeds: []*discordgo.MessageEmbed{embed},
        },
    })
    
    if err != nil {
        slog.Error("Failed to send embed response", "error", err)
    }
}
