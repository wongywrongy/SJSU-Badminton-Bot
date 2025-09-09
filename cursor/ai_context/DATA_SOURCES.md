# Data Sources

**Mac Gym live count (JSON)**

* URL: `https://www.connect2mycloud.com/Widgets/Data/locationCount?type=circle&key=92833ff9-2797-43ed-98ab-8730784a147f&loc_status=false`
* Expect JSON. Parse robustly: tolerate unknown fields and nulls. If schema uncertain, decode into a `map[string]any` then map into typed fields.
* Cache ETag/Last‑Modified if provided. Use conditional GET when possible.

**SJSU Fitness schedule (events)**

* URL: `https://fitness.sjsu.edu/Facility/GetSchedule`
* This endpoint may require query params or cookies. Try GET first; if HTML, parse with `goquery` for rows/cards that indicate badminton. If JSON, decode. Normalize dates to timezone `America/Los_Angeles`.
* Maintain a stable `Event` model with ISO8601 `start`/`end` and a `source_url`.
