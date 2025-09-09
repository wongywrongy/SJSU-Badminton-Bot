# Style Guide

* Idiomatic Go, small packages, clear names. Avoid globals; pass deps via constructors.
* Use contexts with deadlines for all outbound I/O.
* Log with `log/slog` (structured fields).
* Return typed errors with `fmt.Errorf("...: %w", err)`.
* Unit tests: table‑driven, small fixtures in `testdata/`.
