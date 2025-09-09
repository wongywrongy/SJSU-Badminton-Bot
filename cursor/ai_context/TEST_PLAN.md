# Test Plan

* **macgym.go**: Golden JSON fixture → parsed snapshot; missing fields tolerated.
* **fitness_schedule.go**: HTML and JSON fixtures → extract badminton events; timezone normalization.
* **time.go**: Round‑trip conversions, formatting.
* **store**: Dedupe and upsert semantics.
* **commands**: Mock store and verify embeds.
