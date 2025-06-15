# Goのベースイメージを使用
FROM golang:1.21-alpine

# 必要なパッケージをインストール
RUN apk add --no-cache alsa-utils ca-certificates wget

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

# ポートを公開
EXPOSE 8080

# アプリケーションを実行
CMD ["./server"]
