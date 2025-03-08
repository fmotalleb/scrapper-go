# Builder step compiles the application into standalone binary
FROM golang:latest AS builder
RUN mkdir /app
COPY go.mod /app/
COPY go.sum /app/
WORKDIR /app
RUN go mod download
COPY ./ /app
RUN CGO_ENABLED=0 go build -o scrapper-go
RUN chmod +x scrapper-go


FROM library/debian:bookworm-slim AS runner
RUN --mount=type=cache,target=/var/lib/apt/lists/ \
  --mount=type=cache,target=/var/cache/apt/archives/ <<EOF
apt update  || exit 1
apt install -y chromium chromium-driver || exit 1
useradd -m scrapper || exit 1
EOF
USER scrapper


COPY --from=builder /app/scrapper-go /usr/bin/scrapper-go

RUN /usr/bin/scrapper-go setup

WORKDIR /home/scrapper

ENTRYPOINT ["/usr/bin/scrapper-go" ]
CMD ["-c","/config.yaml"]
