# Tasks

1. **Bootstrap**

* Init Go module, add minimal main, load env, start Discord session, register slash commands, graceful shutdown.

2. **HTTP Utilities**

* Build `internal/util/httpx.go` with a shared `*http.Client` (timeout 10s), retry (exponential backoff, 3 tries on 5xx/connection errors), context support, and JSON/HTML helpers.

3. **Scrapers**

* `macgym.go`: Fetch JSON, parse into a struct. Extract badminton‑relevant counters. Return a typed snapshot.
* `fitness_schedule.go`: POST/GET schedule, parse HTML or JSON. Extract events with title, facility/room, start/end, and tags that indicate badminton. Normalize to `America/Los_Angeles`.

4. **Store**

* In‑memory store for latest `MacGymSnapshot` and list of `Event`. Include dedupe by `(title,start,end,location)` key. Provide getter/setter with copy‑on‑read to avoid races; guard with mutex.

5. **Scheduler**

* Cron tasks to refresh MacGym data every 2 minutes and fitness schedule every 30 minutes. Jitter start to avoid thundering herd.

6. **Discord Commands & Handlers**

* `/macgym` returns a neat embed with occupancy and last updated.
* `/badminton events [days:int=7]` lists upcoming events for the window.
* `/subscribe [threshold:int]` + `/unsubscribe` manage a memory set of user IDs.

7. **Alerts**

* When occupancy crosses threshold, send a DM or channel message (configurable). Debounce to avoid spam.

8. **Testing**

* Unit tests for scrapers (use golden files / fixtures), time helpers, and dedupe logic.

9. **DevOps**

* Dockerfile, Makefile, `.env.example`. CI‑friendly build flags.
