# Commands

**Slash commands (guild‑scoped for dev, global for prod):**

* `/macgym` → Show current occupancy (e.g., "Courts in use: 6/8", or best available metric). Include last updated time.
* `/badminton events [days:int=7]` → List upcoming badminton events for the next N days. Paginate if >10.
* `/subscribe [threshold:int=0]` → Subscribe the invoker to alerts. If `threshold>0`, alert only when occupancy ≥ threshold.
* `/unsubscribe` → Remove subscription.

**UX details**

* Use embeds with title, fields, timestamp. Ephemeral replies for command confirmations; public for /macgym and /badminton events (configurable by flag).
