# Project Scope

**Goal:** A Discord bot for the SJSU Badminton club that:

* `/macgym` → shows live Mac Gym badminton court occupancy from the Connect2MyCloud endpoint.
* `/badminton events` → lists upcoming badminton‑related events parsed from fitness.sjsu.edu schedule.
* `/subscribe` and `/unsubscribe` → opt‑in role pings or DM alerts when new badminton events are posted or when occupancy crosses thresholds.
* Background jobs periodically refresh data and cache it.

**Non‑goals (v1):** database persistence, web dashboards, advanced analytics.

**Tech choices:** Go + discordgo, `net/http`, `goquery` (for HTML) or direct JSON decode for APIs, `robfig/cron/v3` for scheduling. Avoid heavy frameworks.
