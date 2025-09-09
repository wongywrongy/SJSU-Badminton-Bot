package discord

import (
    "fmt"
    "time"

    "github.com/bwmarrin/discordgo"
)

func (c *Client) handleMacGym(s *discordgo.Session, i *discordgo.InteractionCreate) {
    snap := c.store.GetMac()
    
    // Create embed
    embed := &discordgo.MessageEmbed{
        Title:       "🏸 Mac Gym — Badminton Occupancy",
        Description: snap.Details,
        Color:       0x00ff00, // Green
        Timestamp:   snap.RetrievedAt.Format(time.RFC3339),
        Footer: &discordgo.MessageEmbedFooter{
            Text: "SJSU Badminton Bot",
        },
    }
    
    // Add capacity information if available
    if snap.Capacity > 0 {
        embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
            Name:   "Courts in Use",
            Value:  fmt.Sprintf("%d / %d", snap.InUse, snap.Capacity),
            Inline: true,
        })
        
        // Add availability percentage
        availability := float64(snap.Capacity-snap.InUse) / float64(snap.Capacity) * 100
        embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
            Name:   "Availability",
            Value:  fmt.Sprintf("%.1f%%", availability),
            Inline: true,
        })
        
        // Set color based on availability
        if availability > 50 {
            embed.Color = 0x00ff00 // Green
        } else if availability > 25 {
            embed.Color = 0xffaa00 // Orange
        } else {
            embed.Color = 0xff0000 // Red
        }
    }
    
    // Add last updated info
    embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
        Name:   "Last Updated",
        Value:  snap.RetrievedAt.Format("Mon, Jan 2 3:04 PM"),
        Inline: true,
    })
    
    c.respondWithEmbed(s, i, embed)
}

func (c *Client) handleBadminton(s *discordgo.Session, i *discordgo.InteractionCreate) {
    opts := i.ApplicationCommandData().Options
    days := 7
    
    // Extract days parameter from subcommand
    if len(opts) > 0 && len(opts[0].Options) > 0 {
        if v := opts[0].Options[0].IntValue(); v > 0 {
            days = int(v)
        }
    }
    
    events := c.store.ListUpcoming(time.Now(), days)
    
    if len(events) == 0 {
        embed := &discordgo.MessageEmbed{
            Title:       "🏸 Upcoming Badminton Events",
            Description: fmt.Sprintf("No badminton events found in the next %d days.", days),
            Color:       0x0099ff,
            Footer: &discordgo.MessageEmbedFooter{
                Text: "SJSU Badminton Bot",
            },
        }
        c.respondWithEmbed(s, i, embed)
        return
    }
    
    // Create embed with events
    embed := &discordgo.MessageEmbed{
        Title:       fmt.Sprintf("🏸 Upcoming Badminton Events (%d days)", days),
        Description: fmt.Sprintf("Found %d upcoming badminton events:", len(events)),
        Color:       0x0099ff,
        Footer: &discordgo.MessageEmbedFooter{
            Text: "SJSU Badminton Bot",
        },
    }
    
    // Add events as fields (Discord limit is 25 fields)
    maxEvents := 10
    if len(events) > maxEvents {
        events = events[:maxEvents]
        embed.Description += fmt.Sprintf(" (showing first %d)", maxEvents)
    }
    
    for _, event := range events {
        fieldValue := fmt.Sprintf("**Time:** %s - %s\n**Location:** %s",
            event.Start.Format("Mon, Jan 2 3:04 PM"),
            event.End.Format("3:04 PM"),
            event.Location)
        
        embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
            Name:   event.Title,
            Value:  fieldValue,
            Inline: false,
        })
    }
    
    c.respondWithEmbed(s, i, embed)
}

func (c *Client) handleSubscribe(s *discordgo.Session, i *discordgo.InteractionCreate) {
    threshold := 0
    
    // Extract threshold parameter
    if len(i.ApplicationCommandData().Options) > 0 {
        threshold = int(i.ApplicationCommandData().Options[0].IntValue())
    }
    
    c.store.Subscribe(i.Member.User.ID, threshold)
    
    var message string
    if threshold > 0 {
        message = fmt.Sprintf("✅ Subscribed to alerts! You'll be notified when Mac Gym occupancy reaches %d or higher.", threshold)
    } else {
        message = "✅ Subscribed to alerts! You'll be notified about new badminton events and Mac Gym updates."
    }
    
    c.ephemeral(s, i, message)
}

func (c *Client) handleUnsubscribe(s *discordgo.Session, i *discordgo.InteractionCreate) {
    c.store.Unsubscribe(i.Member.User.ID)
    c.ephemeral(s, i, "✅ You have been unsubscribed from all badminton alerts.")
}
