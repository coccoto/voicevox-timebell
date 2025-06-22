FROM golang:1.24.4-bookworm

# OS をアップデートする
RUN apt update && apt full-upgrade -y && apt autoremove -y && apt clean
# パッケージをインストールする
RUN apt install -y alsa-utils cron curl

# 作業ディレクトリを設定する
WORKDIR /app

# プログラムをコンテナにコピーする
COPY ./server /app/server
COPY ./client /app/client

# ストレージディレクトリを作成する
RUN mkdir -p /app/storage

# プログラムをビルドして実行可能ファイルを生成する
RUN go build -C /app/server -o /app/dist/app .
