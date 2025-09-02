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


FROM library/debian:bookworm-slim AS slim
RUN --mount=type=cache,target=/var/lib/apt/lists/ \
  --mount=type=cache,target=/var/cache/apt/archives/ <<EOF
set -e
apt-get update
apt-get install -y ca-certificates  ca-certificates \
    fonts-liberation \
    libasound2 \
    libatk-bridge2.0-0 \
    libatk1.0-0 \
    libatspi2.0-0 \
    libc6 \
    libcairo2 \
    libcups2 \
    libdbus-1-3 \
    libdrm2 \
    libexpat1 \
    libgbm1 \
    libglib2.0-0 \
    libgtk-3-0 \
    libnspr4 \
    libnss3 \
    libpango-1.0-0 \
    libx11-6 \
    libx11-xcb1 \
    libxcomposite1 \
    libxdamage1 \
    libxext6 \
    libxfixes3 \
    libxkbcommon0 \
    libxrandr2 \
    wget
rm -rf /var/lib/apt/lists/*
useradd -m scrapper || exit 1
EOF
USER scrapper

COPY scrapper-go /usr/bin/scrapper-go

RUN /usr/bin/scrapper-go setup

WORKDIR /home/scrapper

ENTRYPOINT ["/usr/bin/scrapper-go" ]
CMD ["--config","/config.yaml"]
