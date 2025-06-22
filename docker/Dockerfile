FROM golang:1.24.4-bookworm

# OS をセットアップする
RUN apt update && apt full-upgrade -y && apt autoremove -y && apt clean
RUN apt install -y alsa-utils cron

WORKDIR /app

# コンテナにファイルをコピーする
COPY ./docker/start.sh /app/start.sh
COPY ./docker/crontab /etc/cron.d/voicevox-timebell-cron
COPY ./server /app/server
COPY ./client /app/client

# cron を登録する
RUN chmod 0644 /etc/cron.d/voicevox-timebell-cron
RUN crontab /etc/cron.d/voicevox-timebell-cron

# Go プロジェクトをセットアップする
RUN mkdir -p /app/storage
RUN go build -C /app/server -o /app/dist/app .

# start.sh の権限を設定する
RUN chmod 744 /app/start.sh

CMD ["/app/start.sh"]
