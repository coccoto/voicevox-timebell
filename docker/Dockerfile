FROM golang:1.24.4-bookworm

RUN apt update && apt full-upgrade -y && apt autoremove -y && apt clean
RUN apt install -y alsa-utils cron curl

WORKDIR /app

COPY ./server /app/server
COPY ./client /app/client

RUN mkdir -p /app/storage

RUN go build -C /app/server -o /app/dist/app .

COPY crontab /etc/cron.d/voicevox-timebell-cron

RUN chmod 0644 /etc/cron.d/voicevox-timebell-cron
RUN crontab /etc/cron.d/voicevox-timebell-cron
