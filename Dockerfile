FROM golang:1.22 as build
WORKDIR /app
COPY . .
RUN --mount=type=cache,id=go-mod,target=/go/pkg/mod \
    --mount=type=cache,id=go-build,target=/root/.cache/go-build \
    go build -o /bot ./cmd/bot

FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=build /bot /bot
ENV TIMEZONE=America/Los_Angeles
CMD ["/bot"]
