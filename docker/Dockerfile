# Goのベースイメージを使用
FROM golang:1.21-alpine

# 必要なパッケージをインストール
RUN apk add --no-cache alsa-utils pulseaudio-utils ffmpeg ca-certificates wget

# audioグループを作成し、非rootユーザーを追加
RUN addgroup -g 29 audio && adduser -D -u 1000 appuser && adduser appuser audio

# 作業ディレクトリを設定
WORKDIR /app

# バックエンドファイルをコピー
COPY backend/ ./

# フロントエンドファイルをコピー
COPY frontend/ ./frontend/

# goモジュールを初期化（存在しない場合）
RUN if [ ! -f go.mod ]; then go mod init voicevox-timebell; fi

# アプリケーションをビルド
RUN go build -o server .

# 実行権限を設定
RUN chown appuser:appuser /app/server

# 非rootユーザーに切り替え
USER appuser

# ポートを公開
EXPOSE 8080

# アプリケーションを実行
CMD ["./server"]
