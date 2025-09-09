FROM golang:1.22 as build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bot ./cmd/bot

FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=build /bot /bot
ENV TIMEZONE=America/Los_Angeles
CMD ["/bot"]
